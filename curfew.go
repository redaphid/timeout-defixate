package main

import (
	"fmt"
	"os"
	"time"
)

// StartCurfew re-locks the screen every relock interval while the local clock
// sits inside [startHour, endHour) (wrapping past midnight). Touching bypassPath
// suppresses it for the current night; the file is cleared once the window ends,
// so a forgotten override re-arms itself the next night.
func StartCurfew(startHour, endHour int, relock time.Duration, bypassPath string) {
	logMsg(fmt.Sprintf("Curfew armed: %02d:00–%02d:00, relock every %s (bypass: %s)", startHour, endHour, relock, bypassPath))

	for {
		if !inCurfew(time.Now(), startHour, endHour) {
			clearBypass(bypassPath)
			time.Sleep(relock)
			continue
		}

		if bypassed(bypassPath) {
			time.Sleep(relock)
			continue
		}

		logWarn("Curfew - locking. Go to bed.")
		lockScreen()
		time.Sleep(relock)
	}
}

// inCurfew reports whether the hour of t falls in [startHour, endHour),
// handling windows that wrap past midnight (e.g. 22 to 8).
func inCurfew(t time.Time, startHour, endHour int) bool {
	h := t.Hour()
	if startHour <= endHour {
		return h >= startHour && h < endHour
	}
	return h >= startHour || h < endHour
}

func bypassed(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func clearBypass(path string) {
	os.Remove(path)
}
