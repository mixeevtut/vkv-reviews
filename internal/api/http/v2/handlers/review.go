package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	apigen "vkvideo-reviews/internal/api/gen"
	"vkvideo-reviews/internal/api/http/v2/dto"
	"vkvideo-reviews/internal/entity"
	"vkvideo-reviews/internal/usecase"
)

type HTTPReviewHandler struct {
	uc *usecase.ReviewUseCase
}

func NewReviewHandler(uc *usecase.ReviewUseCase) *HTTPReviewHandler {
	return &HTTPReviewHandler{uc: uc}
}

func (h *HTTPReviewHandler) GetReviewById(c *gin.Context, id apigen.IdPath) {
	r, err := h.uc.GetReview(c.Request.Context(), id)
	if err != nil {
		writeReviewError(c, err)
		return
	}
	c.JSON(http.StatusOK, dto.ToGetReviewResponse(r))
}

func (h *HTTPReviewHandler) CreateReview(c *gin.Context) {
	var in apigen.ReviewCreate
	if err := c.ShouldBindJSON(&in); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, apigen.Error{Code: "bad_request", Message: "invalid body"})
		return
	}
	rw := dto.CreateReviewToEntity(in)
	r, err := h.uc.CreateReview(c.Request.Context(), rw)
	if err != nil {
		//TODO нужен обработчик ошибок
		c.AbortWithStatusJSON(http.StatusInternalServerError, apigen.Error{Code: "err", Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.ToCreateReviewResponse(r))
}

func (h *HTTPReviewHandler) DeleteReview(c *gin.Context, id apigen.IdPath) {
	err := h.uc.DeleteReview(c.Request.Context(), id)
	if err != nil {
		writeReviewError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *HTTPReviewHandler) UpdateReview(c *gin.Context, id apigen.IdPath) {
	var in apigen.ReviewUpdate
	if err := c.ShouldBindJSON(&in); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, apigen.Error{Code: "bad_request", Message: "invalid body"})
		return
	}
	rw := dto.UpdateReviewToEntity(id, in)
	r, err := h.uc.UpdateReview(c.Request.Context(), rw)
	if err != nil {
		writeReviewError(c, err)
		return
	}

	c.JSON(http.StatusOK, dto.ToUpdateReviewResponse(r))
}

func writeReviewError(c *gin.Context, err error) {
	if errors.Is(err, entity.ErrReviewNotFound) {
		c.AbortWithStatusJSON(http.StatusNotFound, apigen.Error{Code: "not_found", Message: "Review not found"})
		return
	}
	c.AbortWithStatusJSON(http.StatusInternalServerError, apigen.Error{Code: "internal_server_error", Message: "Review operation failed"})
}
