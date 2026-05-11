package handlers

import (
	"encoding/json"
	"go-api/internal/models"
	"net/http"
)

// GetBalance mengirimkan data saldo dalam format JSON
func GetBalance(w http.ResponseWriter, r *http.Request) {
	// 1. Ambil Query Param (Di Node: req.query.username)
	username := r.URL.Query().Get("username")

	// 2. Validasi Sederhana
	if username == "" {
		// Kirim error (Di Node: res.status(400).send("..."))
		http.Error(w, "Username tidak boleh kosong", http.StatusBadRequest)
		return
	}

	balance := models.CoinBalance{
		Username: username,
		Amount:   1000, // Ceritanya ambil dari DB
	}

	w.Header().Set("Content-Type", "application/json")
	
	// 3. Error Handling saat Encode
	// json.NewEncoder().Encode() mengembalikan 'error'
	err := json.NewEncoder(w).Encode(balance)
	if err != nil {
		// Jika gagal encode (sangat jarang, tapi harus ditangani)
		http.Error(w, "Gagal memproses data", http.StatusInternalServerError)
		return
	}
}

