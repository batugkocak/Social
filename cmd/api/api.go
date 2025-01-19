package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/batugkocak/social/docs"
	"github.com/batugkocak/social/internal/store"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// Application
type application struct {
	config config
	store  store.Storage
	logger *zap.SugaredLogger
}

func (app *application) run(mux http.Handler) error {

	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.config.apiURL
	docs.SwaggerInfo.BasePath = "/v1"

	srv := http.Server{
		Addr:         app.config.addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	app.logger.Infow("Server has started at", "addr", app.config.addr, "env", app.config.env)

	return srv.ListenAndServe()
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	docsURL := fmt.Sprintf("%s/swagger/doc.json", app.config.addr)
	// /v1
	r.Route("/v1", func(r chi.Router) {
		// v1/swagger
		r.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(docsURL),
		))
		// v1/swagger

		// v1/health
		r.Get("/health", app.healthCheckHandler)
		r.Route("/posts", func(r chi.Router) {
			r.Post("/", app.createPostHandler)
			r.Route("/{postID}", func(r chi.Router) {
				r.Use(app.postsContextMiddleware)
				r.Get("/", app.getPostHandler)
				r.Delete("/", app.deletePostHandler)
				r.Patch("/", app.patchPostHandler)
				r.Post("/comments", app.createCommentHandler)
			})
		})
		// v1/health

		// /v1/users
		r.Route("/users", func(r chi.Router) {

			r.Route("/{userID}", func(r chi.Router) {
				r.Use(app.usersContextMiddleware)
				r.Get("/", app.getUserHandler)
				r.Put("/follow", app.followUserHandler)
				r.Put("/unfollow", app.unfollowUserHandler)
			})
			r.Group(func(r chi.Router) {
				r.Get("/feed", app.getUserFeedHandler)
			})
		})
		// /v1/users

		// /v1/authentication
		r.Route("/authentication", func(r chi.Router) {
			// r.Post("/user", app.registerUserHandler)
			// r.Post("/token", app.createTokenHandler)
		})
		// /v1/authentication

	})
	// /v1

	return r
}

// Config
type config struct {
	addr     string
	dbConfig dbConfig
	env      string
	apiURL   string
}

// DB Config

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}
