package main

import (
	"os"
	"time"

	"github.com/mamaart/statusbar/internal/api"
	"github.com/mamaart/statusbar/internal/ui"
	"github.com/mamaart/statusbar/modules/batterymodule"
	"github.com/mamaart/statusbar/modules/brightnessmodule"
	"github.com/mamaart/statusbar/modules/diskmodule"
	"github.com/mamaart/statusbar/modules/textmodule"
	"github.com/mamaart/statusbar/modules/timemodule"
	"github.com/mamaart/statusbar/modules/volumemodule"
	"github.com/mamaart/statusbar/modules/vpnmodule"
)

func main() {
	socketPath := os.Getenv("SOCKET_PATH")
	if socketPath == "" {
		socketPath = "/run/wg-helper/wg-helper.sock"
	}

	api := api.New()
	tim := timemodule.New()
	txt := textmodule.New(textmodule.Options{
		WindowWidth: 80,
		Delay:       time.Millisecond * 150,
	})

	go api.Run(tim, txt)

	ui.Run([]ui.Module{
		vpnmodule.New(socketPath),
		diskmodule.New(),
		brightnessmodule.New(),
		volumemodule.New(),
		batterymodule.New(),
		tim,
		txt,
	})
}
