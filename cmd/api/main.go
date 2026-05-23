package main

import (
	"fmt"
	"go-api/internal/handlers"
	"go-api/internal/middleware"
	"go-api/internal/repository"
	"go-api/internal/routes"
	"go-api/internal/services"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	// 1. Inisialisasi Repository
	// Nanti kalau pakai Postgres, tinggal ganti baris ini jadi NewPostgresRepository()
	productRepo := repository.NewMemoryProductRepository()

	// 2. Inisialisasi Service (Suntikkan Repo ke Service)
	productService := services.NewProductService(productRepo)

	// 3. Inisialisasi Handler (Suntikkan Service ke Handler)
	productHandler := handlers.NewProductHandler(productService)

	// 4. Daftarkan semua routes
	routes.SetupRoutes(mux, productHandler)

	fmt.Println("Server berjalan di http://localhost:8080")

	handler := middleware.Logger(mux)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Printf("Gagal menyalakan server: %s\n", err)
	}
}

