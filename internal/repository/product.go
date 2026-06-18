package repository

import (
	"context"
	"errors"
	"go-api/internal/models"
)

type ProductRepository interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
}

type memoryProductRepo struct {
	data []models.Product
}

func NewMemoryProductRepository() ProductRepository {
	return &memoryProductRepo{
		data: []models.Product{
			{ID: 1, Name: "Kopi Hitam", Price: 15000, Stock: 100},
			{ID: 2, Name: "Susu Murni", Price: 20000, Stock: 50},
		},
	}
}

func (r *memoryProductRepo) GetAll(ctx context.Context) ([]models.Product, error) {
	return r.data, nil
}

func (r *memoryProductRepo) GetByID(ctx context.Context, id int) (*models.Product, error) {
	for _, p := range r.data {
		if p.ID == id {
			return &p, nil
		}
	}
	return nil, errors.New("produk tidak ditemukan")
}

func (r *memoryProductRepo) Create(ctx context.Context, p *models.Product) error {
	p.ID = len(r.data) + 1
	r.data = append(r.data, *p)
	return nil
}

func (r *memoryProductRepo) Update(ctx context.Context, p *models.Product) error {
	for i, existing := range r.data {
		if existing.ID == p.ID {
			r.data[i] = *p
			return nil
		}
	}
	return errors.New("produk tidak ditemukan")
}

func (r *memoryProductRepo) Delete(ctx context.Context, id int) error {
	for i, p := range r.data {
		if p.ID == id {
			r.data = append(r.data[:i], r.data[i+1:]...)
			return nil
		}
	}
	return errors.New("produk tidak ditemukan")
}
