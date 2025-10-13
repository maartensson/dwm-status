package batterymodule

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type BatteryModule struct {
	output <-chan []byte
}

func New() *BatteryModule {
	output := make(chan []byte, 10)
	go func() {

		udevCh := UdevPowerSupplyEventsThrottled(time.Second)

		old, err := Get()
		if err != nil {
			log.Fatal(err)
		}
		output <- []byte(old.String(false))

		duration := time.Minute
		for {

			select {
			case <-udevCh:
			case <-time.After(duration):
			}

			new, err := Get()
			if err != nil {
				log.Println(err)
				continue
			}

			if old.Capacity != new.Capacity || old.Charging != new.Charging {
				flash := false
				for range 10 {
					flash = !flash
					output <- []byte(new.String(flash))
					time.Sleep(time.Millisecond * 200)
				}
			}
			if !new.Charging && new.Capacity < 25 {
				flash := false
				for range 50 {
					flash = !flash
					output <- []byte(new.String(flash))
					time.Sleep(time.Millisecond * 200)
				}
				duration = time.Second * 5
			} else {
				duration = time.Minute
			}

			output <- []byte(new.String(false))
			old = new
		}
	}()

	return &BatteryModule{
		output: output,
	}
}

func (b *BatteryModule) Reader() <-chan []byte {
	return b.output
}

func Get() (b Battery, err error) {
	capacity, err := os.ReadFile("/sys/class/power_supply/BAT0/capacity")
	if err != nil {
		return b, fmt.Errorf("failed to open capacity file: %s", err)
	}
	stat, err := os.ReadFile("/sys/class/power_supply/BAT0/status")
	if err != nil {
		return b, fmt.Errorf("failed to open status file: %s", err)
	}

	value, err := strconv.Atoi(string(strings.TrimSpace(string(capacity))))
	if err != nil {
		return b, fmt.Errorf("failed to parse capacity: %s", err)
	}

	return Battery{
		Charging: strings.TrimSpace(string(stat)) == "Charging",
		Capacity: value,
	}, nil
}
