package service

import (
	"context"
	"time"

	"arm_back/internal/model"
	"arm_back/internal/repository"

	"github.com/google/uuid"
)

type HabitService struct {
	repo repository.HabitRepository
}

func NewHabitService(repo repository.HabitRepository) *HabitService {
	return &HabitService{repo: repo}
}

func (s *HabitService) List(ctx context.Context, userID uuid.UUID) ([]model.Habit, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *HabitService) GetByID(ctx context.Context, id, userID uuid.UUID) (*model.Habit, error) {
	h, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if h.UserID != userID {
		return nil, model.ErrForbidden
	}
	return h, nil
}

func (s *HabitService) Create(ctx context.Context, userID uuid.UUID, req model.CreateHabitRequest) (*model.Habit, error) {
	id := uuid.New()
	if req.ID != "" {
		if parsed, err := uuid.Parse(req.ID); err == nil {
			id = parsed
		}
	}

	repeatType, repeatMode, repeatWeekday, err := normalizeHabitRepeat(req.RepeatType, req.RepeatMode, req.RepeatWeekday)
	if err != nil {
		return nil, err
	}
	color := req.Color
	if color == "" {
		color = "#E0A018"
	}

	h := &model.Habit{
		ID:            id,
		UserID:        userID,
		Title:         req.Title,
		Color:         color,
		RepeatType:    repeatType,
		RepeatMode:    repeatMode,
		RepeatWeekday: repeatWeekday,
	}
	if err := s.repo.Create(ctx, h); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, h.ID)
}

func (s *HabitService) Update(ctx context.Context, id, userID uuid.UUID, req model.UpdateHabitRequest) (*model.Habit, error) {
	h, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if h.UserID != userID {
		return nil, model.ErrForbidden
	}

	if req.Title != nil {
		h.Title = *req.Title
	}
	if req.Color != nil {
		h.Color = *req.Color
	}

	repeatType := h.RepeatType
	repeatMode := h.RepeatMode
	repeatWeekday := h.RepeatWeekday
	if req.RepeatType != nil {
		repeatType = *req.RepeatType
	}
	if req.RepeatMode != nil {
		repeatMode = *req.RepeatMode
	}
	if req.RepeatWeekday != nil {
		repeatWeekday = req.RepeatWeekday
	}
	repeatType, repeatMode, repeatWeekday, err = normalizeHabitRepeat(repeatType, repeatMode, repeatWeekday)
	if err != nil {
		return nil, err
	}
	h.RepeatType = repeatType
	h.RepeatMode = repeatMode
	h.RepeatWeekday = repeatWeekday

	if err := s.repo.Update(ctx, h); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, h.ID)
}

func (s *HabitService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	h, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if h.UserID != userID {
		return model.ErrForbidden
	}
	return s.repo.SoftDelete(ctx, id)
}

func (s *HabitService) CreateCompletion(ctx context.Context, habitID, userID uuid.UUID, req model.CreateHabitCompletionRequest) (*model.Habit, error) {
	h, err := s.repo.GetByID(ctx, habitID)
	if err != nil {
		return nil, err
	}
	if h.UserID != userID {
		return nil, model.ErrForbidden
	}

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
	completion := &model.HabitCompletion{
		ID:      id,
		HabitID: habitID,
		Date:    date,
	}
	if err := s.repo.CreateCompletion(ctx, completion); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, habitID)
}

func (s *HabitService) DeleteCompletion(ctx context.Context, habitID, completionID, userID uuid.UUID) error {
	h, err := s.repo.GetByID(ctx, habitID)
	if err != nil {
		return err
	}
	if h.UserID != userID {
		return model.ErrForbidden
	}
	completion, err := s.repo.GetCompletionByID(ctx, completionID)
	if err != nil {
		return err
	}
	if completion.HabitID != habitID {
		return model.ErrForbidden
	}
	return s.repo.SoftDeleteCompletion(ctx, completionID)
}

func normalizeHabitRepeat(repeatType, repeatMode string, repeatWeekday *int) (string, string, *int, error) {
	if repeatType == "" {
		repeatType = "daily"
	}
	switch repeatType {
	case "daily":
		return repeatType, "", nil, nil
	case "weekly":
		if repeatMode == "" {
			repeatMode = "any_day"
		}
		switch repeatMode {
		case "any_day":
			return repeatType, repeatMode, nil, nil
		case "specific_day":
			if repeatWeekday == nil || *repeatWeekday < 1 || *repeatWeekday > 7 {
				return "", "", nil, model.ErrInvalidInput
			}
			return repeatType, repeatMode, repeatWeekday, nil
		default:
			return "", "", nil, model.ErrInvalidInput
		}
	case "monthly":
		return repeatType, "", nil, nil
	default:
		return "", "", nil, model.ErrInvalidInput
	}
}
