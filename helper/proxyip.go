package helper

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	CS "github.com/LalatinaHub/LatinaServer/constant"
)

type ProxyInfo struct {
	Host string
	Port int
}

type ProxyIPInfo struct {
	Proxy          string `json:"proxy"`
	Port           int    `json:"port"`
	ProxyIP        bool   `json:"proxyip"`
	Delay          int64  `json:"delay"`
	IP             string `json:"ip"`
	Colo           string `json:"colo"`
	Longitude      string `json:"longitude"`
	HTTPProtocol   string `json:"httpProtocol"`
	Continent      string `json:"continent"`
	Asn            int    `json:"asn"`
	Country        string `json:"country"`
	TLSVersion     string `json:"tlsVersion"`
	City           string `json:"city"`
	Timezone       string `json:"timezone"`
	PostalCode     string `json:"postalCode"`
	Region         string `json:"region"`
	Latitude       string `json:"latitude"`
	RegionCode     string `json:"regionCode"`
	AsOrganization string `json:"asOrganization"`
	Message        string `json:"message,omitempty"`
}

func sendProxyRequest(host, path string, proxy *ProxyInfo) (string, error) {
	var address string
	if proxy != nil {
		address = fmt.Sprintf("%s:%d", proxy.Host, proxy.Port)
	} else {
		address = fmt.Sprintf("%s:%d", host, 443)
	}

	conn, err := tls.DialWithDialer(&net.Dialer{Timeout: 5 * time.Second}, "tcp", address, &tls.Config{ServerName: host, InsecureSkipVerify: true})
	if err != nil {
		return "", err
	}
	defer conn.Close()

	request := fmt.Sprintf(
		"GET %s HTTP/1.1\r\nHost: %s\r\nUser-Agent: Mozilla/5.0\r\nConnection: close\r\n\r\n",
		path, host,
	)

	_, err = conn.Write([]byte(request))
	if err != nil {
		return "", err
	}

	response, err := io.ReadAll(conn)
	if err != nil {
		return "", err
	}

	parts := strings.SplitN(string(response), "\r\n\r\n", 2)
	if len(parts) > 1 {
		return parts[1], nil
	}
	return "", nil
}

func CheckProxyIP(proxyIP string) (ProxyIPInfo, error) {
	if proxyIP == "" {
		return ProxyIPInfo{}, nil
	}

	proxyParts := strings.Split(proxyIP, ":")
	if len(proxyParts) != 2 {
		return ProxyIPInfo{Message: "Invalid proxy format"}, fmt.Errorf("invalid proxy format")
	}

	proxyAddress := proxyParts[0]
	proxyPort, _ := strconv.Atoi(proxyParts[1])
	proxy := &ProxyInfo{Host: proxyAddress, Port: proxyPort}

	start := time.Now()
	ipinfo, err1 := sendProxyRequest(CS.IP_RESOLVER_DOMAIN, CS.IP_RESOLVER_PATH, proxy)
	finish := time.Now()

	if err1 != nil {
		return ProxyIPInfo{
			ProxyIP: false,
			Message: fmt.Sprintf("%v", err1),
		}, nil
	}

	proxyIpInfo := ProxyIPInfo{}
	json.Unmarshal([]byte(ipinfo), &proxyIpInfo)

	myip := GetIpInfo()
	if proxyIpInfo.IP != "" && proxyIpInfo.IP != myip.Ip {
		proxyIpInfo.Proxy = proxy.Host
		proxyIpInfo.Port = proxy.Port
		proxyIpInfo.ProxyIP = true
		proxyIpInfo.Delay = finish.Sub(start).Milliseconds()

		return proxyIpInfo, nil
	}

	return ProxyIPInfo{
		ProxyIP: false,
	}, nil
}
