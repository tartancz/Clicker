package main

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

func FormatTime(t time.Time) string {
	return t.Format(time.TimeOnly)
}

func parseDuration(lenght int, unit string) time.Duration {
	var dur time.Duration
	switch unit {
	case "seconds":
		dur = time.Second * time.Duration(lenght)
	case "minutes":
		dur = time.Minute * time.Duration(lenght)
	case "hours":
		dur = time.Hour * time.Duration(lenght)
	}
	return dur
}

func FormatDuration(dur time.Duration) (int, string) {
	if sec := int(dur.Seconds()); sec < 0 {
		return 0, "seconds"
	} else if sec < 60 {
		return sec, "seconds"
	} else if sec < 3600 {
		return sec / 60, "minutes"
	} else {
		return sec / 3600, "hours"
	}
}

// returns time.Time from string in format "HH:MM:SS" with today Date
func NewDateWithTimeFromString(t string) (time.Time, error) {
	times := strings.Split(t, ":")
	if len(times) != 3 {
		return time.Time{}, errors.New("invalid time format")
	}
	hour, err := strconv.Atoi(times[0])
	if err != nil || hour < 0 || hour > 23 {
		return time.Time{}, errors.New("invalid hour")
	}
	min, err := strconv.Atoi(times[1])
	if err != nil || min < 0 || min > 59 {
		return time.Time{}, errors.New("invalid minute")
	}
	sec, err := strconv.Atoi(times[2])
	if err != nil || sec < 0 || sec > 59 {
		return time.Time{}, errors.New("invalid second")
	}
	today := time.Now()
	result := time.Date(today.Year(), today.Month(), today.Day(), hour, min, sec, 0, time.Local)
	return result, nil
}
