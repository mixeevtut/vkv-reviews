package handlers

import (
	"net/http"
	"strconv"
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

// GetMovie обрабатывает запрос на получение фильма по ID
// @Summary Получить фильм по ID
// @Description Возвращает информацию о конкретном фильме
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "ID фильма"
// @Success 200 {object} dto.ModelMovie
// @Failure 400
// @Failure 404
// @Failure 500
// @Router /video-{id} [get]
func (h HTTPMoviesHandler) GetMovie(c *gin.Context) {
	idStr := c.Request.URL.Path[len("/video-"):]
	ID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	movie, err := h.uc.GetMovie(c.Request.Context(), ID)
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

// GetMovies обрабатывает запрос на получение списка всех фильмов
// @Summary Получить список всех фильмов
// @Description Возвращает массив всех фильмов в системе
// @Tags Movies
// @Accept json
// @Produce json
// @Success 200 {object} dto.MoviesResponse "Успешный ответ"
// @Failure 500
// @Router /movies [get]
func (h HTTPMoviesHandler) GetMovies(c *gin.Context) {
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
