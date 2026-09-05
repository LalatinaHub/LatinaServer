package geoip

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/LalatinaHub/LatinaServer/pkg/constant"
)

type MyIp struct {
	Ip  string `json:"ip,omitempty"`
	CC  string `json:"country,omitempty"`
	Org string `json:"asOrganization,omitempty"`
}

type Countries struct {
	Name   string `json:"name"`
	Code   string `json:"code"`
	Region string `json:"region"`
}

type GeoIpJson struct {
	Ip          string `json:"ip,omitempty"`
	CountryName string `json:"country_name,omitempty"`
	CountryCode string `json:"country,omitempty"`
	Region      string `json:"region,omitempty"`
	Org         string `json:"org,omitempty"`
}

var (
	symbolRegex = regexp.MustCompile("[^a-zA-Z0-9 ]")
	ipinfo      = GeoIpJson{}
)

func Parse(myIp MyIp) GeoIpJson {
	result := GeoIpJson{
		Ip:          myIp.Ip,
		CountryName: "Unknown",
		CountryCode: "XX",
		Region:      "Unknown",
		Org:         "LalatinaHub",
	}

	for _, country := range CountryList {
		if country.Code == myIp.CC {
			result.CountryName = country.Name
			result.CountryCode = country.Code
			result.Region = country.Region
			result.Org = symbolRegex.ReplaceAllString(myIp.Org, "")
			return result
		}
	}

	return result
}

func GetIpInfo() GeoIpJson {
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

	resp, err := httpClient.Get("http://" + constant.IPResolverDomain + constant.IPResolverPath)
	if err != nil {
		return ipinfo
	}
	defer resp.Body.Close()

	io.Copy(buf, resp.Body)
	if resp.StatusCode == 200 {
		myIp := MyIp{}
		if err := json.Unmarshal([]byte(buf.String()), &myIp); err == nil {
			ipinfo = Parse(myIp)
		}
	}

	return ipinfo
}

