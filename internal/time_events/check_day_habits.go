package time_events

import (
	"context"
	"time"
)

func (t *TimeEvents) StartDailyCheck(ctx context.Context) {
	for {
		now := time.Now()

		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		duration := next.Sub(now)
		time.Sleep(duration)

		t.checkHabitsAndMarkCalendar(ctx)
	}
}

func (t *TimeEvents) checkHabitsAndMarkCalendar(ctx context.Context) {
	today := time.Now().Add(-24 * time.Hour).Format("2006-01-02")

	users, err := t.db.GetAllChatIDs(ctx)
	if err != nil {
		t.logger.Error(err.Error())
		return
	}
	for _, chatId := range users {
		habits, err := t.db.GetAccountsHabits(ctx, chatId)
		if err != nil {
			t.logger.Error(err.Error())
			return
		}
		allDone := true
		for _, habit := range habits {
			if !habit.Mark {
				allDone = false
				break
			}
		}

		t.db.MarkCalendarDay(ctx, chatId, today, allDone)
	}
}
