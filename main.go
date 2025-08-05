package main

import (
	"media-svc/config"
	"media-svc/internal/port"

	mongodb "github.com/dtome123/go-mongo-generic"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	db, err := mongodb.NewDatabase(
		mongodb.WithDatabase(cfg.DB.Mongo.Database),
		mongodb.WithSingleURL(cfg.DB.Mongo.DSN),
	)
	if err != nil {
		panic(err)
	}

	server := port.NewServer(cfg, db)
	server.Run()
}
