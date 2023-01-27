package models

import (
	"fmt"
)

// ErrMovieNotFound is an error raise when a movie can't found in the database
var ErrMovieNotFound = fmt.Errorf("Movie not found")

// Movie defines the structure for an API Movie
// swagger:model
type Movie struct {
	// the id for the Movie
	//
	// required: false
	// min: 1
	ID int `json:"id" gorm: "primary_key"` // Unique identifier for the Movie

	// the name for this poduct
	//
	// required: true
	// max length: 255
	Name string `json:"name" validate:"required"`

	// the description for this poduct
	//
	// required: false
	// max length: 10000
	Description string `json:"description"`

	// the price for the Movie
	//
	// required: true
	// min: 0.01
	Price float32 `json:"price"`

	// the link of the Movie
	// required: false
	Link string `json:"link" `

	// the Length of the Movie
	// required: false
	Length int `json:"length"`

	// the Date that create Movie
	// required: false
	DateCreate string `json:"datecreate"`

	// the Date that upload Movie
	// required: false
	CreateAt string `json:"createat"`

	// the Size of Movie
	//required: true
	Size int64 `json:"size"`

	// the real Path of Movie
	//required: false
	Path string `json:"path"`
}

// Movies defines a slice of Movie
type Movies []*Movie

// GetMovies returns all Movies from the database
func GetMovies() Movies {
	var ListMovies Movies
	DB.Find(&ListMovies)
	return ListMovies
}

// GetMovieByID returns a single Movie which matches the id from the
// database.
// If a Movie is not found this function returns a MovieNotFound error
func GetMovieByID(id int) (*Movie, error) {
	//i := findIndexByMovieID(id)

	var movie Movie
	err := DB.Where("id = ?", id).First(&movie).Error
	if err != nil {
		return nil, ErrMovieNotFound
	}
	return &movie, nil
}

// UpdateMovie replaces a Movie in the database with the given
// item.
// If a Movie with the given id does not exist in the database
// this function returns a MovieNotFound error
func UpdateMovie(p Movie) error {
	/* i := findIndexByMovieID(p.ID)
	if i == -1 {
		return ErrMovieNotFound
	}

	// update the Movie in the DB
	MovieList[i] = &p
	*/

	var movie Movie
	err := DB.Where("id = ?", p.ID).First(&movie).Error
	if err != nil {
		return ErrMovieNotFound
	}
	movie = p
	DB.Save(&p)
	return nil
}

// AddMovie adds a new Movie to the database
func AddMovie(p *Movie) {
	// get the next id in sequence
	// maxID := MovieList[len(MovieList)-1].ID
	// p.ID = maxID + 1
	// MovieList = append(MovieList, &p)
	DB.Create(p)
}

// DeleteMovie deletes a Movie from the database
func DeleteMovie(id int) error {
	// i := findIndexByMovieID(id)
	// if i == -1 {
	// 	return ErrMovieNotFound
	// }

	// MovieList = append(MovieList[:i], MovieList[i+1])
	var movie Movie
	err := DB.Where("id = ?", id).First(&movie)
	if err != nil {
		return ErrMovieNotFound
	}
	DB.Delete(&movie)
	return nil
}

/*
// findIndex finds the index of a Movie in the database
// returns -1 when no Movie can be found
func findIndexByMovieID(id int) int {
	for i, p := range MovieList {
		if p.ID == id {
			return i
		}
	}

	return -1
}

var MovieList = []*Movie{
	&Movie{
		ID:          1,
		Name:        "Latte",
		Description: "Frothy milky coffee",
		Price:       2.45,
		SKU:         "abc323",
		Link:        "http://test.link/demo1.mp4",
	},
	&Movie{
		ID:          2,
		Name:        "Esspresso",
		Description: "Short and strong coffee without milk",
		Price:       1.99,
		SKU:         "fjd34",
		Link:        "http://test.link/demo2.mp4",
	},
}

*/
