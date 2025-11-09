package volumemodule

import "fmt"

type Volume int

func (v Volume) String() string {
	if v < 10 {
		return fmt.Sprintf("  %02d%% ", v)
	}
	if v < 65 {
		return fmt.Sprintf("  %d%% ", v)
	}
	return fmt.Sprintf("  %d%% ", v)
}
