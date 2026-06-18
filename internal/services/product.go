package services

import (
	"context"
	"errors"
	"go-api/internal/models"
	"go-api/internal/repository"
)

type ProductService interface {
	GetAll(ctx context.Context) ([]models.Product, error)
	GetByID(ctx context.Context, id int) (*models.Product, error)
	Create(ctx context.Context, product *models.Product) error
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{
		repo: repo,
	}
}

func (s *productService) GetAll(ctx context.Context) ([]models.Product, error) {
	return s.repo.GetAll(ctx)
}

func (s *productService) GetByID(ctx context.Context, id int) (*models.Product, error) {
	if id <= 0 {
		return nil, errors.New("id tidak valid")
	}
	return s.repo.GetByID(ctx, id)
}

func (s *productService) Create(ctx context.Context, product *models.Product) error {
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
