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
	output := make(chan []byte)
	go func() {
		for {
			w, err := Get()
			if err != nil {
				log.Println(err)
			} else {
				output <- []byte(w.String())
			}
			time.Sleep(time.Second * 5)
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
