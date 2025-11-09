package main

import (
	"encoding/json"
	"log"
	"net"
	"os"

	"golang.zx2c4.com/wireguard/wgctrl"
)

type Response struct {
	Devices []string `json:"devices"`
	Error   string   `json:"error,omitempty"`
}

func getWireGuardDevices() ([]string, error) {
	client, err := wgctrl.New()
	if err != nil {
		return nil, err
	}
	defer client.Close()
	devices, err := client.Devices()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(devices))
	for _, d := range devices {
		names = append(names, d.Name)
	}
	return names, nil
}

func main() {
	socketPath := os.Getenv("SOCKET_PATH")
	if socketPath == "" {
		socketPath = "/run/wg-helper/wg-helper.sock"
	}

	// Remove existing socket
	os.Remove(socketPath)

	l, err := net.Listen("unix", socketPath)
	if err != nil {
		log.Fatal("listen error:", err)
	}
	defer l.Close()
	os.Chmod(socketPath, 0666) // allow user access via group if needed

	log.Println("WireGuard helper listening on", socketPath)

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}

		go func(c net.Conn) {
			defer c.Close()

			devices, err := getWireGuardDevices()
			resp := Response{Devices: devices}
			if err != nil {
				resp.Error = err.Error()
			}

			enc := json.NewEncoder(c)
			enc.Encode(resp)
		}(conn)
	}
}
