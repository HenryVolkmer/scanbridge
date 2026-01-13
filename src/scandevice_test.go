package main

import (
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