package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"vkvideo-reviews/internal/api/http/v1/dto"
	"vkvideo-reviews/internal/usecase"
)

type HTTPUserHandler struct {
	UseCase *usecase.UserUseCase
}

func NewUserHandler(uc *usecase.UserUseCase) *HTTPUserHandler {
	return &HTTPUserHandler{UseCase: uc}
}

func (h *HTTPUserHandler) CreateUser(c *gin.Context) {

	var model dto.ModelUser
	if err := c.ShouldBindJSON(&model); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := model.Validate(); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := model.ToUser()
	if err := h.UseCase.CreateUser(c.Request.Context(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "internal_error",
			"message": "Failed to create user",
		})
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully",
		"user": gin.H{
			"id":    model.ID,
			"name":  model.Name,
			"email": model.Email,
		},
	})
}
