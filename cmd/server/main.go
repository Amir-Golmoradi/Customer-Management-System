package main

import (
	"context"
	"log"
	"net/http"

	"github.com/Amir-Golmoradi/Customer-Management-System/internal/config"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/customer"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/database"
	model "github.com/Amir-Golmoradi/Customer-Management-System/internal/database/generated"
	"github.com/Amir-Golmoradi/Customer-Management-System/internal/handler"
)

func main() {
	ctx := context.Background()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Config error", err)
	}

	// Create pgx connection pool
	pool, err := database.NewConnectionPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Config error", err)
	}
	defer pool.Close()

	queries := model.New(pool)

	initializeHandler(queries)
}

func initializeHandler(queries *model.Queries) {

	customerRepo := customer.NewCustomerRepository(queries)
	customerService := customer.NewService(customerRepo)
	customerHandler := handler.NewHandler(customerService)

	mux := http.NewServeMux()
	mux.HandleFunc("/customers", customerHandler.CreateCustomer)
	mux.HandleFunc("/customer", customerHandler.GetCustomers)
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Fatal("Running on port 8080 ", server.ListenAndServe())
}
