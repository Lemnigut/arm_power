package handler

import (
	"net/http"

	"arm_back/internal/middleware"
	"arm_back/internal/model"
	"arm_back/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type HabitHandler struct {
	svc *service.HabitService
}

func NewHabitHandler(svc *service.HabitService) *HabitHandler {
	return &HabitHandler{svc: svc}
}

func (h *HabitHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	includeDeleted := c.Query("includeDeleted") == "true"
	habits, err := h.svc.List(c.Request.Context(), userID, includeDeleted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "failed to list habits"})
		return
	}
	if habits == nil {
		habits = []model.Habit{}
	}
	c.JSON(http.StatusOK, habits)
}

func (h *HabitHandler) GetByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}

	habit, err := h.svc.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, habit)
}

func (h *HabitHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req model.CreateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	habit, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, habit)
}

func (h *HabitHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}

	var req model.UpdateHabitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request"})
		return
	}

	habit, err := h.svc.Update(c.Request.Context(), id, userID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, habit)
}

func (h *HabitHandler) Delete(c *gin.Context) {
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

func (h *HabitHandler) CreateCompletion(c *gin.Context) {
	userID := middleware.GetUserID(c)
	habitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid habit id"})
		return
	}

	var req model.CreateHabitCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	habit, err := h.svc.CreateCompletion(c.Request.Context(), habitID, userID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, habit)
}

func (h *HabitHandler) DeleteCompletion(c *gin.Context) {
	userID := middleware.GetUserID(c)
	habitID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid habit id"})
		return
	}
	completionID, err := uuid.Parse(c.Param("completionId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid completion id"})
		return
	}

	if err := h.svc.DeleteCompletion(c.Request.Context(), habitID, completionID, userID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
