package helper

import "time"

func IsToday(t time.Time) bool {
	now := time.Now()
	if t.Day() == now.Day() && t.Month() == now.Month() && t.Year() == now.Year() {
		return true
	}
	return false
}

func IsNeedUpdate(lastUpdate time.Time, threshold time.Duration) bool {
	now := time.Now()
	lastUpdate = lastUpdate.Add(threshold)
	return lastUpdate.Before(now)
}
