package helper

import (
	"time"
	"twfinder/logger"
)

// StringtoDate : convert string to datetime.
// default layout : Ruby format "Mon Jan 02 15:04:05 -0700 2006"
func StringtoDate(date string, layout string) time.Time {
	if layout == "" {
		layout = time.RubyDate
	}
	res, err := time.Parse(layout, date)
	if err != nil {
		logger.Errorf("Error while converting string to date input:%v layout:%v", date, layout)
		return time.Time{}
	}
	return res
}
