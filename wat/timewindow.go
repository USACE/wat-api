package wat

import (
	"fmt"
	"time"
)

type TimeWindow struct {
	StartTime time.Time `json:"starttime"`
	EndTime   time.Time `json:"endtime"`
}

func (tw TimeWindow) IsValid() error {
	if tw.StartTime.Before(tw.EndTime) {
		return nil
	}
	return fmt.Errorf("invalid time window, StarTime %v must be before EndTIme %v", tw.StartTime, tw.EndTime)
}
