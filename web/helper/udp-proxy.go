package web_helper

import (
	"fmt"
	"net"
)

// Proxy Request struct
type ProxyRequest struct {
	Address string `json:"address"`
	Port    uint16 `json:"port"`
	Payload []byte `json:"payload"`
}

// UDP Proxy (UDP ASSOCIATE)
func HandleUDPForwarding(req ProxyRequest) []byte {
	serverAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", req.Address, req.Port))
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return []byte{}
	}

	// Membuat koneksi UDP
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		fmt.Println("Failed to connect to UDP target:", err)
		return []byte{}
	}
	defer conn.Close()

	// Kirim payload
	_, err = conn.Write(req.Payload)
	if err != nil {
		fmt.Println("Failed to send UDP data:", err)
		return []byte{}
	}

	// Membaca respons dari server UDP
	buffer := make([]byte, 4096)
	n, _, err := conn.ReadFrom(buffer)
	if err != nil {
		fmt.Println("Failed to read from UDP target:", err)
		return []byte{}
	}

	fmt.Println("UDP response received:", string(buffer[:n]))

	return buffer
}
