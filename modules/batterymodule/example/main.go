package main

import (
	"fmt"

	"github.com/mamaart/statusbar/modules/batterymodule"
)

func main() {
	x := batterymodule.New()
	for x := range x.Reader() {
		fmt.Println(string(x))
	}
}
