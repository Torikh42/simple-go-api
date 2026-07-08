package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"go-api/internal/db"
	"go-api/internal/handlers"
	"go-api/internal/middleware"
	"go-api/internal/repository"
	"go-api/internal/routes"
	"go-api/internal/services"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Peringatan: File .env tidak ditemukan, menggunakan variabel environment sistem")
	}

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
		os.Exit(1)
	}
	defer dbPool.Close()

	if err := dbPool.Ping(context.Background()); err != nil {
		fmt.Printf("Database tidak dapat dijangkau: %s\n", err)
		os.Exit(1)
	}
	fmt.Println("✅ Koneksi database berhasil!")

	productRepo := repository.NewPostgresProductRepository(dbPool)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	// DI Auth
	queries := db.New(dbPool)
	userRepo := repository.NewPostgresUserRepository(queries)
	authService := services.NewAuthService(userRepo)
	authHandler := handlers.NewAuthHandler(authService)

	mux := http.NewServeMux()
	routes.SetupRoutes(mux, productHandler, authHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server berjalan di http://localhost:%s\n", port)

	handler := middleware.Recoverer(middleware.Logger(mux))
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		fmt.Printf("Gagal menyalakan server: %s\n", err)
	}
}
