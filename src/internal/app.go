package app

import (
	"log"

	"film-library/src/internal/config"
	"film-library/src/internal/db"
	"film-library/src/internal/film"
	"film-library/src/internal/models"
	"film-library/src/internal/router"
	"film-library/src/internal/user"
)

type App struct {
	Router *router.Router
	Config *config.Config
}

func NewApp(cfg *config.Config) *App {
	database, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatal(err)
	}

	userRepo := user.NewRepository(database.GetDB())
	userService := user.NewService(userRepo, cfg)
	userHandler := user.NewHandler(userService)

	actorRepo := models.NewRepository(database.GetDB())
	actorService := models.NewService(actorRepo)
	actorHandler := models.NewHandler(actorService)

	filmRepo := film.NewRepository(database.GetDB())
	filmService := film.NewService(filmRepo)
	filmHandler := film.NewHandler(filmService)

	router := router.NewRouter(cfg, userHandler, actorHandler, filmHandler)

	return &App{
		Router: router,
		Config: cfg,
	}
}

func (a *App) Run() {
	log.Printf("server running %s", a.Config.Addr())
	if err := a.Router.Run(a.Config.Addr()); err != nil {
		log.Fatal(err)
	}
}
