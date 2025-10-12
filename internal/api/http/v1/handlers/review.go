package handlers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"vkvideo-reviews/internal/api/http/v1/dto"
	"vkvideo-reviews/internal/entity"
	"vkvideo-reviews/internal/usecase"
)

type HTTPReviewHandler struct {
	Svc *usecase.ReviewUseCase
}

func NewReviewHandler(svc *usecase.ReviewUseCase) *HTTPReviewHandler {
	return &HTTPReviewHandler{Svc: svc}
}

// GetById
// @Summary Получение рецензии по ID
// @Param id path int true "ID отзыва"
// @Success 200 {object} dto.ReviewDTO
// @Produce json
// @Failure 400
// @Failure 500
// @Router /reviews/{id} [get]
func (h *HTTPReviewHandler) GetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		writeError(c, http.StatusBadRequest, "bad_request", "invalid_id")
	}

	review, err := h.Svc.GetReview(c.Request.Context(), id)
	if err != nil {
		mapDomainError(c, err)
	}

	c.JSON(http.StatusOK, dto.ToReviewResponse(review))
}

// POST /reviews
func (h *HTTPReviewHandler) Create(c *gin.Context) {
	//TODO
}

// PATCH /reviews/:id
func (h *HTTPReviewHandler) Update(c *gin.Context) {
	//TODO
}

// DELETE /reviews/:id
func (h *HTTPReviewHandler) Delete(c *gin.Context) {
	//TODO
}

func writeError(c *gin.Context, status int, code, msg string) {
	c.AbortWithStatusJSON(status, dto.ErrorReview{Error: msg, Code: code})
}

func mapDomainError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, entity.ErrReviewNotFound):
		writeError(c, http.StatusNotFound, "not_found", "review not found")
	default:
		writeError(c, http.StatusInternalServerError, "internal", "internal server error")
	}
}
