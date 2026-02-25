package handler

import (
	"errors"
	"net/http"

	"arm_back/internal/middleware"
	"arm_back/internal/model"
	"arm_back/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ExerciseHandler struct {
	svc *service.ExerciseService
}

func NewExerciseHandler(svc *service.ExerciseService) *ExerciseHandler {
	return &ExerciseHandler{svc: svc}
}

func (h *ExerciseHandler) List(c *gin.Context) {
	userID := middleware.GetUserID(c)
	exercises, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "failed to list exercises"})
		return
	}
	if exercises == nil {
		exercises = []model.Exercise{}
	}
	c.JSON(http.StatusOK, exercises)
}

func (h *ExerciseHandler) GetByID(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}

	ex, err := h.svc.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, ex)
}

func (h *ExerciseHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req model.CreateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request", Message: err.Error()})
		return
	}

	ex, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "failed to create exercise"})
		return
	}
	c.JSON(http.StatusCreated, ex)
}

func (h *ExerciseHandler) Update(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}

	var req model.UpdateExerciseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request"})
		return
	}

	ex, err := h.svc.Update(c.Request.Context(), id, userID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, ex)
}

func (h *ExerciseHandler) Delete(c *gin.Context) {
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

// Comments

func (h *ExerciseHandler) ListComments(c *gin.Context) {
	userID := middleware.GetUserID(c)
	exerciseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}

	comments, err := h.svc.ListComments(c.Request.Context(), exerciseID, userID)
	if err != nil {
		handleError(c, err)
		return
	}
	if comments == nil {
		comments = []model.ExerciseComment{}
	}
	c.JSON(http.StatusOK, comments)
}

func (h *ExerciseHandler) CreateComment(c *gin.Context) {
	userID := middleware.GetUserID(c)
	exerciseID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid id"})
		return
	}

	var req model.CreateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request"})
		return
	}

	comment, err := h.svc.CreateComment(c.Request.Context(), exerciseID, userID, req)
	if err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, comment)
}

func (h *ExerciseHandler) UpdateComment(c *gin.Context) {
	userID := middleware.GetUserID(c)
	commentID, err := uuid.Parse(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid comment id"})
		return
	}

	var req model.UpdateCommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid request"})
		return
	}

	if err := h.svc.UpdateComment(c.Request.Context(), commentID, userID, req); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "updated"})
}

func (h *ExerciseHandler) DeleteComment(c *gin.Context) {
	userID := middleware.GetUserID(c)
	commentID, err := uuid.Parse(c.Param("commentId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid comment id"})
		return
	}

	if err := h.svc.DeleteComment(c.Request.Context(), commentID, userID); err != nil {
		handleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "deleted"})
}

func handleError(c *gin.Context, err error) {
	if errors.Is(err, model.ErrNotFound) {
		c.JSON(http.StatusNotFound, model.ErrorResponse{Error: "not found"})
		return
	}
	if errors.Is(err, model.ErrForbidden) {
		c.JSON(http.StatusForbidden, model.ErrorResponse{Error: "forbidden"})
		return
	}
	if errors.Is(err, model.ErrInvalidInput) {
		c.JSON(http.StatusBadRequest, model.ErrorResponse{Error: "invalid input"})
		return
	}
	c.JSON(http.StatusInternalServerError, model.ErrorResponse{Error: "internal error"})
}
