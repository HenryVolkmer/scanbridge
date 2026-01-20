package main

import "encoding/xml"

type ScannerCapabilities struct {
	XMLName xml.Name `xml:"ScannerCapabilities"`

	Version        string `xml:"Version"`
	MakeAndModel   string `xml:"MakeAndModel"`
	SerialNumber   string `xml:"SerialNumber"`

	UUID      string `xml:"UUID"`
	AdminURI  string `xml:"AdminURI"`
	IconURI   string `xml:"IconURI"`

	SettingProfiles SettingProfiles `xml:"SettingProfiles"`

	Platen *Platen `xml:"Platen"`
	Adf *Adf    `xml:"Adf"`

	StoredJobRequestSupport StoredJobRequestSupport `xml:"StoredJobRequestSupport"`

	BlankPageDetection              bool `xml:"BlankPageDetection"`
	BlankPageDetectionAndRemoval    bool `xml:"BlankPageDetectionAndRemoval"`
	OCRSupport                      bool `xml:"OCRSupport"`
	OCRLanguageSupport              OCRLanguageSupport `xml:"OCRLanguageSupport"`
}

/* ---------------- Setting Profiles ---------------- */

type SettingProfiles struct {
	Profiles []SettingProfile `xml:"SettingProfile"`
}

type SettingProfile struct {
	Name string `xml:"name,attr"`
	Ref  string `xml:"ref,attr"`

	ColorModes         []string `xml:"ColorModes>ColorMode"`
	DocumentFormats    DocumentFormats `xml:"DocumentFormats"`
	SupportedResolutions SupportedResolutions `xml:"SupportedResolutions"`

	ColorSpaces     []ColorSpace `xml:"ColorSpaces>ColorSpace"`
	CcdChannels     []CcdChannel `xml:"CcdChannels>CcdChannel"`
	BinaryRenderings []BinaryRendering `xml:"BinaryRenderings>BinaryRendering"`
}

type DocumentFormats struct {
	DocumentFormat    []string `xml:"DocumentFormat"`
	DocumentFormatExt []string `xml:"DocumentFormatExt"`
}

/* ---------------- Resolutions ---------------- */

type SupportedResolutions struct {
	DiscreteResolutions []DiscreteResolution `xml:"DiscreteResolutions>DiscreteResolution"`
	ResolutionRange     *ResolutionRange     `xml:"ResolutionRange"`
}

type DiscreteResolution struct {
	XResolution int `xml:"XResolution"`
	YResolution int `xml:"YResolution"`
}

type ResolutionRange struct {
	XResolutionRange ResolutionAxis `xml:"XResolutionRange"`
	YResolutionRange ResolutionAxis `xml:"YResolutionRange"`
}

type ResolutionAxis struct {
	Min    int `xml:"Min"`
	Max    int `xml:"Max"`
	Normal int `xml:"Normal"`
	Step   int `xml:"Step"`
}

/* ---------------- Color / Channels ---------------- */

type ColorSpace struct {
	Default bool   `xml:"default,attr"`
	Value   string `xml:",chardata"`
}

type CcdChannel struct {
	Default bool   `xml:"default,attr"`
	Value   string `xml:",chardata"`
}

type BinaryRendering struct {
	Default bool   `xml:"default,attr"`
	Value   string `xml:",chardata"`
}

/* ---------------- Platen ---------------- */

type Platen struct {
	InputCaps PlatenInputCaps `xml:"PlatenInputCaps"`
}

type PlatenInputCaps struct {
	MinWidth        int `xml:"MinWidth"`
	MaxWidth        int `xml:"MaxWidth"`
	MinHeight       int `xml:"MinHeight"`
	MaxHeight       int `xml:"MaxHeight"`
	MaxScanRegions  int `xml:"MaxScanRegions"`

	SettingProfiles SettingProfiles `xml:"SettingProfiles"`

	SupportedResolutions SupportedResolutions `xml:"SupportedResolutions"`

	MaxOpticalXResolution int `xml:"MaxOpticalXResolution"`
	MaxOpticalYResolution int `xml:"MaxOpticalYResolution"`

	RiskyLeftMargin   int `xml:"RiskyLeftMargin"`
	RiskyRightMargin  int `xml:"RiskyRightMargin"`
	RiskyTopMargin    int `xml:"RiskyTopMargin"`
	RiskyBottomMargin int `xml:"RiskyBottomMargin"`
}

/* ---------------- ADF ---------------- */

type Adf struct {
	SimplexInputCaps AdfSimplexInputCaps `xml:"AdfSimplexInputCaps"`
	FeederCapacity   int `xml:"FeederCapacity"`
	AdfOptions       []string `xml:"AdfOptions>AdfOption"`
}

type AdfSimplexInputCaps struct {
	MinWidth        int `xml:"MinWidth"`
	MaxWidth        int `xml:"MaxWidth"`
	MinHeight       int `xml:"MinHeight"`
	MaxHeight       int `xml:"MaxHeight"`

	SettingProfiles SettingProfiles `xml:"SettingProfiles"`

	EdgeAutoDetection EdgeAutoDetection `xml:"EdgeAutoDetection"`

	MaxOpticalXResolution int `xml:"MaxOpticalXResolution"`
	MaxOpticalYResolution int `xml:"MaxOpticalYResolution"`

	RiskyLeftMargin   int `xml:"RiskyLeftMargin"`
	RiskyRightMargin  int `xml:"RiskyRightMargin"`
	RiskyTopMargin    int `xml:"RiskyTopMargin"`
	RiskyBottomMargin int `xml:"RiskyBottomMargin"`
}

type EdgeAutoDetection struct {
	SupportedEdges []string `xml:"SupportedEdge"`
}

/* ---------------- Jobs / OCR ---------------- */

type StoredJobRequestSupport struct {
	MaxStoredJobRequests int `xml:"MaxStoredJobRequests"`
	TimeoutInSeconds     int `xml:"TimeoutInSeconds"`
}

type OCRLanguageSupport struct {
	Languages []string `xml:"NaturalLanguageSupported"`
}
