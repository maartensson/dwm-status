package vpnmodule

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/mamaart/statusbar/pkg/geoip"
)

type NetStats struct {
	RxBytes uint64
	TxBytes uint64
}

type NetModule struct {
	output       chan []byte
	cachedIP     string
	ipLastUpdate time.Time
	lastIface    string
}

func New() *NetModule {
	n := &NetModule{
		output: make(chan []byte),
	}
	go n.run()
	return n
}

func (n *NetModule) Reader() <-chan []byte {
	return n.output
}

func (n *NetModule) run() {

	iface, localIp, hasInternet := n.getActiveIface()
	if !hasInternet || iface == "" {
		n.output <- []byte(" No internet ")
	} else {
		if iface != n.lastIface || time.Since(n.ipLastUpdate) > time.Hour {
			ip, err := getPublicIP(localIp)
			if err != nil {
				fmt.Fprintln(os.Stderr, "getPublicIP:", err)
			} else {
				n.cachedIP = ip
				n.ipLastUpdate = time.Now()
			}
			n.lastIface = iface
		}
		n.output <- []byte(" " + n.buildStatus(iface) + " ")
	}
	for range time.NewTicker(2 * time.Second).C {
		iface, localIp, hasInternet := n.getActiveIface()
		if !hasInternet || iface == "" {
			n.output <- []byte(" No internet ")
			continue
		}

		if iface != n.lastIface || time.Since(n.ipLastUpdate) > time.Hour {
			ip, err := getPublicIP(localIp)
			if err != nil {
				fmt.Fprintln(os.Stderr, "getPublicIP:", err)
			} else {
				n.cachedIP = ip
				n.ipLastUpdate = time.Now()
			}
			n.lastIface = iface
		}

		n.output <- []byte(" " + n.buildStatus(iface) + " ")
	}
}

func (n *NetModule) getActiveIface() (ifaceName string, localIP net.IP, hasInternet bool) {
	conn, err := net.DialTimeout("udp", "8.8.8.8:53", 1*time.Second)
	if err != nil {
		fmt.Fprintln(os.Stderr, "UDP dial error:", err)
		return "", nil, false
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	localIP = localAddr.IP

	ifaces, err := net.Interfaces()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Interfaces error:", err)
		return "", nil, false
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
			if ip.Equal(localAddr.IP) {
				ifaceName = iface.Name
				break
			}
		}
	}

	return ifaceName, localIP, ifaceName != ""
}

func (n *NetModule) buildStatus(iface string) string {
	stats, err := readTotalBytes(iface)
	if err != nil {
		fmt.Fprintln(os.Stderr, "readTotalBytes error:", err)
		return fmt.Sprintf(" %s", iface)
	}

	vpnActive := isVPNActive(iface)
	flag := getCountry(n.cachedIP)

	emoji := ""
	if vpnActive {
		emoji = "󰒃"
	}

	return fmt.Sprintf("%s %s (%s)  %s  %s", emoji, flag, iface,
		humanizeBytes(stats.RxBytes), humanizeBytes(stats.TxBytes))
}

func humanizeBytes(b uint64) string {
	const unit = 1024
	val := float64(b)
	suffix := "B"
	if b >= unit*unit*unit {
		val /= unit * unit * unit
		suffix = "GiB"
	} else if b >= unit*unit {
		val /= unit * unit
		suffix = "MiB"
	} else if b >= unit {
		val /= unit
		suffix = "KiB"
	}
	return fmt.Sprintf("%.1f%s", val, suffix)
}

func getCountry(ipStr string) string {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return "??"
	}
	geo, err := geoip.New()
	if err != nil {
		fmt.Fprintln(os.Stderr, "geoip.New error:", err)
		return "??"
	}
	c, err := geo.GetCountry(ip)
	if err != nil {
		fmt.Fprintln(os.Stderr, "geo.GetCountry error:", err)
		return "??"
	}
	return c.ISOCode
}

func readTotalBytes(iface string) (NetStats, error) {
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return NetStats{}, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if !strings.HasPrefix(line, iface+":") {
			continue
		}
		parts := strings.Split(line, ":")
		fields := strings.Fields(parts[1])
		if len(fields) < 16 {
			continue
		}
		var rx, tx uint64
		fmt.Sscanf(fields[0], "%d", &rx)
		fmt.Sscanf(fields[8], "%d", &tx)
		return NetStats{RxBytes: rx, TxBytes: tx}, nil
	}
	return NetStats{}, fmt.Errorf("interface %s not found", iface)
}

func getPublicIP(localIP net.IP) (string, error) {
	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{IP: localIP},
		Timeout:   5 * time.Second,
	}

	transport := &http.Transport{
		DialContext: dialer.DialContext,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   5 * time.Second,
	}

	resp, err := client.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := io.ReadAll(resp.Body)
	return string(ip), err
}

func isVPNActive(iface string) bool {
	wgDevices, err := getWireGuardDevices()
	if err != nil {
		fmt.Fprintln(os.Stderr, "WireGuard error:", err)
		return false
	}
	return slices.Contains(wgDevices, iface)
}

func getWireGuardDevices() ([]string, error) {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9999")
	if err != nil {
		return nil, fmt.Errorf("resolve error: %w", err)
	}

	conn, err := net.DialUDP("udp", nil, udpAddr)
	if err != nil {
		return nil, fmt.Errorf("dial error: %w", err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte("list"))
	if err != nil {
		return nil, fmt.Errorf("write error: %w", err)
	}

	buf := make([]byte, 1024)
	n, _, err := conn.ReadFromUDP(buf)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	resp := string(buf[:n])

	if strings.HasPrefix(resp, "ERROR:") {
		return nil, errors.New(resp)
	}
	// Trim "OK: " prefix and split into slice
	resp = strings.TrimPrefix(resp, "OK:")
	resp = strings.TrimSpace(resp)
	if resp == "" {
		return []string{}, nil
	}

	return strings.Fields(resp), nil
}
