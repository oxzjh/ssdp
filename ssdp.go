package ssdp

import (
	"bufio"
	"bytes"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	ssdpAddr = "239.255.255.250:1900"

	SearchTargetAll             = "ssdp:all"
	SearchTargetRootDevice      = "upnp:rootdevice"
	SearchTargetMediaRenderer   = "urn:schemas-upnp-org:service:MediaRenderer:1"
	SearchTargetAVTransport     = "urn:schemas-upnp-org:service:AVTransport:1"
	SearchTargetWANIPConnection = "urn:schemas-upnp-org:service:WANIPConnection:1"
)

var (
	headerPrefix = []byte("HTTP")
	headerSuffix = []byte("\r\n\r\n")
)

type Target struct {
	Header   http.Header
	Location string
}

func Search(timeout time.Duration, searchTarget string) ([]*Target, error) {
	addr, err := net.ResolveUDPAddr("udp4", ssdpAddr)
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenPacket("udp4", ":0")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var b strings.Builder
	b.WriteString("M-SEARCH * HTTP/1.1\r\n")
	b.WriteString("HOST: " + ssdpAddr + "\r\n")
	b.WriteString("MAN: \"ssdp:discover\"\r\n")
	b.WriteString("MX: " + strconv.Itoa(int(timeout/time.Second)) + "\r\n")
	b.WriteString("ST: " + searchTarget + "\r\n")
	b.WriteString("\r\n")

	if _, err = conn.WriteTo([]byte(b.String()), addr); err != nil {
		return nil, err
	}

	if err = conn.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}

	var targets []*Target
	buf := make([]byte, 1<<16)
	for {
		n, _, err := conn.ReadFrom(buf)
		if err != nil {
			break
		}
		data := buf[:n]
		if bytes.HasPrefix(data, headerPrefix) && bytes.HasSuffix(data, headerSuffix) {
			if resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(data)), nil); err == nil {
				targets = append(targets, &Target{
					Header:   resp.Header,
					Location: resp.Header.Get("Location"),
				})
			}
		}
	}
	return targets, nil
}
