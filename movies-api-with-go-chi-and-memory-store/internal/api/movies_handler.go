package api

import (
	"errors"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/mikemonzo/play-chi-go/movies-api-with-go-chi-and-memory-store/internal/models"
)

type movieResponse struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	ReleaseDate time.Time `json:"release_date"`
	TicketPrice float64   `json:"ticket_price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewMovieResponse(m models.Movie) movieResponse {
	return movieResponse{
		ID:          m.ID,
		Title:       m.Title,
		Director:    m.Director,
		ReleaseDate: m.ReleaseDate,
		TicketPrice: m.TicketPrice,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func (hr movieResponse) Render(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func NewMovieListResponse(movies []models.Movie) []render.Renderer {
	list := []render.Renderer{}
	for _, movie := range movies {
		mr := NewMovieResponse(movie)
		list = append(list, mr)
	}
	return list
}

func (s *Server) handleListMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := s.store.GetAll()
	if err != nil {
		render.Render(w, r, ErrInternalServerError)
		return
	}

	render.RenderList(w, r, NewMovieListResponse(movies))
}

func (s *Server) handleGetMovie(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	movie, err := s.store.GetByID(id)
	if err != nil {
		var rnfErr *models.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	mr := NewMovieResponse(movie)
	render.Render(w, r, mr)
}

type CreateMovieRequest struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	ReleaseDate time.Time `json:"release_date"`
	TicketPrice float64   `json:"ticket_price"`
}

func (mr *CreateMovieRequest) Bind(r *http.Request) error {
	return nil
}

func (s *Server) handleCreateMovie(w http.ResponseWriter, r *http.Request) {
	data := &CreateMovieRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	createMovieParams := models.CreateMovieParams{
		ID:          uuid.MustParse(data.ID),
		Title:       data.Title,
		Director:    data.Director,
		ReleaseDate: data.ReleaseDate,
		TicketPrice: data.TicketPrice,
	}
	err := s.store.Create(createMovieParams)
	if err != nil {
		var dupKeyErr *models.DuplicateKeyError
		if errors.As(err, &dupKeyErr) {
			render.Render(w, r, ErrConflict(err))
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	w.WriteHeader(200)
	w.Write(nil)
}

type updateMovieRequest struct {
	Title       string    `json:"title"`
	Director    string    `json:"director"`
	ReleaseDate time.Time `json:"release_date"`
	TicketPrice float64   `json:"ticket_price"`
}

func (mr *updateMovieRequest) Bind(r *http.Request) error {
	return nil
}

func (s *Server) handleUpdateMovie(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	data := &updateMovieRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	updateMovieParams := models.UpdateMovieParams{
		Title:       data.Title,
		Director:    data.Director,
		ReleaseDate: data.ReleaseDate,
		TicketPrice: data.TicketPrice,
	}
	err = s.store.Update(id, updateMovieParams)
	if err != nil {
		var rnfErr *models.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	w.WriteHeader(200)
	w.Write(nil)
}

func (s *Server) handleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		render.Render(w, r, ErrBadRequest)
		return
	}

	err = s.store.Delete(id)
	if err != nil {
		var rnfErr *models.RecordNotFoundError
		if errors.As(err, &rnfErr) {
			render.Render(w, r, ErrNotFound)
		} else {
			render.Render(w, r, ErrInternalServerError)
		}
		return
	}

	w.WriteHeader(200)
	w.Write(nil)
}
