package handlers

import (
	"net/http"

	apigen "vkvideo-reviews/internal/api/gen"
	"vkvideo-reviews/internal/api/http/v1/dto"
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

	movieDTO := dto.FromMovie(*movie)
	c.JSON(http.StatusOK, movieDTO)
}

func (h HTTPMoviesHandler) GetMovies(c *gin.Context, params apigen.GetMoviesParams) {
	movies, err := h.uc.GetMovies(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var moviesDTO []*dto.ModelMovie
	for _, movie := range movies {
		moviesDTO = append(moviesDTO, dto.FromMovie(*movie))
	}

	response := dto.MoviesResponse{Movies: moviesDTO}
	c.JSON(http.StatusOK, response)
}
