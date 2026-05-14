package domain

import (
	"fmt"
	"time"
)

type CustomDuration struct {
	time.Duration
}

func (d CustomDuration) String() string {
	t := time.Time{}.Add(time.Duration(d.Duration))
	return fmt.Sprintf("%s", t.Format("15:04:05"))
}

func AverageDuration(durations []CustomDuration) CustomDuration {
	if len(durations) == 0 {
		return CustomDuration{0}
	}
	var sum CustomDuration
	for _, d := range durations {
		sum.Duration += d.Duration
	}
	return CustomDuration{sum.Duration / time.Duration(len(durations))}
}
