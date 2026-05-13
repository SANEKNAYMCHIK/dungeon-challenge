package domain

import (
	"fmt"
	"time"
)

type CustomTime struct {
	time.Time
}

func (t *CustomTime) UnmarshalJSON(data []byte) error {
	strTime := string(data[1 : len(data)-1])
	parsedTime, err := time.Parse(time.TimeOnly, strTime)
	if err != nil {
		return err
	}
	t.Time = parsedTime
	return nil
}

func (t CustomTime) String() string {
	return fmt.Sprintf("[%s]", t.Format("15:04:05"))
}
