package timemodule

import (
	"errors"
	"fmt"
	"time"

	"github.com/mamaart/statusbar/modules/timemodule/watchface"
)

type TimeModule struct {
	output chan []byte
	toggle chan struct{}
}

func New() *TimeModule {
	output := make(chan []byte)
	toggle := make(chan struct{})
	go func() {
		state := watchface.Clock
		t := time.Now()
		output <- []byte(Clock(t.Format(fmt.Sprintf(" %s 15:04", getClockIcon(t)))).String())
		for {
			t := time.Now()
			select {
			case <-toggle:
				switch state {
				case watchface.Clock:
					state = watchface.Date
					output <- []byte(Calendar(t.Format("02 Jan 2006")).String())
				case watchface.Date:
					state = watchface.Weeknumber
					_, w := t.ISOWeek()
					output <- []byte(WeekNo(fmt.Sprintf("Week: %d", w)).String())
				case watchface.Weeknumber:
					state = watchface.Day
					output <- []byte(Day(t.Weekday().String()).String())
				case watchface.Day:
					state = watchface.Clock
					output <- []byte(Clock(t.Format("15:04")).String())
				}
			case <-time.After(time.Second * 5):
				switch state {
				case watchface.Clock:
					output <- []byte(Clock(t.Format(fmt.Sprintf(" %s 15:04", getClockIcon(t)))).String())
				case watchface.Date:
					output <- []byte(Calendar(t.Format("02 Jan 2006")).String())
				case watchface.Weeknumber:
					_, w := t.ISOWeek()
					output <- []byte(WeekNo(fmt.Sprintf("Week: %d", w)).String())
				case watchface.Day:
					output <- []byte(Day(t.Weekday().String()).String())
				}
			}
		}
	}()
	return &TimeModule{
		output: output,
		toggle: toggle,
	}
}

func (t *TimeModule) Reader() <-chan []byte {
	return t.output
}

func (d *TimeModule) Toggle() error {
	select {
	case d.toggle <- struct{}{}:
		return nil
	case <-time.After(time.Second * 5):
		return errors.New("timeout while trying to toggle")
	}
}
