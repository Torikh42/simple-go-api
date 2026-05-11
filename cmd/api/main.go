package main

import (
	"fmt"
	"go-api/internal/handlers"
	"go-api/internal/middleware"
	"net/http"
)

func main() {
	// Di Node.js: app = express()
	// Di Go 1.22+: Kita pakai ServeMux bawaan yang sudah mendukung Method & Path Params
	mux := http.NewServeMux()

	// Handler inline (seperti anonymous function di Express)
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		// Di Node: res.send("Hello World")
		fmt.Fprint(w, "Selamat datang di Go Backend!!")
	})

	mux.HandleFunc("GET /health", handlers.HealthHandler)

	mux.HandleFunc("GET /balance", handlers.GetBalance)

	fmt.Println("Server berjalan di http://localhost:8080")

	// Di Node: app.listen(8080) -> Non-blocking
	// Di Go: ListenAndServe itu BLOCKING (mengeblok thread utama)
	// Bungkus mux dengan middleware Logger
	handler := middleware.Logger(mux)
	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Printf("Gagal menyalakan server: %s\n", err)
	}
}
