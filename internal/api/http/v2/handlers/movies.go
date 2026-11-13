package handlers

import (
	"errors"
	"net/http"
	"vkvideo-reviews/internal/api/http/v2/dto"

	apigen "vkvideo-reviews/internal/api/gen"
	"vkvideo-reviews/internal/entity"
	"vkvideo-reviews/internal/usecase"

	"github.com/gin-gonic/gin"
)

type HTTPMoviesHandler struct {
	uc *usecase.Movies
}

func NewMoviesHandler(uc *usecase.Movies) *HTTPMoviesHandler {
	return &HTTPMoviesHandler{uc: uc}
}

func (h HTTPMoviesHandler) GetMovie(c *gin.Context, id apigen.IdPath) {

	movie, err := h.uc.GetMovie(c.Request.Context(), id)
	if errors.Is(err, entity.ErrMovieNotFound) {
		c.JSON(http.StatusNotFound, apigen.Error{Code: "not_found", Message: "Movie not found"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	movieDTO := dto.ToMovieResponse(*movie)
	c.JSON(http.StatusOK, movieDTO)
}

func (h HTTPMoviesHandler) GetMovies(c *gin.Context, params apigen.GetMoviesParams) {
	movies, err := h.uc.GetMovies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	moviesResponse := make([]apigen.Movie, 0, len(movies))
	for _, movie := range movies {
		moviesResponse = append(moviesResponse, dto.ToMovieResponse(*movie))
	}

	c.JSON(http.StatusOK, moviesResponse)
}
