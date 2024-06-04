package main

import "time"

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
