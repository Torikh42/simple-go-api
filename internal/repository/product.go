package repository

import (
	"errors"
	"go-api/internal/models"
)

// 1. Interface: Kontrak yang harus dipenuhi oleh database apa pun
// Nanti jika kita ganti ke PostgreSQL, kodenya tetap sama!
type ProductRepository interface {
	GetAll() ([]models.Product, error)
	GetByID(id int) (*models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(id int) error
}

// 2. Struct Penyimpanan (In-Memory)
type memoryProductRepo struct {
	data []models.Product
}

// 3. Constructor (Fungsi Pembuat)
func NewMemoryProductRepository() ProductRepository {
	return &memoryProductRepo{
		data: []models.Product{
			{ID: 1, Name: "Kopi Hitam", Price: 15000, Stock: 100},
			{ID: 2, Name: "Susu Murni", Price: 20000, Stock: 50},
		},
	}
}

// 4. Implementasi Method GetAll
func (r *memoryProductRepo) GetAll() ([]models.Product, error) {
	return r.data, nil
}

// Implementasi Method GetByID
func (r *memoryProductRepo) GetByID(id int) (*models.Product, error) {
	for _, p := range r.data {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, errors.New("produk tidak ditemukan")
}

// 5. Implementasi Method Create
func (r *memoryProductRepo) Create(p *models.Product) error {
	p.ID = len(r.data) + 1
	r.data = append(r.data, *p)
	return nil
}

// Implementasi Method Update
func (r *memoryProductRepo) Update(p *models.Product) error {
	for i, existing := range r.data {
		if existing.ID == p.ID {
			r.data[i] = *p
			return nil
		}
	}
	return errors.New("produk tidak ditemukan")
}

// Implementasi Method Delete
func (r *memoryProductRepo) Delete(id int) error {
	for i, p := range r.data {
		if p.ID == id {
			// Menghapus elemen dari slice (array)
			r.data = append(r.data[:i], r.data[i+1:]...)
			return nil
		}
	}
	return errors.New("produk tidak ditemukan")
}
