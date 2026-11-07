package dlna

import (
	"strings"

	"github.com/oxzjh/ssdp"
)

type DLNA struct {
	avtransportSCPDURL    string
	avtransportControlURL string
	renderingControlURL   string
}

func (d *DLNA) GetAllowSpeeds() ([]string, error) {
	resp, err := ssdp.GetSCPD(d.avtransportSCPDURL)
	if err != nil {
		return nil, err
	}
	for _, v := range resp.StateVariable {
		if v.Name == "TransportPlaySpeed" {
			speeds := make([]string, 0, len(v.AllowedValues))
			for _, a := range v.AllowedValues {
				s := strings.TrimSpace(a)
				if s != "vendor-defined" {
					speeds = append(speeds, s)
				}
			}
			return speeds, nil
		}
	}
	return []string{"1"}, nil
}

func (d *DLNA) GetTransportInfo() (*GetTransportInfoResponse, error) {
	var response GetTransportInfoResponse
	if err := ssdp.SoapPost(d.avtransportControlURL, "GetTransportInfo", "AVTransport", instanceID, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (d *DLNA) GetPositionInfo() (*GetPositionInfoResponse, error) {
	var response GetPositionInfoResponse
	if err := ssdp.SoapPost(d.avtransportControlURL, "GetPositionInfo", "AVTransport", instanceID, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (d *DLNA) SetAVTransportURI(url string) error {
	return ssdp.SoapPost(d.avtransportControlURL, "SetAVTransportURI", "AVTransport", buildSetAVTransportURI(url), nil)
}

func (d *DLNA) Play(speed string) error {
	return ssdp.SoapPost(d.avtransportControlURL, "Play", "AVTransport", buildPlay(speed), nil)
}

func (d *DLNA) Pause() error {
	return ssdp.SoapPost(d.avtransportControlURL, "Pause", "AVTransport", instanceID, nil)
}

func (d *DLNA) Stop() error {
	return ssdp.SoapPost(d.avtransportControlURL, "Stop", "AVTransport", instanceID, nil)
}

// reltime: 00:00:00
func (d *DLNA) Seek(reltime string) error {
	return ssdp.SoapPost(d.avtransportControlURL, "Seek", "AVTransport", buildSeek(reltime), nil)
}

func (d *DLNA) GetMute() (*GetMuteResponse, error) {
	var response GetMuteResponse
	if err := ssdp.SoapPost(d.renderingControlURL, "GetMute", "RenderingControl", channel, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (d *DLNA) SetMute(mute bool) error {
	return ssdp.SoapPost(d.renderingControlURL, "SetMute", "RenderingControl", buildSetMute(mute), nil)
}

func (d *DLNA) GetVolume() (*GetVolumeResponse, error) {
	var response GetVolumeResponse
	if err := ssdp.SoapPost(d.renderingControlURL, "GetVolume", "RenderingControl", channel, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

// volume: 0-100
func (d *DLNA) SetVolume(volume int) error {
	return ssdp.SoapPost(d.renderingControlURL, "SetVolume", "RenderingControl", buildSetVolume(volume), nil)
}

func New(device *ssdp.Device) *DLNA {
	d := &DLNA{}
	for _, service := range device.Services {
		switch service.ServiceId {
		case "urn:upnp-org:serviceId:AVTransport":
			d.avtransportSCPDURL = device.GetURL(service.SCPDURL)
			d.avtransportControlURL = device.GetURL(service.ControlURL)
		case "urn:upnp-org:serviceId:RenderingControl":
			d.renderingControlURL = device.GetURL(service.ControlURL)
		}
	}
	return d
}
