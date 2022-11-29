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

func IsLast30Days(t time.Time) bool {
	timeNow := time.Now()
	dateNow := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, timeNow.Location())
	last30Days := dateNow.Add(-30 * time.Hour * 24)
	return t.After(last30Days)
}

func GetLast30Days() time.Time {
	timeNow := time.Now()
	dateNow := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, timeNow.Location())
	last30Days := dateNow.Add(-30 * time.Hour * 24)
	return last30Days
}

func IsFutureDate(t time.Time) bool {
	timeNow := time.Now()
	dateNow := time.Date(timeNow.Year(), timeNow.Month(), timeNow.Day(), 0, 0, 0, 0, timeNow.Location())
	return t.After(dateNow)
}

func ToDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}
