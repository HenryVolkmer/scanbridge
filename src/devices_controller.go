package main

import (
	"context"
	"time"
	"log"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"github.com/grandcat/zeroconf"
	"github.com/joho/godotenv"
)

type DevicesController struct {
	config *Config
	Devices []*ScanDevice `json:"devices"`
}

func (dc *DevicesController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    var err error
    var dev []*ScanDevice
	dec := json.NewEncoder(w)
    if dc.config.IsAutodiscovery == true {
    	dev, err = dc.discover()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Err: %s", err)
			dec.Encode(&Notification{Data: "cant fetch devices!", Title: "KO!"})
			return
		}
    } else {
    	dev = dc.config.Devices
    }
	dec.Encode(dev)
}

func (dc *DevicesController) discover() ([]*ScanDevice, error) {

	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, err
	}
	
	var buffer []*ScanDevice = []*ScanDevice{}
	entries := make(chan *zeroconf.ServiceEntry)
	go func(results <-chan *zeroconf.ServiceEntry) {
		for entry := range results {
			device, err := dc.newScanDevice(entry)
			if err != nil {
				log.Println(err)
				continue
			}
			buffer = append(buffer, device)
		}
	}(entries)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(4))
	defer cancel()
	err = resolver.Browse(ctx, "_uscan._tcp", "local", entries)
	if err != nil {
		return nil, err
	}

	<-ctx.Done()

	return buffer, nil
}

func (dc *DevicesController) newScanDevice(se *zeroconf.ServiceEntry) (*ScanDevice, error) {

	eSCL := strings.Replace(strings.Join(se.Text, "\n"), "-", "", -1)
	eSCLCapabilitiesMap, err := godotenv.Unmarshal(eSCL)

	if err != nil {
		log.Println("godotenv:",err)
		return nil, err
	}

	dev := &ScanDevice{AddrIPv4: se.AddrIPv4[0]}

	if ty, ok := eSCLCapabilitiesMap["ty"]; ok {
		dev.Ty = ty
	}

	if rep, ok := eSCLCapabilitiesMap["representation"]; ok {
 		url, err := url.Parse(rep)
 		if err == nil {
 			url.Host = se.AddrIPv4[0].String()
			dev.Representation = url.String()
		}
	}

	if cs, ok := eSCLCapabilitiesMap["cs"]; ok {
		dev.Cs = strings.Split(cs, ",")
	}
	
	if is, ok := eSCLCapabilitiesMap["is"]; ok {
		dev.Is = strings.Split(is, ",")
	}

	return  dev, nil
}

func NewDevicesController(c *Config) *DevicesController {
	return  &DevicesController{config: c}
}