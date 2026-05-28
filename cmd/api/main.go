package main

import (
	"context"
	"fmt"
	"go-api/internal/handlers"
	"go-api/internal/middleware"
	"go-api/internal/repository"
	"go-api/internal/routes"
	"go-api/internal/services"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	// 1. Load variabel dari file .env
	if err := godotenv.Load(); err != nil {
		fmt.Println("Peringatan: File .env tidak ditemukan, menggunakan variabel environment sistem")
	}

	// 2. Buat koneksi ke PostgreSQL menggunakan Connection Pool
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	dbPool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		fmt.Printf("Gagal terhubung ke database: %s\n", err)
		os.Exit(1) // Hentikan server jika database tidak bisa dikoneksikan
	}
	defer dbPool.Close()

	// Cek apakah koneksi benar-benar berhasil
	if err := dbPool.Ping(context.Background()); err != nil {
		fmt.Printf("Database tidak dapat dijangkau: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Koneksi database berhasil!")

	// 3. Inisialisasi lapisan-lapisan aplikasi (Dependency Injection)
	// HANYA BARIS INI yang berubah dari In-Memory ke PostgreSQL!
	productRepo := repository.NewPostgresProductRepository(dbPool)

	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// 4. Setup routing
	mux := http.NewServeMux()
	routes.SetupRoutes(mux, productHandler)

	// 5. Nyalakan server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server berjalan di http://localhost:%s\n", port)

	handler := middleware.Logger(mux)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		fmt.Printf("Gagal menyalakan server: %s\n", err)
	}
}
