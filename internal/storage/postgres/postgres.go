package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"

	"habit-check-bot/internal/models"
)

type (
	PgxInterface interface {
		Begin(context.Context) (pgx.Tx, error)
		Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
		Query(context.Context, string, ...interface{}) (pgx.Rows, error)
		QueryRow(context.Context, string, ...interface{}) pgx.Row
	}
)

type postgresRepository struct {
	db     PgxInterface
	logger *zap.Logger
}

func New(
	db PgxInterface,
	logger *zap.Logger,
) *postgresRepository {
	return &postgresRepository{
		db:     db,
		logger: logger,
	}
}

func (p *postgresRepository) SaveHabit(
	ctx context.Context,
	name string,
	chatId int64,
) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(ctx2 context.Context) {
		_ = tx.Rollback(ctx2)
	}(ctx)

	const insertIntoHabits = `
INSERT INTO habits (name, chatId)
VALUES ($1, $2)`

	_, err = tx.Exec(ctx, insertIntoHabits,
		name,
		chatId,
	)
	if err != nil {
		p.logger.Error("failed to save habit", zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.logger.Error("failed to commit tx", zap.Error(err))
		return err
	}

	return nil
}

func (p *postgresRepository) DeleteHabit(
	ctx context.Context,
	name string,
	chatId int64,
) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(ctx2 context.Context) {
		_ = tx.Rollback(ctx2)
	}(ctx)

	const deleteFromHabits = `
DELETE FROM habits
WHERE name = $1 AND chatId = $2
`

	_, err = tx.Exec(ctx, deleteFromHabits, name, chatId)
	if err != nil {
		p.logger.Error("failed to delete habit", zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.logger.Error("failed to commit tx", zap.Error(err))
		return err
	}

	return nil
}

func (p *postgresRepository) GetAccountsHabits(
	ctx context.Context,
	chatId int64,
) ([]models.Habit, error) {

	const getTodayMarkHabits = `
	SELECT habits.name, habits.chatId
	FROM habits join time_done
	ON habits.name = time_done.name
	AND habits.chatId = time_done.chatId
	AND time_done.created_at = CURRENT_DATE
	WHERE habits.chatId = $1
`
	rows, err := p.db.Query(ctx, getTodayMarkHabits, chatId)
	if err != nil {
		p.logger.Error("failed to get habits", zap.Error(err))
		return []models.Habit{}, err
	}

	habits := make([]models.Habit, 0)
	for rows.Next() {
		habit := models.Habit{
			Mark: true,
		}
		if err := rows.Scan(&habit.Name, &habit.ChatId); err != nil {
			p.logger.Error("failed to scan habit", zap.Error(err))
			return []models.Habit{}, err
		}
		habits = append(habits, habit)
	}
	rows.Close()

	const getNotTodayMarkHabits = `
SELECT h.name, h.chatId
FROM habits h
LEFT JOIN time_done td
  ON h.name = td.name
  AND h.chatId = td.chatId
  AND td.created_at = CURRENT_DATE
WHERE td.name IS NULL and h.chatId = $1`

	rows, err = p.db.Query(ctx, getNotTodayMarkHabits, chatId)
	if err != nil {
		p.logger.Error("failed to get habits", zap.Error(err))
		return []models.Habit{}, err
	}

	for rows.Next() {
		habit := models.Habit{
			Mark: false,
		}
		if err := rows.Scan(&habit.Name, &habit.ChatId); err != nil {
			p.logger.Error("failed to scan habit", zap.Error(err))
			return []models.Habit{}, err
		}
		habits = append(habits, habit)
	}
	rows.Close()

	return habits, nil
}

func (p *postgresRepository) MarkHabit(
	ctx context.Context,
	name string,
	chatId int64,
) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(ctx2 context.Context) {
		_ = tx.Rollback(ctx2)
	}(ctx)

	const deleteFromTime = `
		DELETE FROM time_done
		WHERE name = $1 AND chatId = $2 and created_at = CURRENT_DATE
	`
	cmd, err := tx.Exec(ctx, deleteFromTime, name, chatId)
	if err != nil {
		p.logger.Error("failed to update mark", zap.Error(err))
		return err
	}

	if cmd.RowsAffected() == 0 {
		const insertIntoTime = `
		INSERT INTO time_done (name, chatId, created_at)
		VALUES ($1, $2, CURRENT_DATE)
	`
		_, err = tx.Exec(ctx, insertIntoTime, name, chatId)
		if err != nil {
			p.logger.Error("failed to insert date", zap.Error(err))
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.logger.Error("failed to commit tx", zap.Error(err))
		return err
	}

	return nil
}

func (p *postgresRepository) GetHabitsDates(
	ctx context.Context,
	name string,
	chatId int64,
) ([]string, error) {
	const getTimes = `
SELECT created_at from time_done
where name = $1 AND chatId = $2
`

	rows, err := p.db.Query(ctx, getTimes, name, chatId)
	if err != nil {
		p.logger.Error("failed to get dates of habit", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	ans := make([]string, 0)
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			p.logger.Error("failed to scan date", zap.Error(err))
			return nil, err
		}
		ans = append(ans, date)
	}

	return ans, nil
}

func (p *postgresRepository) GetAllChatIDs(
	ctx context.Context,
) ([]int64, error) {
	const getChatIds = `SELECT DISTINCT chatId FROM habits`
	rows, err := p.db.Query(ctx, getChatIds)
	if err != nil {
		p.logger.Error("failed to get chat IDs", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var chatIDs []int64
	for rows.Next() {
		var chatId int64
		if err := rows.Scan(&chatId); err != nil {
			p.logger.Error("failed to scan chatId", zap.Error(err))
			return nil, err
		}
		chatIDs = append(chatIDs, chatId)
	}
	return chatIDs, nil
}

func (p *postgresRepository) MarkCalendarDay(
	ctx context.Context,
	chatId int64,
	date string,
	mark bool,
) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(ctx2 context.Context) {
		_ = tx.Rollback(ctx2)
	}(ctx)

	const deleteFromCalendarDay = `
		DELETE FROM calendar
		WHERE chatId = $1 AND date = $2 and mark = $3
	`
	cmd, err := tx.Exec(ctx, deleteFromCalendarDay, chatId, date, mark)
	if err != nil {
		p.logger.Error("failed to update mark", zap.Error(err))
		return err
	}

	if cmd.RowsAffected() == 0 {
		const insertIntoCalendarDay = `
		INSERT INTO calendar (chatId, date, mark)
		VALUES ($1, $2, $3)
	`
		_, err = tx.Exec(ctx, insertIntoCalendarDay, chatId, date, mark)
		if err != nil {
			p.logger.Error("failed to insert date", zap.Error(err))
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.logger.Error("failed to commit tx", zap.Error(err))
		return err
	}

	return nil
}

func (p *postgresRepository) GetMarkDay(
	ctx context.Context,
	chatId int64,
	date string,
) (bool, error) {
	const getFromCalendar = `
	select mark from calendar
	where chatId = $1 and date = $2
`
	row := p.db.QueryRow(ctx, getFromCalendar, chatId, date)
	var mark bool
	if err := row.Scan(&mark); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, err
		}
		p.logger.Error("failed to get mark for day", zap.Error(err))
		return false, err
	}
	return mark, nil
}
