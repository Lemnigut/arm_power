package service

import (
	"context"
	"errors"

	"arm_back/internal/model"
	"arm_back/internal/repository"

	"github.com/google/uuid"
)

type ExerciseService struct {
	repo repository.ExerciseRepository
}

func NewExerciseService(repo repository.ExerciseRepository) *ExerciseService {
	return &ExerciseService{repo: repo}
}

func (s *ExerciseService) List(ctx context.Context, userID uuid.UUID) ([]model.Exercise, error) {
	return s.repo.ListByUser(ctx, userID)
}

func (s *ExerciseService) GetByID(ctx context.Context, id, userID uuid.UUID) (*model.Exercise, error) {
	ex, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ex.UserID != userID {
		return nil, model.ErrForbidden
	}
	return ex, nil
}

func (s *ExerciseService) Create(ctx context.Context, userID uuid.UUID, req model.CreateExerciseRequest) (*model.Exercise, error) {
	id := uuid.New()
	if req.ID != "" {
		if parsed, err := uuid.Parse(req.ID); err == nil {
			id = parsed
		}
	}
	ex := &model.Exercise{
		ID:           id,
		UserID:       userID,
		Name:         req.Name,
		Muscles:      req.Muscles,
		Category:     req.Category,
		Description:  req.Description,
		YoutubeLinks: req.YoutubeLinks,
	}
	if ex.Muscles == nil {
		ex.Muscles = []string{}
	}
	if ex.YoutubeLinks == nil {
		ex.YoutubeLinks = []string{}
	}
	if err := s.repo.Create(ctx, ex); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, ex.ID)
}

func (s *ExerciseService) Update(ctx context.Context, id, userID uuid.UUID, req model.UpdateExerciseRequest) (*model.Exercise, error) {
	ex, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if ex.UserID != userID {
		return nil, model.ErrForbidden
	}

	if req.Name != nil {
		ex.Name = *req.Name
	}
	if req.Muscles != nil {
		ex.Muscles = req.Muscles
	}
	if req.Category != nil {
		ex.Category = *req.Category
	}
	if req.Description != nil {
		ex.Description = *req.Description
	}
	if req.YoutubeLinks != nil {
		ex.YoutubeLinks = req.YoutubeLinks
	}

	if err := s.repo.Update(ctx, ex); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, ex.ID)
}

func (s *ExerciseService) Delete(ctx context.Context, id, userID uuid.UUID) error {
	ex, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if ex.UserID != userID {
		return model.ErrForbidden
	}
	return s.repo.SoftDelete(ctx, id)
}

// Comments

func (s *ExerciseService) ListComments(ctx context.Context, exerciseID, userID uuid.UUID) ([]model.ExerciseComment, error) {
	ex, err := s.repo.GetByID(ctx, exerciseID)
	if err != nil {
		return nil, err
	}
	if ex.UserID != userID {
		return nil, model.ErrForbidden
	}
	return s.repo.ListComments(ctx, exerciseID)
}

func (s *ExerciseService) CreateComment(ctx context.Context, exerciseID, userID uuid.UUID, req model.CreateCommentRequest) (*model.ExerciseComment, error) {
	ex, err := s.repo.GetByID(ctx, exerciseID)
	if err != nil {
		return nil, err
	}
	if ex.UserID != userID {
		return nil, model.ErrForbidden
	}

	id := uuid.New()
	if req.ID != "" {
		if parsed, err := uuid.Parse(req.ID); err == nil {
			id = parsed
		}
	}
	c := &model.ExerciseComment{
		ID:         id,
		ExerciseID: exerciseID,
		UserID:     userID,
		Text:       req.Text,
	}
	if err := s.repo.CreateComment(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *ExerciseService) UpdateComment(ctx context.Context, commentID, userID uuid.UUID, req model.UpdateCommentRequest) error {
	c := &model.ExerciseComment{ID: commentID, Text: req.Text}
	return s.repo.UpdateComment(ctx, c)
}

func (s *ExerciseService) DeleteComment(ctx context.Context, commentID, userID uuid.UUID) error {
	_ = userID // ownership check can be added
	return s.repo.SoftDeleteComment(ctx, commentID)
}

func IsNotFound(err error) bool {
	return errors.Is(err, model.ErrNotFound)
}

func IsForbidden(err error) bool {
	return errors.Is(err, model.ErrForbidden)
}
