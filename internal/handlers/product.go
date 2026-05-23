package handlers

import (
	"encoding/json"
	"go-api/internal/models"
	"go-api/internal/repository"
	"net/http"
	"strconv"
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

// Menangani GET /products/{id}
func (h *ProductHandler) GetProductByID(w http.ResponseWriter, r *http.Request) {
	// Menangkap param {id} dari URL menggunakan fitur baru Go 1.22
	idString := r.PathValue("id")
	
	// URL params selalu berbentuk string, kita ubah jadi integer (angka)
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "ID harus berupa angka", http.StatusBadRequest)
		return
	}

	product, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(product)
}

// Menangani PUT /products/{id}
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "ID harus berupa angka", http.StatusBadRequest)
		return
	}

	var updatedProduct models.Product
	if err := json.NewDecoder(r.Body).Decode(&updatedProduct); err != nil {
		http.Error(w, "Data tidak valid", http.StatusBadRequest)
		return
	}
	
	// Pastikan ID di model sama dengan ID dari URL
	updatedProduct.ID = id

	if err := h.repo.Update(&updatedProduct); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProduct)
}

// Menangani DELETE /products/{id}
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idString := r.PathValue("id")
	id, err := strconv.Atoi(idString)
	if err != nil {
		http.Error(w, "ID harus berupa angka", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Mengembalikan 204 No Content untuk operasi Delete yang berhasil
	w.WriteHeader(http.StatusNoContent)
}
