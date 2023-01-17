package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"../models"
	"github.com/go-swagger/go-swagger/examples/oauth2/models"
	"github.com/go-swagger/go-swagger/examples/task-tracker/models"
	"github.com/gorilla/mux"
)

// KeyMovie is a key used for the Movie object in the context
type KeyMovie struct{}

// Movies handler for getting and updating movies
type Movies struct {
	l *log.Logger
	v *models.ValidationError
}

// NewMovies returns a new movies handler with the given logger
func NewMovies(l *log.Logger, v *models.Validation) *Movies {
	return &Movies{l, v}
}

// ErrInvalidMoviePath is an error message when the Movie path is not valid
var ErrInvalidMoviePath = fmt.Errorf("Invalid Path, path should be /Movies/[id]")

// GenericError is a generic error message returned by a server
type GenericError struct {
	Message string `json:"message"`
}

// ValidationError is a collection of validation error messages
type ValidationError struct {
	Messages []string `json:"messages"`
}

// getMovieID returns the Movie ID from the URL
// Panics if cannot convert the id into an integer
// this should never happen as the router ensures that
// this is a valid number
func getMovieID(r *http.Request) int {
	// parse the Movie id from the url
	vars := mux.Vars(r)

	// convert the id into an integer and return
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		// should never happen
		panic(err)
	}

	return id
}

// swagger:route PUT /Movies Movies updateMovie
// Update a Movies details
//
// responses:
//	201: noContentResponse
//  404: errorResponse
//  422: errorValidation

// Update handles PUT requests to update Movies
func (p *Movies) Update(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	// fetch the Movie from the context
	prod := r.Context().Value(KeyMovie{}).(models.Movie)
	p.l.Println("[DEBUG] updating record id", prod.ID)

	err := models.UpdateMovie(prod)
	if err == models.ErrMovieNotFound {
		p.l.Println("[ERROR] Movie not found", err)

		rw.WriteHeader(http.StatusNotFound)
		models.ToJSON(&GenericError{Message: "Movie not found in database"}, rw)
		return
	}

	// write the no content success header
	rw.WriteHeader(http.StatusNoContent)
}

// swagger:route GET /Movies Movies listMovies
// Return a list of Movies from the database
// responses:
//	200: MoviesResponse

// ListAll handles GET requests and returns all current Movies
func (p *Movies) ListAll(rw http.ResponseWriter, r *http.Request) {
	p.l.Println("[DEBUG] get all records")
	rw.Header().Add("Content-Type", "application/json")

	prods := models.GetMovies()

	err := models.ToJSON(prods, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Println("[ERROR] serializing Movie", err)
	}
}

// swagger:route GET /Movies/{id} Movies listSingleMovie
// Return a list of Movies from the database
// responses:
//	200: MovieResponse
//	404: errorResponse

// ListSingle handles GET requests
func (p *Movies) ListSingle(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	id := getMovieID(r)

	p.l.Println("[DEBUG] get record id", id)

	prod, err := models.GetMovieByID(id)

	switch err {
	case nil:

	case models.ErrMovieNotFound:
		p.l.Println("[ERROR] fetching Movie", err)

		rw.WriteHeader(http.StatusNotFound)
		models.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	default:
		p.l.Println("[ERROR] fetching Movie", err)

		rw.WriteHeader(http.StatusInternalServerError)
		models.ToJSON(&GenericError{Message: err.Error()}, rw)
		return
	}

	err = models.ToJSON(prod, rw)
	if err != nil {
		// we should never be here but log the error just incase
		p.l.Println("[ERROR] serializing Movie", err)
	}
}

// swagger:route POST /Movies Movies createMovie
// Create a new Movie
//
// responses:
//	200: MovieResponse
//  422: errorValidation
//  501: errorResponse

// Create handles POST requests to add new Movies
func (p *Movies) Create(rw http.ResponseWriter, r *http.Request) {
	// fetch the Movie from the context
	prod := r.Context().Value(KeyMovie{}).(models.Movie)

	p.l.Printf("[DEBUG] Inserting Movie: %#v\n", prod)
	models.AddMovie(prod)
}
