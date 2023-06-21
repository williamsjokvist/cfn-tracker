package utils

import "time"

// Sleeps for the passed duration or until the breakFunction returns true
func SleepOrBreak(duration time.Duration, breakFunction func() bool) bool {
	didBreak := false
	sleepPeriod := 1 * time.Second
	sleepCyclesLeft := int(duration / sleepPeriod)

	for sleepCyclesLeft > 0 {
		if breakFunction() {
			sleepCyclesLeft = 0
			didBreak = true
		}

		time.Sleep(sleepPeriod)
		sleepCyclesLeft--
	}

	return didBreak
}
