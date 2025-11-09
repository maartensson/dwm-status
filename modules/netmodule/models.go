package netmodule

import "fmt"

type IFace struct {
	Name string
	Addr string
}

func (f IFace) String() string {
	return fmt.Sprintf(" 󰌗 %s ", f.Addr)
}
