package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"

	"github.com/vishvananda/netlink"
	"golang.zx2c4.com/wireguard/wgctrl"
)

func toggleInterface(name string) (string, error) {
	link, err := netlink.LinkByName(name)
	if err != nil {
		return "", err
	}
	if link.Attrs().Flags&net.FlagUp != 0 {
		netlink.LinkSetDown(link)
		return fmt.Sprintf("Interface %s brought down", name), nil
	} else {
		netlink.LinkSetUp(link)
		return fmt.Sprintf("Interface %s brought up", name), nil
	}
}

func listInterfaces() ([]string, error) {
	client, err := wgctrl.New()
	if err != nil {
		return nil, err
	}
	defer client.Close()
	devs, err := client.Devices()
	if err != nil {
		return nil, err
	}
	names := make([]string, len(devs))
	for i, d := range devs {
		names[i] = d.Name
	}
	return names, nil
}

func main() {
	addr := "127.0.0.1:9999"
	if a := os.Getenv("WG_HELPER_ADDR"); a != "" {
		addr = a
	}

	udpAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	log.Println("WireGuard UDP helper listening on", addr)

	buf := make([]byte, 1024)
	for {
		n, clientAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Println("read error:", err)
			continue
		}

		command := strings.TrimSpace(string(buf[:n]))
		parts := strings.Fields(command)
		if len(parts) == 0 {
			conn.WriteToUDP([]byte("ERROR: empty command"), clientAddr)
			continue
		}

		switch parts[0] {
		case "toggle":
			if len(parts) < 2 {
				conn.WriteToUDP([]byte("ERROR: missing device"), clientAddr)
				continue
			}
			msg, err := toggleInterface(parts[1])
			if err != nil {
				conn.WriteToUDP([]byte("ERROR: "+err.Error()), clientAddr)
			} else {
				conn.WriteToUDP([]byte("OK: "+msg), clientAddr)
			}

		case "list":
			devs, err := listInterfaces()
			if err != nil {
				conn.WriteToUDP([]byte("ERROR: "+err.Error()), clientAddr)
			} else {
				conn.WriteToUDP([]byte("OK: "+strings.Join(devs, " ")), clientAddr)
			}

		default:
			conn.WriteToUDP([]byte("ERROR: unknown command "+parts[0]), clientAddr)
		}
	}
}
