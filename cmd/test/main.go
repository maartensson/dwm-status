package main

import (
	"fmt"
	"net"
	"time"
)

func getActiveIface() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:53")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localIP := conn.LocalAddr().(*net.UDPAddr).IP

	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.Equal(localIP) {
				return iface.Name, nil
			}
		}
	}
	return "", fmt.Errorf("interface not found for IP %s", localIP)
}

func main() {
	var lastIface string
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		iface, err := getActiveIface()
		if err != nil {
			fmt.Println("No internet")
			continue
		}
		if iface != lastIface {
			fmt.Println("Active interface:", iface)
			lastIface = iface
		}
	}
}
