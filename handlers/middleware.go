package handlers

import (
	"context"
	"net/http"

	"github.com/thoaidao89/moviemanagementsystem/models"
)

// MiddlewareValidateMovie validates the movie in the request and call next if ok

func (p *Movies) MiddlewareValidateMovie(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		rw.Header().Add("Content-Type", "application/json")
		mov := &models.Movie{}
		err := models.FromJSON(mov, r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing movie", err)
			rw.WriteHeader(http.StatusBadRequest)
			models.ToJSON(&GenericError{Message: err.Error()}, rw)
			return

		}
		p.l.Println("Movie ", mov)
		//validate the movie
		errs := p.v.Validate(mov)
		if len(errs) != 0 {
			p.l.Println("[ERROR] validating product", errs)

			//return validation message as array
			rw.WriteHeader(http.StatusUnprocessableEntity)
			models.ToJSON(&ValidationError{Messages: errs.Errors()}, rw)
			return
		}

		//add the product to the context
		ctx := context.WithValue(r.Context(), KeyMovie{}, mov)
		r = r.WithContext(ctx)

		// Call the next handler
		next.ServeHTTP(rw, r)

	})
}
