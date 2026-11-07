package dlna

import (
	"encoding/xml"
	"strconv"
)

const (
	instanceID = "<InstanceID>0</InstanceID>"
	channel    = instanceID + "<Channel>Master</Channel>"
)

type GetTransportInfoResponse struct {
	XMLName               xml.Name `xml:"Envelope"`
	CurrentTransportState string   `xml:"Body>GetTransportInfoResponse>CurrentTransportState"`
}

type GetPositionInfoResponse struct {
	XMLName       xml.Name `xml:"Envelope"`
	TrackDuration string   `xml:"Body>GetPositionInfoResponse>TrackDuration"`
	RelTime       string   `xml:"Body>GetPositionInfoResponse>RelTime"`
}

type GetMuteResponse struct {
	XMLName     xml.Name `xml:"Envelope"`
	CurrentMute bool     `xml:"Body>GetMuteResponse>CurrentMute"`
}

type GetVolumeResponse struct {
	XMLName       xml.Name `xml:"Envelope"`
	CurrentVolume int      `xml:"Body>GetVolumeResponse>CurrentVolume"`
}

func buildSetAVTransportURI(url string) string {
	return instanceID + "<CurrentURI>" + url + "</CurrentURI><CurrentURIMetaData/>"
}

func buildPlay(speed string) string {
	return instanceID + "<Speed>" + speed + "</Speed>"
}

func buildSeek(reltime string) string {
	return instanceID + "<Unit>REL_TIME</Unit><Target>" + reltime + "</Target>"
}

func buildSetMute(mute bool) string {
	m := "0"
	if mute {
		m = "1"
	}
	return channel + "<DesiredMute>" + m + "</DesiredMute>"
}

func buildSetVolume(volume int) string {
	return channel + "<DesiredVolume>" + strconv.Itoa(volume) + "</DesiredVolume>"
}
