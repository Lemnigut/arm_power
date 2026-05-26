package repository

import (
	"context"

	"arm_back/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExerciseRepository interface {
	ListByUser(ctx context.Context, userID uuid.UUID, includeDeleted bool) ([]model.Exercise, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Exercise, error)
	Create(ctx context.Context, ex *model.Exercise) error
	Update(ctx context.Context, ex *model.Exercise) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	// Comments
	ListComments(ctx context.Context, exerciseID uuid.UUID, includeDeleted bool) ([]model.ExerciseComment, error)
	CreateComment(ctx context.Context, c *model.ExerciseComment) error
	UpdateComment(ctx context.Context, c *model.ExerciseComment) error
	SoftDeleteComment(ctx context.Context, id uuid.UUID) error
}

type pgExerciseRepo struct {
	pool *pgxpool.Pool
}

func NewExerciseRepository(pool *pgxpool.Pool) ExerciseRepository {
	return &pgExerciseRepo{pool: pool}
}

func (r *pgExerciseRepo) ListByUser(ctx context.Context, userID uuid.UUID, includeDeleted bool) ([]model.Exercise, error) {
	where := `user_id = $1 AND is_deleted = false`
	if includeDeleted {
		where = `user_id = $1`
	}
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, name, muscles, category, description, youtube_links, weight_unit, is_single_hand, is_deleted, created_at, updated_at
		 FROM exercises WHERE `+where+` ORDER BY name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return collectExercises(rows)
}

func (r *pgExerciseRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Exercise, error) {
	var ex model.Exercise
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, name, muscles, category, description, youtube_links, weight_unit, is_single_hand, is_deleted, created_at, updated_at
		 FROM exercises WHERE id = $1`, id,
	).Scan(&ex.ID, &ex.UserID, &ex.Name, &ex.Muscles, &ex.Category, &ex.Description, &ex.YoutubeLinks, &ex.WeightUnit, &ex.IsSingleHand, &ex.IsDeleted, &ex.CreatedAt, &ex.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return &ex, err
}

func (r *pgExerciseRepo) Create(ctx context.Context, ex *model.Exercise) error {
	createdAt := nullableTime(ex.CreatedAt)
	_, err := r.pool.Exec(ctx,
		`INSERT INTO exercises (id, user_id, name, muscles, category, description, youtube_links, weight_unit, is_single_hand, is_deleted, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, COALESCE($11, now()))`,
		ex.ID, ex.UserID, ex.Name, ex.Muscles, ex.Category, ex.Description, ex.YoutubeLinks, ex.WeightUnit, ex.IsSingleHand, ex.IsDeleted, createdAt,
	)
	return err
}

func (r *pgExerciseRepo) Update(ctx context.Context, ex *model.Exercise) error {
	createdAt := nullableTime(ex.CreatedAt)
	_, err := r.pool.Exec(ctx,
		`UPDATE exercises
		 SET name=$2, muscles=$3, category=$4, description=$5, youtube_links=$6, weight_unit=$7, is_single_hand=$8, is_deleted=$9, created_at=COALESCE($10, created_at), updated_at=now()
		 WHERE id=$1`,
		ex.ID, ex.Name, ex.Muscles, ex.Category, ex.Description, ex.YoutubeLinks, ex.WeightUnit, ex.IsSingleHand, ex.IsDeleted, createdAt,
	)
	return err
}

func (r *pgExerciseRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE exercises SET is_deleted=true, updated_at=now() WHERE id=$1`, id)
	return err
}

// Comments

func (r *pgExerciseRepo) ListComments(ctx context.Context, exerciseID uuid.UUID, includeDeleted bool) ([]model.ExerciseComment, error) {
	where := `exercise_id=$1 AND is_deleted=false`
	if includeDeleted {
		where = `exercise_id=$1`
	}
	rows, err := r.pool.Query(ctx,
		`SELECT id, exercise_id, user_id, text, is_deleted, created_at, updated_at
		 FROM exercise_comments WHERE `+where+` ORDER BY created_at`, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []model.ExerciseComment
	for rows.Next() {
		var c model.ExerciseComment
		if err := rows.Scan(&c.ID, &c.ExerciseID, &c.UserID, &c.Text, &c.IsDeleted, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *pgExerciseRepo) CreateComment(ctx context.Context, c *model.ExerciseComment) error {
	createdAt := nullableTime(c.CreatedAt)
	_, err := r.pool.Exec(ctx,
		`INSERT INTO exercise_comments (id, exercise_id, user_id, text, is_deleted, created_at)
		 VALUES ($1, $2, $3, $4, $5, COALESCE($6, now()))`,
		c.ID, c.ExerciseID, c.UserID, c.Text, c.IsDeleted, createdAt,
	)
	return err
}

func (r *pgExerciseRepo) UpdateComment(ctx context.Context, c *model.ExerciseComment) error {
	createdAt := nullableTime(c.CreatedAt)
	_, err := r.pool.Exec(ctx,
		`UPDATE exercise_comments
		 SET text=$2, is_deleted=$3, created_at=COALESCE($4, created_at), updated_at=now()
		 WHERE id=$1`,
		c.ID, c.Text, c.IsDeleted, createdAt,
	)
	return err
}

func (r *pgExerciseRepo) SoftDeleteComment(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE exercise_comments SET is_deleted=true, updated_at=now() WHERE id=$1`, id)
	return err
}

func collectExercises(rows pgx.Rows) ([]model.Exercise, error) {
	var exercises []model.Exercise
	for rows.Next() {
		var ex model.Exercise
		if err := rows.Scan(&ex.ID, &ex.UserID, &ex.Name, &ex.Muscles, &ex.Category, &ex.Description, &ex.YoutubeLinks, &ex.WeightUnit, &ex.IsSingleHand, &ex.IsDeleted, &ex.CreatedAt, &ex.UpdatedAt); err != nil {
			return nil, err
		}
		exercises = append(exercises, ex)
	}
	return exercises, rows.Err()
}
