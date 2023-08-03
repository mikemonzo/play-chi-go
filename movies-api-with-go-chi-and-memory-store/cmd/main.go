package main

import (
	"context"
	"log"

	"github.com/mikemonzo/play-chi-go/movies-api-with-go-chi-and-memory-store/config"
	"github.com/mikemonzo/play-chi-go/movies-api-with-go-chi-and-memory-store/internal/api"
	"github.com/mikemonzo/play-chi-go/movies-api-with-go-chi-and-memory-store/internal/models"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	store := models.NewMemoryMoviesStore()
	server := api.NewServer(cfg.HTTPServer, store)
	server.Start(ctx)
}
