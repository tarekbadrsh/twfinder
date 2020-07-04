package helper

import "time"

// StringtoDate : convert string to datetime using Ruby format.
func StringtoDate(date string) time.Time {
	if res, err := time.Parse(time.RubyDate, date); err == nil {
		return res
	}
	return time.Time{}
}
