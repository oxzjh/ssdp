package ssdp

import (
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func SoapPost(url, action, service, content string, response any) error {
	body := `<?xml version="1.0" encoding="utf-8"?>
<s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/">
<s:Body>
` + fmt.Sprintf(`<u:%s xmlns:u="urn:schemas-upnp-org:service:%s:1">%s</u:%s>`, action, service, content, action) + `
</s:Body>
</s:Envelope>`
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", `text/xml; charset="utf-8"`)
	req.Header.Set("SOAPAction", fmt.Sprintf(`"urn:schemas-upnp-org:service:%s:1#%s"`, service, action))
	req.Header.Set("Connection", "close")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		decoder := xml.NewDecoder(resp.Body)
		found := false
		for {
			t, err := decoder.Token()
			if err != nil {
				return err
			}
			switch token := t.(type) {
			case xml.StartElement:
				if token.Name.Local == "errorDescription" {
					found = true
				}
			case xml.EndElement:
				found = false
			case xml.CharData:
				if found {
					return errors.New(string(token))
				}
			}
		}
	} else if response != nil {
		return xml.NewDecoder(resp.Body).Decode(response)
	}
	return nil
}
