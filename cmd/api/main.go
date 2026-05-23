package main

import (
	"fmt"
	"go-api/internal/handlers"
	"go-api/internal/middleware"
	"go-api/internal/repository"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// 1. Inisialisasi Repository
	// Nanti kalau pakai Postgres, tinggal ganti baris ini jadi NewPostgresRepository()
	productRepo := repository.NewMemoryProductRepository()

	// 2. Inisialisasi Handler (Suntikkan / Inject repository ke dalam handler)
	productHandler := handlers.NewProductHandler(productRepo)

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Selamat datang di Go Backend!")
	})
	mux.HandleFunc("GET /health", handlers.HealthHandler)

	// 3. Gunakan method dari productHandler
	mux.HandleFunc("GET /products", productHandler.GetProducts)
	mux.HandleFunc("POST /products", productHandler.CreateProduct)

	// Fitur baru: Menangkap parameter ID dinamis
	mux.HandleFunc("GET /products/{id}", productHandler.GetProductByID)
	mux.HandleFunc("PUT /products/{id}", productHandler.UpdateProduct)
	mux.HandleFunc("DELETE /products/{id}", productHandler.DeleteProduct)

	fmt.Println("Server berjalan di http://localhost:8080")

	handler := middleware.Logger(mux)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Printf("Gagal menyalakan server: %s\n", err)
	}
}
