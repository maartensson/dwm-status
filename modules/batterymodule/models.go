package batterymodule

import "fmt"

type Battery struct {
	Charging bool
	Capacity int
}

//⚡, 🔌, 🔋⚡, 🔋🔌

func (b Battery) String(flash bool) string {
	if flash {
		if b.Charging {
			return fmt.Sprintf(" 🔌 %d%% ", b.Capacity)
		} else if b.Capacity < 25 {
			return fmt.Sprintf(" 💀 %d%% ", b.Capacity)
		} else {
			return fmt.Sprintf(" 💡 %d%% ", b.Capacity)
		}
	} else {
		if b.Charging {
			return fmt.Sprintf(" ⚡ %d%% ", b.Capacity)
		} else if b.Capacity < 25 {
			return fmt.Sprintf(" 🪫 %d%% ", b.Capacity)
		} else {
			return fmt.Sprintf(" 🔋 %d%% ", b.Capacity)
		}
	}
}
