package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net"
	"slices"
	"strings"
)

// ScanDevice is modeled against the 
// Mopria Alliance eSCL Technical Specification v2.97
// The eSCL Spec introduces the "Cs", "Is", "Pdl" ... 
// Props
type ScanDevice struct {
	// Host machine IPv4 address
	AddrIPv4 net.IP `json:"IPv4"`
	// human-readable make and model
	Ty string `json:"name"` 
	// URL to a PNG or ICO file containing a graphical
	// representation of the scanner.
	Representation string `json:"representation"`
	// The ColorSpace defines the color capabilities of the scanner:
	// "color" if the Scanner supports color scanning, "grayscale" if
	// the scanner supports grayscale, "binary" if the scanner
	// supports 1-bit monochrome scanning.
	Cs []string `json:"color_spaces"`
	// The InputSource defines the list of scan input options:
	// "platen" for glass flat bed scanning, "adf" for Automatic
	// Document Feeder, "camera" if the Scanner has a non-
	// traditional scan bed (such as a stage)
	Is []string `json:"input_sources"`

	Version string `json:"version"`

	// List of MIME media types supported by the scanner
	// application/pdf,image/jpeg
	Pdl []string
}

// NewScanJob advices the Scanner to enqueue a new Scan-Job
func (sd *ScanDevice) NewScanJob(dto *ScanSettingsDto) error {
	
	if err := sd.Validate(dto); err != nil {
		return  err
	}

	settings := scanSettings{
		XmlnsPwg:  "http://www.pwg.org/schemas/2010/12/sm",
		XmlnsScan: "http://schemas.hp.com/imaging/escl/2011/05/03",
		Version: dto.Version,
		ScanRegions: scanRegions{
			ScanRegion: scanRegion{
				ContentRegionUnits: "escl:ThreeHundredthsOfInches",
				Height: dto.Height,
				Width: dto.Width,
				XOffset: dto.XOffset,
				YOffset: dto.YOffset,
			},
		},
		ColorMode: dto.ColorMode,
		XResolution: dto.XResolution,
		YResolution: dto.YResolution,
		InputSource: dto.InputSource,
		DocumentFormatExt: &documentFormatExt{
			DocumentFormat: dto.DocumentFormat,
		},
	}

	buf, _ := xml.MarshalIndent(settings, "", "  ")
	log.Println(xml.Header, string(buf))
	return nil
}

// internal indicator, whether Scanner supports PDF generation
func (sd *ScanDevice) isPdfSupported() bool {
	return slices.Contains(sd.Pdl, "application/pdf")
}

// Validate validates the dto against ScanDevice configuration
func (sd *ScanDevice) Validate(dto *ScanSettingsDto) error {

	if !slices.Contains(sd.Cs, dto.ColorMode) {
		return fmt.Errorf(
			"unsupported ColorMode %s, supported are: %s", 
			dto.ColorMode, 
			strings.Join(sd.Cs, ","),
		)
	}

	if !slices.Contains(sd.Is, dto.InputSource) {
		return fmt.Errorf(
			"unsupported InputSource %s, supported are: %s", 
			dto.InputSource, 
			strings.Join(sd.Is, ","),
		)
	}

	if dto.Version != sd.Version {
		return  fmt.Errorf("Given Version %s dont matched support Version %s", dto.Version, sd.Version)
	}

	return  nil
}

// ScanSettingsDto is a shorthand sibling to scanSettings.
// Why not use scanSettings directly? Due the XML-Structure,
// its bloated. We only need the 
// Core-Settings (Version, ColorMode, Resulutions...) in order
// to create a full functional scanSettings-Struct.
type ScanSettingsDto struct {
	Version string
	DocumentFormat string
	ColorMode  string
	InputSource string
	XResolution int
	YResolution int
	// ScanRegion Settings
	Height int
	Width int
	XOffset int
	YOffset int
}

type scanSettings struct {
	XMLName xml.Name `xml:"scan:ScanSettings"`
	XmlnsPwg  string `xml:"xmlns:pwg,attr"`
	XmlnsScan string `xml:"xmlns:scan,attr"`
	Version string `xml:"pwg:Version"`
	ScanRegions scanRegions `xml:"pwg:ScanRegions"`
	ColorMode  string `xml:"scan:ColorMode"`
	XResolution int `xml:"scan:XResolution"`
	YResolution int `xml:"scan:YResolution"`
	InputSource string `xml:"pwg:InputSource"`
	DocumentFormatExt *documentFormatExt `xml:"scan:DocumentFormatExt,omitempty"`
	CompressionFactor *int `xml:"scan:CompressionFactor,omitempty"`
}

type scanRegions struct {
	ScanRegion scanRegion `xml:"pwg:ScanRegion"`
}

type scanRegion struct {
	ContentRegionUnits string `xml:"pwg:ContentRegionUnits"`
	Height             int    `xml:"pwg:Height"`
	Width              int    `xml:"pwg:Width"`
	XOffset            int    `xml:"pwg:XOffset"`
	YOffset            int    `xml:"pwg:YOffset"`
}

type documentFormatExt struct {
	DocumentFormat string `xml:"scan:DocumentFormat"`
}
