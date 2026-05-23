package handlers

import (
	"encoding/json"
	"go-api/internal/models"
	"go-api/internal/repository"
	"net/http"
)

// 1. Buat Struct Handler yang menampung Repository (Dependency Injection)
type ProductHandler struct {
	repo repository.ProductRepository
}

// 2. Constructor untuk Handler
func NewProductHandler(repo repository.ProductRepository) *ProductHandler {
	return &ProductHandler{
		repo: repo,
	}
}

// 3. Ubah fungsi biasa menjadi "Method" dari struct ProductHandler
// Perhatikan tambahan `(h *ProductHandler)` sebelum nama fungsi
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	// Panggil logika dari database/repository
	products, _ := h.repo.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var newProduct models.Product

	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		http.Error(w, "Data tidak valid", http.StatusBadRequest)
		return
	}

	// Panggil fungsi simpan ke database/repository
	// Kita passing pointer (&newProduct) agar ID-nya otomatis terisi oleh Repo
	h.repo.Create(&newProduct)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProduct)
}
