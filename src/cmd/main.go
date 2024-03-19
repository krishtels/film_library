package main

import (
	"film-library/src/internal"
	"film-library/src/internal/config"
)

func main() {
	cfg := config.New()
	app := app.NewApp(cfg)

	app.Run()
}
