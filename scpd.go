package ssdp

import (
	"encoding/xml"
	"net/http"
)

type SCPDResponse struct {
	Actions       []*Action        `xml:"actionList>action"`
	StateVariable []*StateVariable `xml:"serviceStateTable>stateVariable"`
}

type Action struct {
	Name      string      `xml:"name"`
	Arguments []*Argument `xml:"argumentList>argument"`
}

type Argument struct {
	Name                 string `xml:"name"`
	Direction            string `xml:"direction"`
	RelatedStateVariable string `xml:"relatedStateVariable"`
}

type StateVariable struct {
	Name          string   `xml:"name"`
	DataType      string   `xml:"dataType"`
	AllowedValues []string `xml:"allowedValueList>allowedValue"`
}

func GetSCPD(url string) (*SCPDResponse, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var response SCPDResponse
	if err := xml.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	return &response, nil
}
