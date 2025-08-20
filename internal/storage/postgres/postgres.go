package postgres

import (
	"context"
	"habit-check-bot/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"go.uber.org/zap"
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
	habit *models.Habit,
) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer func(ctx2 context.Context) {
		_ = tx.Rollback(ctx2)
	}(ctx)

	const saveHabit = `
INSERT INTO habits (name, chatId)
VALUES ($1, $2)`

	_, err = tx.Exec(ctx, saveHabit,
		habit.Name,
		habit.ChatId,
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

	const deleteHabit = `
DELETE FROM habits
WHERE name = $1 AND chatId = $2
`

	_, err = tx.Exec(ctx, deleteHabit, name, chatId)
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

func (p *postgresRepository) GetHabits(
	ctx context.Context,
	chatId int64,
) []models.Habit {

	const getHabitsQuery = `
	SELECT name, chatId, mark
	FROM habits
	WHERE chatId = $1
`
	rows, err := p.db.Query(ctx, getHabitsQuery, chatId)
	if err != nil {
		p.logger.Error("failed to get habits", zap.Error(err))
		return []models.Habit{}
	}
	defer rows.Close()

	habits := make([]models.Habit, 0)
	for rows.Next() {
		var habit models.Habit
		if err := rows.Scan(&habit.Name, &habit.ChatId, &habit.Mark); err != nil {
			p.logger.Error("failed to scan habit", zap.Error(err))
			return []models.Habit{}
		}
		habits = append(habits, habit)
	}
	return habits
}

func (p *postgresRepository) MarkHabit(
	ctx context.Context,
	name string,
	chatId int64,
	mark bool,
) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func(ctx2 context.Context) {
		_ = tx.Rollback(ctx2)
	}(ctx)

	const updateMark = `
		UPDATE habits
		SET mark = $1
		WHERE name = $2 AND chatId = $3
	`
	_, err = tx.Exec(ctx, updateMark, mark, name, chatId)
	if err != nil {
		p.logger.Error("failed to update mark", zap.Error(err))
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		p.logger.Error("failed to commit tx", zap.Error(err))
		return err
	}

	return nil
}

func (p *postgresRepository) GetDatesHabit(
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
