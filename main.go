package main

import (
	"log"

	"eventhub/config"
	"eventhub/database"
	"eventhub/route"
)

func main() {
	cfg := config.Load()

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("gagal koneksi database: %v", err)
	}
	defer db.Close()

	app := config.NewApp()

	route.Setup(app)

	if err := app.Listen(cfg.ServerAddr()); err != nil {
		log.Fatalf("server berhenti: %v", err)
	}
}
