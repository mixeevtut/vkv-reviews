package handlers

import (
	"net/http"
	"vkvideo-reviews/internal/api/http/v2/dto"

	apigen "vkvideo-reviews/internal/api/gen"
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
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if movie == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
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

	var moviesResponse []apigen.Movie
	for _, movie := range movies {
		moviesResponse = append(moviesResponse, dto.ToMovieResponse(*movie))
	}

	c.JSON(http.StatusOK, moviesResponse)
}
