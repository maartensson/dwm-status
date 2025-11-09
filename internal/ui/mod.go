package ui

import (
	"reflect"
	"strings"

	"github.com/mamaart/statusbar/pkg/bar"
)

type Module interface {
	Reader() <-chan []byte
}

func Run(
	modules []Module,
) {
	state := make(chan []byte)
	go func() {
		chunks := make([]string, len(modules))

		cases := make([]reflect.SelectCase, len(modules))
		for i, module := range modules {
			cases[i] = reflect.SelectCase{
				Dir:  reflect.SelectRecv,
				Chan: reflect.ValueOf(module.Reader()),
			}
		}

		for {
			chosen, value, ok := reflect.Select(cases)
			if !ok {
				continue
			}
			chunks[chosen] = string(value.Bytes())
			state <- []byte(strings.Join(chunks, "|"))
		}
	}()

	bar := bar.New()
	for x := range state {
		bar.Update(x)
	}
}
