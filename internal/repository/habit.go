package repository

import (
	"context"

	"arm_back/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HabitRepository interface {
	ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Habit, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Habit, error)
	Create(ctx context.Context, h *model.Habit) error
	Update(ctx context.Context, h *model.Habit) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	CreateCompletion(ctx context.Context, c *model.HabitCompletion) error
	GetCompletionByID(ctx context.Context, id uuid.UUID) (*model.HabitCompletion, error)
	SoftDeleteCompletion(ctx context.Context, id uuid.UUID) error
}

type pgHabitRepo struct {
	pool *pgxpool.Pool
}

func NewHabitRepository(pool *pgxpool.Pool) HabitRepository {
	return &pgHabitRepo{pool: pool}
}

func (r *pgHabitRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Habit, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, title, color, repeat_type, repeat_mode, repeat_weekday, is_deleted, created_at, updated_at
		 FROM habits WHERE user_id=$1 AND is_deleted=false ORDER BY created_at`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	habits, err := collectHabits(rows)
	if err != nil {
		return nil, err
	}

	for i := range habits {
		completions, err := r.loadCompletions(ctx, habits[i].ID)
		if err != nil {
			return nil, err
		}
		habits[i].Completions = completions
	}
	return habits, nil
}

func (r *pgHabitRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Habit, error) {
	var h model.Habit
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, title, color, repeat_type, repeat_mode, repeat_weekday, is_deleted, created_at, updated_at
		 FROM habits WHERE id=$1`, id,
	).Scan(&h.ID, &h.UserID, &h.Title, &h.Color, &h.RepeatType, &h.RepeatMode, &h.RepeatWeekday, &h.IsDeleted, &h.CreatedAt, &h.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	completions, err := r.loadCompletions(ctx, h.ID)
	if err != nil {
		return nil, err
	}
	h.Completions = completions
	return &h, nil
}

func (r *pgHabitRepo) Create(ctx context.Context, h *model.Habit) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO habits (id, user_id, title, color, repeat_type, repeat_mode, repeat_weekday)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		h.ID, h.UserID, h.Title, h.Color, h.RepeatType, h.RepeatMode, h.RepeatWeekday,
	)
	return err
}

func (r *pgHabitRepo) Update(ctx context.Context, h *model.Habit) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE habits
		 SET title=$2, color=$3, repeat_type=$4, repeat_mode=$5, repeat_weekday=$6, updated_at=now()
		 WHERE id=$1`,
		h.ID, h.Title, h.Color, h.RepeatType, h.RepeatMode, h.RepeatWeekday,
	)
	return err
}

func (r *pgHabitRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx,
		`UPDATE habit_completions SET is_deleted=true, updated_at=now() WHERE habit_id=$1`, id); err != nil {
		return err
	}
	if _, err := tx.Exec(ctx,
		`UPDATE habits SET is_deleted=true, updated_at=now() WHERE id=$1`, id); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func (r *pgHabitRepo) CreateCompletion(ctx context.Context, c *model.HabitCompletion) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO habit_completions (id, habit_id, date)
		 VALUES ($1, $2, $3)
		 ON CONFLICT (habit_id, date) DO UPDATE
		 SET id=EXCLUDED.id, is_deleted=false, updated_at=now()`,
		c.ID, c.HabitID, c.Date,
	)
	return err
}

func (r *pgHabitRepo) GetCompletionByID(ctx context.Context, id uuid.UUID) (*model.HabitCompletion, error) {
	var c model.HabitCompletion
	err := r.pool.QueryRow(ctx,
		`SELECT id, habit_id, date, is_deleted, created_at, updated_at
		 FROM habit_completions WHERE id=$1`, id,
	).Scan(&c.ID, &c.HabitID, &c.Date, &c.IsDeleted, &c.CreatedAt, &c.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return &c, err
}

func (r *pgHabitRepo) SoftDeleteCompletion(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE habit_completions SET is_deleted=true, updated_at=now() WHERE id=$1`, id)
	return err
}

func (r *pgHabitRepo) loadCompletions(ctx context.Context, habitID uuid.UUID) ([]model.HabitCompletion, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, habit_id, date, is_deleted, created_at, updated_at
		 FROM habit_completions WHERE habit_id=$1 AND is_deleted=false ORDER BY date`, habitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var completions []model.HabitCompletion
	for rows.Next() {
		var c model.HabitCompletion
		if err := rows.Scan(&c.ID, &c.HabitID, &c.Date, &c.IsDeleted, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		completions = append(completions, c)
	}
	return completions, rows.Err()
}

func collectHabits(rows pgx.Rows) ([]model.Habit, error) {
	var habits []model.Habit
	for rows.Next() {
		var h model.Habit
		if err := rows.Scan(&h.ID, &h.UserID, &h.Title, &h.Color, &h.RepeatType, &h.RepeatMode, &h.RepeatWeekday, &h.IsDeleted, &h.CreatedAt, &h.UpdatedAt); err != nil {
			return nil, err
		}
		habits = append(habits, h)
	}
	return habits, rows.Err()
}
