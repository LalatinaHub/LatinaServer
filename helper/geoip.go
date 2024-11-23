package helper

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/LalatinaHub/LatinaSub-go/geoip"
)

var (
	ipinfo = geoip.GeoIpJson{}
)

func GetIpInfo() geoip.GeoIpJson {
	if ipinfo.Ip != "" {
		return ipinfo
	}

	buf := new(strings.Builder)
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := httpClient.Get("http://ipinfo.io/ip")
	if err != nil {
		return ipinfo
	}
	defer resp.Body.Close()

	io.Copy(buf, resp.Body)
	if resp.StatusCode == 200 {
		myIp := geoip.MyIp{}
		if err := json.Unmarshal([]byte(buf.String()), &myIp); err == nil {
			ipinfo = geoip.Parse(myIp)
		}
	}

	return ipinfo
}
