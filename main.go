package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/thoaidao89/moviemanagementsystem/handlers"
	"github.com/thoaidao89/moviemanagementsystem/models"

	gohandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	l := log.New(os.Stdout, "movie-management", log.LstdFlags)
	v := models.NewValidation()
	if err != nil {
		l.Fatal(err)
	}

	var bindAddress = os.Getenv("BIND_ADDRESS") + ":" + os.Getenv("PORT")
	ph := handlers.NewMovies(l, v)

	//create serve mux
	// create a new serve mux and register the handlers
	sm := mux.NewRouter()
	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/movies", ph.ListAll)
	getR.HandleFunc("/movies/{id:[0-9]+}", ph.ListSingle)

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/movies", ph.Update)
	putR.Use(ph.MiddlewareValidateMovie)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/movies", ph.Create)
	postR.Use(ph.MiddlewareValidateMovie)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/movies/{id:[0-9]+}", ph.Delete)

	//CORS
	ch := gohandlers.CORS(gohandlers.AllowedOrigins([]string{os.Getenv("ALLOW_CORS"), os.Getenv("ALLOW_CORS2")}))
	// create a new server
	s := http.Server{
		Addr:         bindAddress,       // configure the bind address
		Handler:      ch(sm),            // set the default handler
		ErrorLog:     l,                 // set the logger for the server
		ReadTimeout:  5 * time.Second,   // max time to read request from the client
		WriteTimeout: 10 * time.Second,  // max time to write response to the client
		IdleTimeout:  120 * time.Second, // max time for connections using TCP Keep-Alive
	}
	// connect with database
	models.ConnectDatabase()
	// start the server
	go func() {
		l.Println("Starting server on port ", os.Getenv("PORT"))

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
