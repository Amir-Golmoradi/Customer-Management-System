package main

import (
	"context"
	"log"

	"github.com/Amir-Golmoradi/Customer-Management-System/internal/config"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/database"
	model "github.com/Amir-Golmoradi/Customer-Management-System/internal/database/generated"
)

func main() {
	ctx := context.Background()
	cfg, _ := config.Load()

	// Create pgx connection pool
	pool, err := database.NewConnectionPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Config error", err)
	}
	defer pool.Close()

	queries := model.New(pool)
	_ = queries
}
