package timemodule

import (
	"time"
)

type Time interface {
	String() string
}

type Calendar string

func (c Calendar) String() string {
	return " 📅 " + string(c) + " "
}

type Clock string

func (c Clock) String() string {
	return string(c) + " "
}

type WeekNo string

func (w WeekNo) String() string {
	return " 📅 " + string(w) + " "
}

type Day string

func (d Day) String() string {
	return " 📅 " + string(d) + " "
}

func getClockIcon(t time.Time) string {
	hour := t.Hour() % 12 // convert to 12-hour format
	if hour == 0 {
		hour = 12
	}

	switch hour {
	case 1:
		return "󱑋"
	case 2:
		return "󱑌"
	case 3:
		return "󱑍"
	case 4:
		return "󱑎"
	case 5:
		return "󱑏"
	case 6:
		return "󱑐"
	case 7:
		return "󱑑"
	case 8:
		return "󱑒"
	case 9:
		return "󱑓"
	case 10:
		return "󱑔"
	case 11:
		return "󱑕"
	case 12:
		return "󱑖"
	default:
		return "󱑆" // fallback (shouldn't happen)
	}
}
