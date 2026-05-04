package handlers

import (
	"encoding/json"
	"go-api/internal/models"
	"net/http"
)

// GetBalance mengirimkan data saldo dalam format JSON
func GetBalance(w http.ResponseWriter, r *http.Request) {
	// Inisialisasi data (Mock data)
	balance := models.CoinBalance{
		Username: "antigravity",
		Amount:   1000,
	}

	// Set Header agar klien tahu ini JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode struct ke JSON dan tulis langsung ke ResponseWriter
	json.NewEncoder(w).Encode(balance)
}
