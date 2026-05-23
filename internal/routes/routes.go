package routes

import (
	"fmt"
	"go-api/internal/handlers"
	"net/http"
)

// SetupRoutes mendaftarkan semua endpoint aplikasi
func SetupRoutes(mux *http.ServeMux, productHandler *handlers.ProductHandler) {
	// Rute Dasar
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Selamat datang di Go Backend!")
	})
	mux.HandleFunc("GET /health", handlers.HealthHandler)

	// Rute CRUD Product
	mux.HandleFunc("GET /products", productHandler.GetProducts)
	mux.HandleFunc("POST /products", productHandler.CreateProduct)
	mux.HandleFunc("GET /products/{id}", productHandler.GetProductByID)
	mux.HandleFunc("PUT /products/{id}", productHandler.UpdateProduct)
	mux.HandleFunc("DELETE /products/{id}", productHandler.DeleteProduct)
}
