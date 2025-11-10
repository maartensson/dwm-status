package main

import (
	"fmt"
	"net"
	"os"
)

func sendCommand(addr string, cmd string) (string, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return "", fmt.Errorf("resolve error: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return "", fmt.Errorf("dial error: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(cmd))
	if err != nil {
		return "", fmt.Errorf("write error: %w", err)
	}

	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return "", fmt.Errorf("read error: %w", err)
	}

	return string(buf[:n]), nil
}

func toggleWG(addr string, device string) (string, error) {
	if device == "" {
		return "", fmt.Errorf("device name required")
	}
	return sendCommand(addr, "toggle "+device)
}

func listWG(addr string) (string, error) {
	return sendCommand(addr, "list")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: wgctl <toggle|list> [device]")
		os.Exit(1)
	}

	addr := "127.0.0.1:9999"
	cmd := os.Args[1]

	var resp string
	var err error

	switch cmd {
	case "toggle":
		if len(os.Args) < 3 {
			fmt.Println("usage: wgctl toggle <device>")
			os.Exit(1)
		}
		device := os.Args[2]
		resp, err = toggleWG(addr, device)

	case "list":
		resp, err = listWG(addr)

	default:
		fmt.Println("unknown command:", cmd)
		os.Exit(1)
	}

	if err != nil {
		fmt.Println("error:", err)
		os.Exit(1)
	}

	fmt.Println(resp)
}
