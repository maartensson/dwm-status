package batterymodule

import (
	"log"
	"strings"
	"syscall"
	"time"
)

// UdevPowerSupplyEventsThrottled returns a channel that emits
// at most once per second, coalescing rapid udev events.
func UdevPowerSupplyEventsThrottled(rateLimitDuration time.Duration) <-chan struct{} {
	raw := UdevPowerSupplyEvents()      // raw udev events
	throttled := make(chan struct{}, 1) // throttled output

	go func() {
		ticker := time.NewTicker(rateLimitDuration)
		defer ticker.Stop()

		var pending bool
		for {
			select {
			case <-raw:
				pending = true // mark that an event happened
			case <-ticker.C:
				if pending {
					select {
					case throttled <- struct{}{}:
					default: // drop if consumer is busy
					}
					pending = false
				}
			}
		}
	}()

	return throttled
}

// UdevPowerSupplyEvents returns a channel emitting a signal whenever
// a power_supply event occurs, throttled to at most once per second.
func UdevPowerSupplyEvents() <-chan struct{} {
	rawCh := make(chan struct{}, 1)       // raw udev events
	throttledCh := make(chan struct{}, 1) // throttled output

	go func() {
		fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, syscall.NETLINK_KOBJECT_UEVENT)
		if err != nil {
			log.Println("socket error:", err)
			close(rawCh)
			close(throttledCh)
			return
		}
		defer syscall.Close(fd)

		sa := &syscall.SockaddrNetlink{
			Family: syscall.AF_NETLINK,
			Pid:    0, // automatic unique PID
			Groups: 1, // subscribe to broadcast messages
		}

		if err := syscall.Bind(fd, sa); err != nil {
			log.Println("bind error:", err)
			close(rawCh)
			close(throttledCh)
			return
		}

		buf := make([]byte, 4096)
		for {
			n, _, err := syscall.Recvfrom(fd, buf, 0)
			if err != nil {
				log.Println("recvfrom error:", err)
				continue
			}

			msg := string(buf[:n])
			if msg != "" && isPowerSupplyEvent(msg) {
				// non-blocking send to raw channel
				select {
				case rawCh <- struct{}{}:
				default:
				}
			}
		}
	}()

	// throttler goroutine
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()

		var pending bool
		for {
			select {
			case <-rawCh:
				pending = true
			case <-ticker.C:
				if pending {
					select {
					case throttledCh <- struct{}{}:
					default: // drop if receiver is busy
					}
					pending = false
				}
			}
		}
	}()

	return throttledCh
}

// helper: filter for battery / power_supply events
func isPowerSupplyEvent(msg string) bool {
	return strings.Contains(msg, "POWER_SUPPLY") || strings.Contains(msg, "BAT")
}
