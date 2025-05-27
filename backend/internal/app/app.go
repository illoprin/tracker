package app

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"tracker-backend/internal/app/dependencies"
	"tracker-backend/internal/config"
	"tracker-backend/internal/server"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type App struct {
	router *chi.Mux
	port   string
}

func NewApp(port string, deps *dependencies.Dependencies) *App {
	master := chi.NewRouter()

	// use some middleware stack
	master.Use(middleware.Logger)
	master.Use(middleware.RealIP)
	master.Use(middleware.Recoverer)

	// permit access to public images
	avatarsFS := http.FileServer(http.Dir(
		path.Join(os.Getenv(config.PublicDirPathEnvName), config.AvatarsDir),
	))
	coversFS := http.FileServer(http.Dir(
		path.Join(os.Getenv(config.PublicDirPathEnvName), config.CoversDir),
	))
	master.Handle("/public/avatars/*", http.StripPrefix("/public/avatars/", avatarsFS))
	master.Handle("/public/covers/*", http.StripPrefix("/public/covers/", coversFS))

	// mount api routes
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
