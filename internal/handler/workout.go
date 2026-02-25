package handler

import (
	"net/http"

	"arm_back/internal/middleware"
	"arm_back/internal/model"
	"arm_back/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WorkoutHandler struct {
	svc *service.WorkoutService
}

func NewWorkoutHandler(svc *service.WorkoutService) *WorkoutHandler {
	return &WorkoutHandler{svc: svc}
}

func (h *WorkoutHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	workouts, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "failed to list workouts"})
		return
	}
	if workouts == nil {
		workouts = []model.Workout{}
	}
	c.JSON(http.StatusOK, workouts)
}

func (h *WorkoutHandler) GetByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}

	w, err := h.svc.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, w)
}

func (h *WorkoutHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req model.CreateWorkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	w, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, w)
}

func (h *WorkoutHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}

	var req model.UpdateWorkoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request"})
		return
	}

	w, err := h.svc.Update(c.Request.Context(), id, userID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, w)
}

func (h *WorkoutHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id, userID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// Exercises in workout

func (h *WorkoutHandler) AddExercise(c *gin.Context) {
	userID := middleware.GetUserID(c)
	workoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid workout id"})
		return
	}

	var req model.AddWorkoutExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	w, err := h.svc.AddExercise(c.Request.Context(), workoutID, userID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, w)
}

func (h *WorkoutHandler) RemoveExercise(c *gin.Context) {
	userID := middleware.GetUserID(c)
	workoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid workout id"})
		return
	}
	exerciseID, err := uuid.Parse(c.Param("exerciseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid exercise id"})
		return
	}

	if err := h.svc.RemoveExercise(c.Request.Context(), workoutID, exerciseID, userID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

// Sets

func (h *WorkoutHandler) AddSet(c *gin.Context) {
	userID := middleware.GetUserID(c)
	workoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid workout id"})
		return
	}
	weID, err := uuid.Parse(c.Param("exerciseId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid exercise id"})
		return
	}

	var req model.CreateSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request"})
		return
	}

	w, err := h.svc.AddSet(c.Request.Context(), workoutID, weID, userID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, w)
}

func (h *WorkoutHandler) UpdateSet(c *gin.Context) {
	userID := middleware.GetUserID(c)
	workoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid workout id"})
		return
	}
	setID, err := uuid.Parse(c.Param("setId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid set id"})
		return
	}

	var req model.UpdateSetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request"})
		return
	}

	w, err := h.svc.UpdateSet(c.Request.Context(), workoutID, setID, userID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, w)
}

func (h *WorkoutHandler) DeleteSet(c *gin.Context) {
	userID := middleware.GetUserID(c)
	workoutID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid workout id"})
		return
	}
	setID, err := uuid.Parse(c.Param("setId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid set id"})
		return
	}

	if err := h.svc.DeleteSet(c.Request.Context(), workoutID, setID, userID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func (h *WorkoutHandler) Copy(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}

	w, err := h.svc.Copy(c.Request.Context(), id, userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, w)
}
