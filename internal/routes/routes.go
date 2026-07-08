package routes

import (
	"fmt"
	"go-api/internal/handlers"
	"go-api/internal/middleware"
	"net/http"
)

func SetupRoutes(mux *http.ServeMux, productHandler *handlers.ProductHandler, authHandler *handlers.AuthHandler) {
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Selamat datang di Go Backend!")
	})
	mux.HandleFunc("GET /health", handlers.HealthHandler)

	// Auth Routes (Public)
	mux.HandleFunc("POST /register", authHandler.Register)
	mux.HandleFunc("POST /login", authHandler.Login)

	// Product Routes (Protected dengan middleware Auth)
	mux.HandleFunc("GET /products", middleware.Auth(productHandler.GetProducts))
	mux.HandleFunc("POST /products", middleware.Auth(productHandler.CreateProduct))
	mux.HandleFunc("GET /products/{id}", middleware.Auth(productHandler.GetProductByID))
	mux.HandleFunc("PUT /products/{id}", middleware.Auth(productHandler.UpdateProduct))
	mux.HandleFunc("DELETE /products/{id}", middleware.Auth(productHandler.DeleteProduct))
}
