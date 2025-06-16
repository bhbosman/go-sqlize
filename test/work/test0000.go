package work

import "time"

func init() {
	now01 := time.Now()
	now02 := time.Now()
	time.Time.Before(now01, now02)
}
