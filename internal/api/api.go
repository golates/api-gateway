package api

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/golates/api-gateway/internal/handlers"
	custom_middleware "github.com/golates/api-gateway/internal/middlewares"
	"github.com/golates/api-gateway/pkg/config"
	"log"
	"net/http"
	"time"
)

type API struct {
	router *chi.Mux
	cfg    *config.Config
}

func NewAPI(cfg *config.Config) *API {
	r := chi.NewRouter()

	return &API{
		router: r,
		cfg:    cfg,
	}
}

func (api *API) RunServer() error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%v", api.cfg.ApiPort),
		Handler:      api.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	log.Println("Server started:", api.cfg.ApiPort)

	return srv.ListenAndServe()
}

func (api *API) SetupMiddlewares() {
	api.router.Use(middleware.Logger)
	api.router.Use(middleware.Recoverer)
	api.router.Use(custom_middleware.ValidatorMiddleware)
	api.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
}

func (api *API) SetupRoutes() {
	h := handlers.NewHandlers(api.cfg)

	api.router.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.AuthHandler.Login)
		r.Post("/login/oauth/google", h.AuthHandler.OAuthGoogleLogin)
		r.Post("/login/oauth/facebook", h.AuthHandler.OAuthFacebookLogin)
		r.Post("/check-email", h.AuthHandler.CheckEmail)
		r.Post("/register", h.AuthHandler.Register)
		r.Post("/forgot-password", h.AuthHandler.ForgotPassword)
	})
}
