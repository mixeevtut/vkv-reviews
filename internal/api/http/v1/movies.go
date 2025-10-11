package v1

import (
	"net/http"
	"strconv"
	"vkvideo-reviews/internal/usecase"
)

type MoviesHandler struct {
	uc usecase.Movies
}

func (m MoviesHandler) GetMovies(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	movieID, err := strconv.Atoi(r.URL.Path[len("/video-"):])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movies, err := m.uc.GetMovie(r.Context(), movieID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
