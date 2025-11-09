package ui

import (
	"strings"

	"github.com/mamaart/statusbar/pkg/bar"
)

func Run(
	vpn <-chan []byte,
	tim <-chan []byte,
	bat <-chan []byte,
	vol <-chan []byte,
	bri <-chan []byte,
	wtr <-chan []byte,
	net <-chan []byte,
	dsk <-chan []byte,
	txt <-chan []byte,
) {
	state := make(chan []byte)

	go func() {
		chunks := make([]string, 9)
		for {
			select {
			case data := <-vpn:
				chunks[0] = string(data)
			case data := <-tim:
				chunks[1] = string(data)
			case data := <-bat:
				chunks[2] = string(data)
			case data := <-vol:
				chunks[3] = string(data)
			case data := <-bri:
				chunks[4] = string(data)
			case data := <-wtr:
				chunks[5] = string(data)
			case data := <-net:
				chunks[6] = string(data)
			case data := <-dsk:
				chunks[7] = string(data)
			case data := <-txt:
				chunks[8] = string(data)
			}
			state <- []byte(strings.Join(chunks, "|"))
		}
	}()

	bar := bar.New()
	for x := range state {
		bar.Update(x)
	}
}
