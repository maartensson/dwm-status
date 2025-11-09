package main

import (
	"fmt"

	"github.com/mamaart/statusbar/modules/vpnmodule"
)

func main() {
	m := vpnmodule.New()
	for r := range m.Reader() {
		fmt.Println(string(r))
	}
}
