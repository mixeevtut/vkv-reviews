package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	apigen "vkvideo-reviews/internal/api/gen"
	"vkvideo-reviews/internal/api/http/v2/dto"
	"vkvideo-reviews/internal/usecase"
)

type HTTPUserHandler struct {
	uc *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *HTTPUserHandler {
	return &HTTPUserHandler{uc: uc}
}

func (h *HTTPUserHandler) CreateUser(c *gin.Context) {
	var model apigen.UserCreate

	if err := c.ShouldBindJSON(&model); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest,
			apigen.Error{Code: "bad_request", Message: "invalid request body"})
		return
	}

	user := dto.ToUserEntity(model)
	id, err := h.uc.CreateUser(c.Request.Context(), user)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError,
			apigen.Error{Code: "internal_server_error", Message: "Failed to create user"})
	}

	response := dto.ToUserCreateResponse(id)
	c.JSON(http.StatusCreated, response)
}
