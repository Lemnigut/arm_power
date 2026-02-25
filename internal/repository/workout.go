package repository

import (
	"context"

	"arm_back/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkoutRepository interface {
	ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Workout, error)
	GetByID(ctx context.Context, id uuid.UUID) (*model.Workout, error)
	Create(ctx context.Context, w *model.Workout) error
	Update(ctx context.Context, w *model.Workout) error
	SoftDelete(ctx context.Context, id uuid.UUID) error
	// Exercises
	AddExercise(ctx context.Context, we *model.WorkoutExercise) error
	UpdateExercise(ctx context.Context, we *model.WorkoutExercise) error
	SoftDeleteExercise(ctx context.Context, id uuid.UUID) error
	// Sets
	AddSet(ctx context.Context, s *model.WorkoutSet) error
	UpdateSet(ctx context.Context, s *model.WorkoutSet) error
	SoftDeleteSet(ctx context.Context, id uuid.UUID) error
	// Copy
	CopyWorkout(ctx context.Context, source *model.Workout, newID uuid.UUID, userID uuid.UUID) (*model.Workout, error)
}

type pgWorkoutRepo struct {
	pool *pgxpool.Pool
}

func NewWorkoutRepository(pool *pgxpool.Pool) WorkoutRepository {
	return &pgWorkoutRepo{pool: pool}
}

func (r *pgWorkoutRepo) ListByUser(ctx context.Context, userID uuid.UUID) ([]model.Workout, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, date, weekday, comment, is_deleted, created_at, updated_at
		 FROM workouts WHERE user_id=$1 AND is_deleted=false ORDER BY date`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workouts []model.Workout
	for rows.Next() {
		var w model.Workout
		if err := rows.Scan(&w.ID, &w.UserID, &w.Date, &w.Weekday, &w.Comment, &w.IsDeleted, &w.CreatedAt, &w.UpdatedAt); err != nil {
			return nil, err
		}
		workouts = append(workouts, w)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Eager load exercises and sets
	for i := range workouts {
		exs, err := r.loadExercises(ctx, workouts[i].ID)
		if err != nil {
			return nil, err
		}
		workouts[i].Exercises = exs
	}
	return workouts, nil
}

func (r *pgWorkoutRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Workout, error) {
	var w model.Workout
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, date, weekday, comment, is_deleted, created_at, updated_at
		 FROM workouts WHERE id=$1`, id,
	).Scan(&w.ID, &w.UserID, &w.Date, &w.Weekday, &w.Comment, &w.IsDeleted, &w.CreatedAt, &w.UpdatedAt)
	if err == pgx.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	exs, err := r.loadExercises(ctx, w.ID)
	if err != nil {
		return nil, err
	}
	w.Exercises = exs
	return &w, nil
}

func (r *pgWorkoutRepo) Create(ctx context.Context, w *model.Workout) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO workouts (id, user_id, date, weekday, comment) VALUES ($1, $2, $3, $4, $5)`,
		w.ID, w.UserID, w.Date, w.Weekday, w.Comment,
	)
	return err
}

func (r *pgWorkoutRepo) Update(ctx context.Context, w *model.Workout) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE workouts SET date=$2, weekday=$3, comment=$4, updated_at=now() WHERE id=$1`,
		w.ID, w.Date, w.Weekday, w.Comment,
	)
	return err
}

func (r *pgWorkoutRepo) SoftDelete(ctx context.Context, id uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Delete sets of exercises in this workout
	_, err = tx.Exec(ctx,
		`UPDATE workout_sets SET is_deleted=true, updated_at=now()
		 WHERE workout_exercise_id IN (SELECT id FROM workout_exercises WHERE workout_id=$1)`, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE workout_exercises SET is_deleted=true, updated_at=now() WHERE workout_id=$1`, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE workouts SET is_deleted=true, updated_at=now() WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// Exercises

func (r *pgWorkoutRepo) AddExercise(ctx context.Context, we *model.WorkoutExercise) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO workout_exercises (id, workout_id, exercise_id, name, sort_order, comment)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		we.ID, we.WorkoutID, we.ExerciseID, we.Name, we.SortOrder, we.Comment,
	)
	return err
}

func (r *pgWorkoutRepo) UpdateExercise(ctx context.Context, we *model.WorkoutExercise) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE workout_exercises SET name=$2, sort_order=$3, comment=$4, updated_at=now() WHERE id=$1`,
		we.ID, we.Name, we.SortOrder, we.Comment,
	)
	return err
}

func (r *pgWorkoutRepo) SoftDeleteExercise(ctx context.Context, id uuid.UUID) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `UPDATE workout_sets SET is_deleted=true, updated_at=now() WHERE workout_exercise_id=$1`, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `UPDATE workout_exercises SET is_deleted=true, updated_at=now() WHERE id=$1`, id)
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// Sets

func (r *pgWorkoutRepo) AddSet(ctx context.Context, s *model.WorkoutSet) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO workout_sets (id, workout_exercise_id, set_number, weight, reps, to_failure)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		s.ID, s.WorkoutExerciseID, s.SetNumber, s.Weight, s.Reps, s.ToFailure,
	)
	return err
}

func (r *pgWorkoutRepo) UpdateSet(ctx context.Context, s *model.WorkoutSet) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE workout_sets SET set_number=$2, weight=$3, reps=$4, to_failure=$5, updated_at=now() WHERE id=$1`,
		s.ID, s.SetNumber, s.Weight, s.Reps, s.ToFailure,
	)
	return err
}

func (r *pgWorkoutRepo) SoftDeleteSet(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE workout_sets SET is_deleted=true, updated_at=now() WHERE id=$1`, id)
	return err
}

// Copy

func (r *pgWorkoutRepo) CopyWorkout(ctx context.Context, source *model.Workout, newID uuid.UUID, userID uuid.UUID) (*model.Workout, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	// Create new workout
	_, err = tx.Exec(ctx,
		`INSERT INTO workouts (id, user_id, date, weekday, comment) VALUES ($1, $2, $3, $4, $5)`,
		newID, userID, source.Date, source.Weekday, source.Comment,
	)
	if err != nil {
		return nil, err
	}

	// Copy exercises and sets
	for _, ex := range source.Exercises {
		newExID := uuid.New()
		_, err = tx.Exec(ctx,
			`INSERT INTO workout_exercises (id, workout_id, exercise_id, name, sort_order, comment)
			 VALUES ($1, $2, $3, $4, $5, $6)`,
			newExID, newID, ex.ExerciseID, ex.Name, ex.SortOrder, ex.Comment,
		)
		if err != nil {
			return nil, err
		}

		for _, s := range ex.Sets {
			_, err = tx.Exec(ctx,
				`INSERT INTO workout_sets (id, workout_exercise_id, set_number, weight, reps, to_failure)
				 VALUES ($1, $2, $3, $4, $5, $6)`,
				uuid.New(), newExID, s.SetNumber, s.Weight, s.Reps, s.ToFailure,
			)
			if err != nil {
				return nil, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, newID)
}

// helpers

func (r *pgWorkoutRepo) loadExercises(ctx context.Context, workoutID uuid.UUID) ([]model.WorkoutExercise, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, workout_id, exercise_id, name, sort_order, comment, is_deleted, created_at, updated_at
		 FROM workout_exercises WHERE workout_id=$1 AND is_deleted=false ORDER BY sort_order`, workoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exs []model.WorkoutExercise
	for rows.Next() {
		var we model.WorkoutExercise
		if err := rows.Scan(&we.ID, &we.WorkoutID, &we.ExerciseID, &we.Name, &we.SortOrder, &we.Comment, &we.IsDeleted, &we.CreatedAt, &we.UpdatedAt); err != nil {
			return nil, err
		}
		exs = append(exs, we)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range exs {
		sets, err := r.loadSets(ctx, exs[i].ID)
		if err != nil {
			return nil, err
		}
		exs[i].Sets = sets
	}
	return exs, nil
}

func (r *pgWorkoutRepo) loadSets(ctx context.Context, weID uuid.UUID) ([]model.WorkoutSet, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, workout_exercise_id, set_number, weight, reps, to_failure, is_deleted, created_at, updated_at
		 FROM workout_sets WHERE workout_exercise_id=$1 AND is_deleted=false ORDER BY set_number`, weID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sets []model.WorkoutSet
	for rows.Next() {
		var s model.WorkoutSet
		if err := rows.Scan(&s.ID, &s.WorkoutExerciseID, &s.SetNumber, &s.Weight, &s.Reps, &s.ToFailure, &s.IsDeleted, &s.CreatedAt, &s.UpdatedAt); err != nil {
			return nil, err
		}
		sets = append(sets, s)
	}
	return sets, rows.Err()
}
