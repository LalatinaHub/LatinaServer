package helper

import (
	"crypto/tls"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/LalatinaHub/LatinaSub-go/ipapi"
)

var (
	ipinfo = ipapi.Ipapi{}
)

func GetIpInfo() ipapi.Ipapi {
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

	resp, err := httpClient.Get("http://ipinfo.io/json")
	if err != nil {
		return ipinfo
	}
	defer resp.Body.Close()

	io.Copy(buf, resp.Body)
	if resp.StatusCode == 200 {
		ipinfo = ipapi.Parse(buf.String())
		return ipinfo
	}

	return ipinfo
}
