package model

import (
	"time"

	"github.com/google/uuid"
)

type Habit struct {
	ID            uuid.UUID         `json:"id"`
	UserID        uuid.UUID         `json:"userId"`
	Title         string            `json:"title"`
	Color         string            `json:"color"`
	RepeatType    string            `json:"repeatType"`
	RepeatMode    string            `json:"repeatMode"`
	RepeatWeekday *int              `json:"repeatWeekday"`
	IsArchived    bool              `json:"isArchived"`
	IsDeleted     bool              `json:"isDeleted"`
	Completions   []HabitCompletion `json:"completions"`
	CreatedAt     time.Time         `json:"createdAt"`
	UpdatedAt     time.Time         `json:"updatedAt"`
}

type HabitCompletion struct {
	ID        uuid.UUID `json:"id"`
	HabitID   uuid.UUID `json:"habitId"`
	Date      time.Time `json:"date"`
	IsDeleted bool      `json:"isDeleted"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type CreateHabitRequest struct {
	ID            string `json:"id"`
	Title         string `json:"title" binding:"required"`
	Color         string `json:"color"`
	RepeatType    string `json:"repeatType"`
	RepeatMode    string `json:"repeatMode"`
	RepeatWeekday *int   `json:"repeatWeekday"`
	IsArchived    *bool  `json:"isArchived"`
	IsDeleted     *bool  `json:"isDeleted"`
	CreatedAt     string `json:"createdAt"`
}

type UpdateHabitRequest struct {
	Title         *string `json:"title"`
	Color         *string `json:"color"`
	RepeatType    *string `json:"repeatType"`
	RepeatMode    *string `json:"repeatMode"`
	RepeatWeekday *int    `json:"repeatWeekday"`
	IsArchived    *bool   `json:"isArchived"`
	IsDeleted     *bool   `json:"isDeleted"`
	CreatedAt     *string `json:"createdAt"`
}

type CreateHabitCompletionRequest struct {
	ID        string `json:"id"`
	Date      string `json:"date" binding:"required"`
	IsDeleted *bool  `json:"isDeleted"`
	CreatedAt string `json:"createdAt"`
}
