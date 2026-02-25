package model

import (
	"time"

	"github.com/google/uuid"
)

type Exercise struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	Name         string    `json:"name"`
	Muscles      []string  `json:"muscles"`
	Category     string    `json:"category"`
	Description  string    `json:"description"`
	YoutubeLinks []string  `json:"youtubeLinks"`
	IsDeleted    bool      `json:"isDeleted"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type ExerciseComment struct {
	ID         uuid.UUID `json:"id"`
	ExerciseID uuid.UUID `json:"exerciseId"`
	UserID     uuid.UUID `json:"userId"`
	Text       string    `json:"text"`
	IsDeleted  bool      `json:"isDeleted"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type CreateExerciseRequest struct {
	ID           string   `json:"id"`
	Name         string   `json:"name" binding:"required"`
	Muscles      []string `json:"muscles"`
	Category     string   `json:"category"`
	Description  string   `json:"description"`
	YoutubeLinks []string `json:"youtubeLinks"`
}

type UpdateExerciseRequest struct {
	Name         *string  `json:"name"`
	Muscles      []string `json:"muscles"`
	Category     *string  `json:"category"`
	Description  *string  `json:"description"`
	YoutubeLinks []string `json:"youtubeLinks"`
}

type CreateCommentRequest struct {
	ID   string `json:"id"`
	Text string `json:"text" binding:"required"`
}

type UpdateCommentRequest struct {
	Text string `json:"text" binding:"required"`
}
