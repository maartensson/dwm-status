package batterymodule

import "fmt"

type Battery struct {
	Charging bool
	Capacity int
}

func (b Battery) String(flash bool) string {
	if flash {
		if b.Charging {
			if b.Capacity < 80 {
				return fmt.Sprintf(" 󱊦 %d%% ", b.Capacity)
			}

			return fmt.Sprintf(" 󱊥 %d%% ", b.Capacity)
		}

		if b.Capacity < 25 {
			return fmt.Sprintf("  %d%% ", b.Capacity)
		}

		if b.Capacity < 50 {
			return fmt.Sprintf(" 󰁻 %d%% ", b.Capacity)
		}

		return fmt.Sprintf(" 󰁾 %d%% ", b.Capacity)

	} else {
		if b.Charging {

			if b.Capacity < 20 {
				return fmt.Sprintf(" 󰢟 %d%% ", b.Capacity)
			}
			if b.Capacity < 40 {
				return fmt.Sprintf(" 󱊤 %d%% ", b.Capacity)
			}
			if b.Capacity < 80 {
				return fmt.Sprintf(" 󱊥 %d%% ", b.Capacity)
			}
			return fmt.Sprintf(" 󱊦 %d%% ", b.Capacity)

		}

		if b.Capacity < 10 {
			return fmt.Sprintf(" 󰁺 %d%% ", b.Capacity)
		}
		if b.Capacity < 20 {
			return fmt.Sprintf(" 󰁻 %d%% ", b.Capacity)
		}
		if b.Capacity < 30 {
			return fmt.Sprintf(" 󰁼 %d%% ", b.Capacity)
		}
		if b.Capacity < 40 {
			return fmt.Sprintf(" 󰁽 %d%% ", b.Capacity)
		}
		if b.Capacity < 50 {
			return fmt.Sprintf(" 󰁾 %d%% ", b.Capacity)
		}
		if b.Capacity < 60 {
			return fmt.Sprintf(" 󰁿 %d%% ", b.Capacity)
		}
		if b.Capacity < 70 {
			return fmt.Sprintf(" 󰂀 %d%% ", b.Capacity)
		}
		if b.Capacity < 80 {
			return fmt.Sprintf(" 󰂁 %d%% ", b.Capacity)
		}
		if b.Capacity < 90 {
			return fmt.Sprintf(" 󰂂 %d%% ", b.Capacity)
		}
		return fmt.Sprintf(" 󱈑 %d%% ", b.Capacity)
	}
}
