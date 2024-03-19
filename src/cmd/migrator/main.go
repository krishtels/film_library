package main

import (
	"log"

	"film-library/src/internal/config"
	"film-library/src/internal/db"
)

func main() {
	cfg := config.New()
	m := db.NewMigrator(cfg)

	log.Println("migrating")
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
}
