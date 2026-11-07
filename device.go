package ssdp

import (
	"encoding/xml"
	"net"
	"net/http"
	"net/url"
	"strings"
)

type root struct {
	Device *Device `xml:"device"`
}

type Device struct {
	FriendlyName string     `xml:"friendlyName"`
	Services     []*Service `xml:"serviceList>service"`
	LocalIP      string
	url          string
}

type Service struct {
	ServiceType string `xml:"serviceType"`
	ServiceId   string `xml:"serviceId"`
	SCPDURL     string `xml:"SCPDURL"`
	ControlURL  string `xml:"controlURL"`
	EventSubURL string `xml:"eventSubURL"`
}

func (d *Device) GetURL(u string) string {
	if strings.HasPrefix(u, "/") {
		return d.url + u
	}
	return d.url + "/" + u
}

func NewDevice(location string) (*Device, error) {
	u, err := url.Parse(location)
	if err != nil {
		return nil, err
	}
	conn, err := net.Dial("udp4", u.Host)
	if err != nil {
		panic(err)
	}
	localIP := strings.Split(conn.LocalAddr().String(), ":")[0]
	conn.Close()

	req, err := http.NewRequest(http.MethodGet, location, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Connection", "close")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	var r root
	if err = xml.NewDecoder(resp.Body).Decode(&r); err != nil {
		return nil, err
	}
	resp.Body.Close()
	r.Device.url = u.Scheme + "://" + u.Host
	r.Device.LocalIP = localIP
	return r.Device, nil
}
