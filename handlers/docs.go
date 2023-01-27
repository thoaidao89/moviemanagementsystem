// Package classification of Movie API
//
// Documentation for Movie API
//
//	Schemes: http
//	BasePath: /
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
// swagger:meta
package handlers

import "github.com/thoaidao89/moviemanagementsystem/models"

//
// NOTE: Types defined here are purely for documentation purposes
// these types are not used by any of the handers

// Generic error message returned as a string
// swagger:response errorResponse
type errorResponseWrapper struct {
	// Description of the error
	// in: body
	Body GenericError
}

// Validation errors defined as an array of strings
// swagger:response errorValidation
type errorValidationWrapper struct {
	// Collection of the errors
	// in: body
	Body ValidationError
}

// A list of Movies
// swagger:response MoviesResponse
type MoviesResponseWrapper struct {
	// All current Movies
	// in: body
	Body []models.Movie
}

// Data structure representing a single Movie
// swagger:response MovieResponse
type MovieResponseWrapper struct {
	// Newly created Movie
	// in: body
	Body models.Movie
}

// No content is returned by this API endpoint
// swagger:response noContentResponse
type noContentResponseWrapper struct {
}

// swagger:parameters updateMovie createMovie
type MovieParamsWrapper struct {
	// Movie data structure to Update or Create.
	// Note: the id field is ignored by update and create operations
	// in: body
	// required: true
	Body models.Movie
}

// swagger:parameters listSingleMovie deleteMovie
type MovieIDParamsWrapper struct {
	// The id of the Movie for which the operation relates
	// in: path
	// required: true
	ID int `json:"id"`
}
