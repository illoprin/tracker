package app

import (
	"fmt"
	"net/http"
	"tracker-backend/internal/app/setup"
	"tracker-backend/internal/server"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	router *chi.Mux
	port   string
}

func NewApp(port string, deps *setup.Dependencies) *App {
	master := chi.NewRouter()

	// use some middleware stack
	master.Use(middleware.Logger)
	master.Use(middleware.RealIP)
	master.Use(middleware.Recoverer)

	// TODO: add access to public data

	// mount user routes
	master.Mount("/api", server.NewAppRouter(deps))

	// TODO: generate and mount swagger docs

	app := &App{
		router: master,
		port:   port,
	}

	return app
}

func (a *App) Run() {
	fmt.Printf("server started on address %s\n", a.port)
	if err := http.ListenAndServe(":"+a.port, a.router); err != nil {
		fmt.Printf("error occurred %s\n", err.Error())
	}
}
