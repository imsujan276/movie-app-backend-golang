package main

import (
	"backend/pkg/handlers"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
)

func Routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Recoverer)
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"Authorization", "Content-Type"},
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
	}).Handler)

	r.Get("/status", handlers.Repo.StatusHandler)

	r.Route("/v1", func(r chi.Router) {
		r.Post("/login", handlers.Repo.LoginHandler)
		r.With(paginate).Get("/movies", handlers.Repo.GetAllMovies)
		r.Get("/movie/{id:[0-9]+}", handlers.Repo.GetOneMovie)
		r.With(paginate).Get("/genres", handlers.Repo.GetAllGenres)
		r.Get("/movies/{genre_id:[0-9]+}", handlers.Repo.GetAllMoviesByGenre)
		// r.With(handlers.Repo.CheckToken).Post("/admin/movie/add", handlers.Repo.AddUpdateMovie)
		// r.Get("/admin/movie/delete/{id:[0-9]+}", handlers.Repo.DeleteMovie)

		r.Route("/admin", func(r chi.Router) {
			// securing /admin route with checktoken middleware
			r.Use(handlers.Repo.CheckToken)
			r.Post("/movie/add", handlers.Repo.AddUpdateMovie)
			r.Get("/movie/delete/{id:[0-9]+}", handlers.Repo.DeleteMovie)
		})

		r.Route("/graphql", func(r chi.Router) {
			r.Post("/", handlers.Repo.GetAllMoviesGraphQL)
		})
	})

	return r
}

// paginate is a stub, but very possible to implement middleware logic
// to handle the request params for handling a paginated request.
func paginate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// just a stub.. some ideas are to look at URL query params for something like
		// the page number, or the limit, and send a query cursor down the chain
		next.ServeHTTP(w, r)
	})
}
