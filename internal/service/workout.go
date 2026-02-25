package service

import (
	"context"
	"fmt"
	"time"

	"arm_back/internal/model"
	"arm_back/internal/repository"

	"github.com/google/uuid"
)

type WorkoutService struct {
	repo repository.WorkoutRepository
}

func NewWorkoutService(repo repository.WorkoutRepository) *WorkoutService {
	return &WorkoutService{repo: repo}
}

func (s *WorkoutService) List(ctx context.Context, userID uuid.UUID) ([]model.Workout, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *WorkoutService) GetByID(ctx context.Context, id, userID uuid.UUID) (*model.Workout, error) {
	w, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if w.UserID != userID {
		return nil, model.ErrForbidden
	}
	return w, nil
}

func (s *WorkoutService) Create(ctx context.Context, userID uuid.UUID, req model.CreateWorkoutRequest) (*model.Workout, error) {
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, model.ErrInvalidInput
	}

	id := uuid.New()
	if req.ID != "" {
		if parsed, err := uuid.Parse(req.ID); err == nil {
			id = parsed
		}
	}
	w := &model.Workout{
		ID:      id,
		UserID:  userID,
		Date:    date,
		Weekday: req.Weekday,
		Comment: req.Comment,
	}

	if err := s.repo.Create(ctx, w); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, w.ID)
}

func (s *WorkoutService) Update(ctx context.Context, id, userID uuid.UUID, req model.UpdateWorkoutRequest) (*model.Workout, error) {
	w, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if w.UserID != userID {
		return nil, model.ErrForbidden
	}

	if req.Date != nil {
		d, err := time.Parse("2006-01-02", *req.Date)
		if err != nil {
			return nil, model.ErrInvalidInput
		}
		w.Date = d
	}
	if req.Weekday != nil {
		w.Weekday = *req.Weekday
	}
	if req.Comment != nil {
		w.Comment = *req.Comment
	}

	if err := s.repo.Update(ctx, w); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, w.ID)
}

func (s *WorkoutService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	w, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if w.UserID != userID {
		return model.ErrForbidden
	}
	return s.repo.SoftDelete(ctx, id)
}

// Exercises

func (s *WorkoutService) AddExercise(ctx context.Context, workoutID, userID uuid.UUID, req model.AddWorkoutExerciseRequest) (*model.Workout, error) {
	w, err := s.repo.GetByID(ctx, workoutID)
	if err != nil {
		return nil, err
	}
	if w.UserID != userID {
		return nil, model.ErrForbidden
	}

	exID, _ := uuid.Parse(req.ExerciseID)
	weID := uuid.New()
	if req.ID != "" {
		if parsed, err := uuid.Parse(req.ID); err == nil {
			weID = parsed
		}
	}
	we := &model.WorkoutExercise{
		ID:         weID,
		WorkoutID:  workoutID,
		ExerciseID: &exID,
		Name:       req.Name,
		SortOrder:  len(w.Exercises),
		Comment:    req.Comment,
	}

	if err := s.repo.AddExercise(ctx, we); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, workoutID)
}

func (s *WorkoutService) RemoveExercise(ctx context.Context, workoutID, exerciseID, userID uuid.UUID) error {
	w, err := s.repo.GetByID(ctx, workoutID)
	if err != nil {
		return err
	}
	if w.UserID != userID {
		return model.ErrForbidden
	}
	return s.repo.SoftDeleteExercise(ctx, exerciseID)
}

// Sets

func (s *WorkoutService) AddSet(ctx context.Context, workoutID, weID, userID uuid.UUID, req model.CreateSetRequest) (*model.Workout, error) {
	w, err := s.repo.GetByID(ctx, workoutID)
	if err != nil {
		return nil, err
	}
	if w.UserID != userID {
		return nil, model.ErrForbidden
	}

	// Find exercise to determine set number
	setNum := 1
	for _, ex := range w.Exercises {
		if ex.ID == weID {
			setNum = len(ex.Sets) + 1
			break
		}
	}

	setID := uuid.New()
	if req.ID != "" {
		if parsed, err := uuid.Parse(req.ID); err == nil {
			setID = parsed
		}
	}
	set := &model.WorkoutSet{
		ID:                setID,
		WorkoutExerciseID: weID,
		SetNumber:         setNum,
		Weight:            req.Weight,
		Reps:              req.Reps,
		ToFailure:         req.ToFailure,
	}

	if err := s.repo.AddSet(ctx, set); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, workoutID)
}

func (s *WorkoutService) UpdateSet(ctx context.Context, workoutID, setID, userID uuid.UUID, req model.UpdateSetRequest) (*model.Workout, error) {
	w, err := s.repo.GetByID(ctx, workoutID)
	if err != nil {
		return nil, err
	}
	if w.UserID != userID {
		return nil, model.ErrForbidden
	}

	// Find existing set
	var existing *model.WorkoutSet
	for _, ex := range w.Exercises {
		for _, st := range ex.Sets {
			if st.ID == setID {
				existing = &st
				break
			}
		}
	}
	if existing == nil {
		return nil, model.ErrNotFound
	}

	if req.Weight != nil {
		existing.Weight = *req.Weight
	}
	if req.Reps != nil {
		existing.Reps = *req.Reps
	}
	if req.ToFailure != nil {
		existing.ToFailure = *req.ToFailure
	}

	if err := s.repo.UpdateSet(ctx, existing); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, workoutID)
}

func (s *WorkoutService) DeleteSet(ctx context.Context, workoutID, setID, userID uuid.UUID) error {
	w, err := s.repo.GetByID(ctx, workoutID)
	if err != nil {
		return err
	}
	if w.UserID != userID {
		return model.ErrForbidden
	}
	return s.repo.SoftDeleteSet(ctx, setID)
}

// Copy

func (s *WorkoutService) Copy(ctx context.Context, workoutID, userID uuid.UUID) (*model.Workout, error) {
	source, err := s.repo.GetByID(ctx, workoutID)
	if err != nil {
		return nil, err
	}
	if source.UserID != userID {
		return nil, model.ErrForbidden
	}

	now := time.Now()
	weekdays := []string{"Воскресенье", "Понедельник", "Вторник", "Среда", "Четверг", "Пятница", "Суббота"}

	copy := &model.Workout{
		Date:    now,
		Weekday: weekdays[now.Weekday()],
		Comment: fmt.Sprintf("[Копия] %s", source.Date.Format("2006-01-02")),
		Exercises: source.Exercises,
	}

	return s.repo.CopyWorkout(ctx, copy, uuid.New(), userID)
}
