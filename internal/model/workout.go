package model

import (
	"time"

	"github.com/google/uuid"
)

type Workout struct {
	ID          uuid.UUID         `json:"id"`
	UserID      uuid.UUID         `json:"userId"`
	Date        time.Time         `json:"date"`
	Weekday     string            `json:"weekday"`
	Comment     string            `json:"comment"`
	WorkoutType string            `json:"workoutType"`
	IsDeleted   bool              `json:"isDeleted"`
	Exercises   []WorkoutExercise `json:"exercises"`
	CreatedAt   time.Time         `json:"createdAt"`
	UpdatedAt   time.Time         `json:"updatedAt"`
}

type WorkoutExercise struct {
	ID           uuid.UUID    `json:"id"`
	WorkoutID    uuid.UUID    `json:"workoutId"`
	ExerciseID   *uuid.UUID   `json:"exerciseId"`
	Name         string       `json:"name"`
	SortOrder    int          `json:"sortOrder"`
	Comment      string       `json:"comment"`
	IsSingleHand bool         `json:"isSingleHand"`
	WeightUnit   string       `json:"weightUnit"`
	IsDeleted    bool         `json:"isDeleted"`
	Sets         []WorkoutSet `json:"sets"`
	CreatedAt    time.Time    `json:"createdAt"`
	UpdatedAt    time.Time    `json:"updatedAt"`
}

type WorkoutSet struct {
	ID                uuid.UUID `json:"id"`
	WorkoutExerciseID uuid.UUID `json:"workoutExerciseId"`
	SetNumber         int       `json:"setNumber"`
	Weight            float64   `json:"weight"`
	Reps              int       `json:"reps"`
	ToFailure         bool      `json:"toFailure"`
	Hand              string    `json:"hand"`
	IsDeleted         bool      `json:"isDeleted"`
	CreatedAt         time.Time `json:"createdAt"`
	UpdatedAt         time.Time `json:"updatedAt"`
}

type CreateWorkoutRequest struct {
	ID          string `json:"id"`
	Date        string `json:"date" binding:"required"`
	Weekday     string `json:"weekday"`
	Comment     string `json:"comment"`
	WorkoutType string `json:"workoutType"`
}

type UpdateWorkoutRequest struct {
	Date        *string `json:"date"`
	Weekday     *string `json:"weekday"`
	Comment     *string `json:"comment"`
	WorkoutType *string `json:"workoutType"`
}

type AddWorkoutExerciseRequest struct {
	ID           string `json:"id"`
	ExerciseID   string `json:"exerciseId" binding:"required"`
	Name         string `json:"name" binding:"required"`
	Comment      string `json:"comment"`
	IsSingleHand bool   `json:"isSingleHand"`
	WeightUnit   string `json:"weightUnit"`
}

type CreateSetRequest struct {
	ID        string  `json:"id"`
	Weight    float64 `json:"weight"`
	Reps      int     `json:"reps"`
	ToFailure bool    `json:"toFailure"`
	Hand      string  `json:"hand"`
}

type UpdateSetRequest struct {
	Weight    *float64 `json:"weight"`
	Reps      *int     `json:"reps"`
	ToFailure *bool    `json:"toFailure"`
	Hand      *string  `json:"hand"`
}
