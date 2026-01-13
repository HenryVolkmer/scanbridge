package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"image/png"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/netip"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

//go:embed web/dist/app.js
var webJs embed.FS
//go:embed web/src/index.html
var entryIndexFile embed.FS

var debug *bool
var config *Config

const pdfStorageDir string = "/var/tmp/scanbridge"

var env *Environment
var scanimageBin *string
var deviceURI *string
var deviceSource *string
var scanFormat string = "png"
var scanResolution int16 = 200

func main() {

	//scanimageBin = flag.String("scanimage", "", "scanimage binary path. If not set, autodiscover is used. E.g.: /usr/bin/scanimage")
	//deviceSource = flag.String("source", "ADF", "scan source to query from")
	debug = flag.Bool("debug", false, "enable debug mode")
	bindingAddrPort := flag.String("bind", "127.0.0.1:8080", "Binding to")
	configFile := flag.String("config", "", "configfile")

	flag.Parse()

	if *debug == true {
		log.Println("** DEBUG MODE **")
	}

	if len(*configFile) == 0 {
		log.Fatalln("config parameter missing! -config=/my/conf/config.json")
	}
	var err error
	config, err = NewConfig(*configFile)
	if err != nil {
		log.Fatalln("Error loading configfile:", err)
	}
	config.IsDebug = *debug

	if *debug == true && config.Smtp == nil {
		log.Println("DEBUG:no smpt-config provided, we wont send Mails!")
	}

	env = NewEnvironment(config)

	bindAddrPort := netip.MustParseAddrPort(*bindingAddrPort)
	log.Printf("Starting webserver on %s...", bindAddrPort.String())

	sub, _ := fs.Sub(entryIndexFile, "web/src")
	http.Handle("/", http.FileServer(http.FS(sub)))

	webJsSub, _ := fs.Sub(webJs, "web/dist")
	http.Handle("/app.js", http.FileServer(http.FS(webJsSub)))
	
	http.Handle("/api/devices", NewDevicesController(config))
	http.HandleFunc("/api/env", envCtrl)
	http.HandleFunc("/api/scan", scanCtrl)
	http.HandleFunc("/api/download/", pdfDownloadCtrl)
	log.Fatalln(http.ListenAndServe(bindAddrPort.String(), nil))
}

func pdfDownloadCtrl(w http.ResponseWriter, r *http.Request) {

	uuid := strings.TrimPrefix(r.URL.Path, "/api/download/")
	if uuid == "" || strings.Contains(uuid, "..") {
		http.NotFound(w, r)
		return
	}
	filename := fmt.Sprintf("%s.pdf", uuid)
	pdfPath := filepath.Join(pdfStorageDir, filename)

	f, err := os.Open(pdfPath)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	defer f.Close()
	
	stat, err := f.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(`attachment; filename="%s"`, filename),
	)
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))

	if _, err := io.Copy(w, f); err != nil {
		log.Printf("download aborted: %v", err)
		return
	}
}

type Notification struct {
	Title string
	Data string
	URL string `json:"url"`
}

func scanCtrl(w http.ResponseWriter, r *http.Request) {
	
	mode := r.URL.Query().Get("mode")
	if mode == "" {
		mode = "Color"
	}
	uuid, err := scan(mode)
	if err != nil {
		log.Printf("Err: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(&Notification{Data: "Prüfe, ob der Scanner eingeschalten ist (Ein/Aus-Taste darf nicht blinken) und Papier im Schnelleinzug liegt. Beim Einlegen des Papiers wird der Scanner ein kurzen Ton wiedergeben.", Title: "Scan kann nicht ausgeführt werden!"})
		return
	}
	json.NewEncoder(w).Encode(&Notification{
		Data: "Der Scan war erfolgreich!", 
		Title: "OK!",
		URL: fmt.Sprintf("/api/download/%s", uuid.String()),
	})
}

func envCtrl(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(env); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func scan(mode string) (uuid.UUID, error) {

	uuid := uuid.New()	
	cwd, err := os.MkdirTemp("", "scanbridge*")
	if err != nil {
		log.Printf("Err: %s", err)
		return uuid, err
	}

	if *debug == false {
		defer os.RemoveAll(cwd)
	}

	log.Println("id", uuid.String(), "scanTo:", cwd, "Mode:", mode)
		
	// Create command.
	cmd := exec.Command(
		*scanimageBin,
		fmt.Sprintf("--device-name=%s", *deviceURI),
		fmt.Sprintf("--source=%s", *deviceSource),
		fmt.Sprintf("--format=%s", scanFormat),
		fmt.Sprintf("--resolution=%d", scanResolution),
		fmt.Sprintf("--batch=%s/%%d.png", cwd),
		fmt.Sprintf("--mode=%s", mode),
		"--batch-start=10",
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Printf("Err: %s | %s", err, stderr.String())
		return uuid, err
	}

	err = os.MkdirAll(pdfStorageDir, 0700)
	if err != nil {
		return uuid, err
	}
	pdfFileName := filepath.Join(pdfStorageDir, fmt.Sprintf("%s.pdf", uuid.String()))
	err = pngsToPDF(cwd, pdfFileName)
	if err != nil {
		log.Printf("Err: %s", err)
		return uuid, err
	}
	
	log.Println("PDF-File generated:", pdfFileName)

	smtpService, err := NewSmtpService(config)

	if err != nil {
		if *debug == true {
			log.Println("DEBUG:send mail to", smtpService.config.Smtp.Recipient)
		}
		err := smtpService.SendMail(pdfFileName)
		if err != nil {
			log.Printf("Err: %s", err)
			return uuid, err
		} else if *debug == true {
			log.Println("DEBUG:mail successfully sent to", smtpService.config.Smtp.Recipient)
		}
	} else if *debug == true {
		log.Println("DEBUG:omit send mail:no smtp configured")
	}

	return uuid, err
}

func mustResolveBinary(bin string) *string {
	path, err := exec.LookPath(bin);
	if err != nil {
		log.Fatalf("cant locate %s. Please install!", bin)
	}
	return  &path
}

func pngsToPDF(cwd string, pdfPath string) error {

	var pngFiles []string
	entries, err := os.ReadDir(cwd)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if filepath.Ext(e.Name()) == fmt.Sprintf(".%s", scanFormat) {
			pngFiles = append(
				pngFiles,
				filepath.Join(cwd, e.Name()),
			)
		}
	}
	sort.Strings(pngFiles)

	pdf := gofpdf.New("P", "pt", "A4", "")
	pdf.SetAutoPageBreak(false, 0)

	for _, file := range pngFiles {
		f, err := os.Open(file)
		if err != nil {
			return err
		}

		img, err := png.DecodeConfig(f)
		f.Close()
		if err != nil {
			return err
		}

		w := float64(img.Width)
		h := float64(img.Height)

		pdf.AddPageFormat("P", gofpdf.SizeType{Wd: w, Ht: h})
		pdf.ImageOptions(
			file,
			0, 0,
			w, h,
			false,
			gofpdf.ImageOptions{ImageType: "PNG"},
			0,
			"",
		)
	}

	return pdf.OutputFileAndClose(pdfPath)
}
