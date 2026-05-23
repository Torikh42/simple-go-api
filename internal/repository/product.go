package repository

import "go-api/internal/models"

// 1. Interface: Kontrak yang harus dipenuhi oleh database apa pun
// Nanti jika kita ganti ke PostgreSQL, kodenya tetap sama!
type ProductRepository interface {
	GetAll() ([]models.Product, error)
	Create(product *models.Product) error
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

// 5. Implementasi Method Create
func (r *memoryProductRepo) Create(p *models.Product) error {
	p.ID = len(r.data) + 1
	r.data = append(r.data, *p)
	return nil
}
