package main

import (
	"log"
	"net/http"
	"time"

	"codealong/internal/order"
	"codealong/internal/product"
	"codealong/internal/user"
	"context"
)

func main() {
	ctx := context.Background()

	userRepo := user.NewMemoryRepository()
	userSvc := user.NewService(userRepo)
	userHandler := user.NewHandler(userSvc)

	productRepo := product.NewMemoryRepository()
	productSvc := product.NewService(productRepo)
	productHandler := product.NewHandler(productSvc)

	orderRepo := order.NewMemoryRepository()
	orderSvc := order.NewService(orderRepo, productRepo, userRepo)
	ordersHandler := order.NewHandler(orderSvc)

	mux := http.NewServeMux()
	userHandler.Register(mux)
	productHandler.Register(mux)
	ordersHandler.Register(mux)

	// Seed initial state
	if err := product.Seed(ctx, productRepo); err != nil {
		log.Fatal(err)
	}
	if err := user.Seed(ctx, userRepo); err != nil {
		log.Fatal(err)
	}

	srv := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Println("Server is running...")

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start", "error", err)
	}
}
