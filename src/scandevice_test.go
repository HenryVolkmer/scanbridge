package main

import (
	"path"
	"net"
	"net/http"
	"net/url"
	"net/http/httptest"
	"os"
	"testing"
)

func TestCanCreateConfigByDto(t *testing.T) {
	dev := &ScanDevice{
		Version: "2.0",
		Cs: []string{"grey"},
		Is: []string{"adf"},
		Pdl: []string{"application/pdf"},
	}

	if !dev.isPdfSupported() {
		t.Fatalf("PDF is not supported MIME-Type, but should be!")
	}

	err := dev.NewScanJob(&ScanSettingsDto{
		Version: "2.0",
		DocumentFormat: "pdf",
		ColorMode: "grey",
		InputSource: "adf",
		XResolution: 300,
		YResolution: 300,
		Height: 3300,
		Width: 2550,
		XOffset: 0,
		YOffset: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCantCreateConfigByInvalidDto(t *testing.T) {
	dev := &ScanDevice{
		Version: "2.0",
		Cs: []string{"color"},
		Is: []string{"platen"},
		Pdl: []string{"image/jpeg"},
	}
	err := dev.NewScanJob(&ScanSettingsDto{
		Version: "2.0",
		DocumentFormat: "pdf",
		ColorMode: "grey",
		InputSource: "adf",
	})
	if err == nil {
		t.Fatal(err)
	}
}

func TestCanFetchCapabilitiesAndCreateDevice(t *testing.T) {

	scannerCapabilitiesXML, err := os.ReadFile(path.Join("testdata/caps.xml"))
	if err != nil {
		t.Fatal(err)
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/eSCL/ScannerCapabilities" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write(scannerCapabilitiesXML)
	}))
	defer server.Close()

	// Server-URL -> IP extrahieren
	u, err := url.Parse(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	host, _, _ := net.SplitHostPort(u.Host)
	ip := net.ParseIP(host)

	dev, err := NewScanDevice(server.Client(), &ip)
	if err != nil {
		t.Fatal(err)
	}

	if dev == nil {
		t.Fatal("device is nil")
	}
}
