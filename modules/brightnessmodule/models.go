package brightnessmodule

import "fmt"

type Brightness int

func (b Brightness) String() string {
	if b < 20 {
		return fmt.Sprintf(" 󰃞 %02d%% ", b)
	}
	if b < 70 {
		return fmt.Sprintf(" 󰃟 %d%% ", b)
	}
	return fmt.Sprintf(" 󰃠 %d%% ", b)
}
