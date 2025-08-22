package time_events

import (
	"context"
	"time"
)

func (t *TimeEvents) StartDailyCheck(ctx context.Context) {
	for {
		now := time.Now()
		// Next midnight
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		duration := next.Sub(now)
		time.Sleep(duration)
		// Call your function here
		t.checkHabitsAndMarkCalendar(ctx)
	}
}
