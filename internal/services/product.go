package services

import (
	"context"
	"errors"
	"go-api/internal/models"
	"go-api/internal/repository"
)

// 1. Kontrak (Interface) untuk Service
type ProductService interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
}

// 2. Struct yang menampung Repository
type productService struct {
	repo repository.ProductRepository
}

// 3. Constructor
func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

// 4. Implementasi Logika Bisnis
func (s *productService) GetAll(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetAll(ctx) // Langsung passing saja ke repo kalau tidak ada logika tambahan
}

func (s *productService) GetByID(ctx context.Context, id int) (*models.Product, error) {
	if id <= 0 {
		return nil, errors.New("id tidak valid")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *productService) Create(ctx context.Context, product *models.Product) error {
	// Ini contoh Logika Bisnis (Business Rule):
	if product.Price < 0 {
		return errors.New("harga tidak boleh negatif")
	}
	if product.Stock < 0 {
		return errors.New("stok tidak boleh negatif")
	}
	if product.Name == "" {
		return errors.New("nama produk wajib diisi")
	}

	return s.repo.Create(ctx, product)
}

func (s *productService) Update(ctx context.Context, product *models.Product) error {
	// Logika bisnis yang sama berlaku untuk update
	if product.Price < 0 {
		return errors.New("harga tidak boleh negatif")
	}
	if product.Stock < 0 {
		return errors.New("stok tidak boleh negatif")
	}
	if product.Name == "" {
		return errors.New("nama produk wajib diisi")
	}

	return s.repo.Update(ctx, product)
}

func (s *productService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}
