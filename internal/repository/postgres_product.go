package repository

import (
	"context"
	"errors"

	"go-api/internal/db"
	"go-api/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type postgresProductRepo struct {
	queries *db.Queries
}

// NewPostgresProductRepository membuat instance baru yang terkoneksi dengan database
func NewPostgresProductRepository(dbPool *pgxpool.Pool) ProductRepository {
	return &postgresProductRepo{
		queries: db.New(dbPool),
	}
}

// GetAll mengambil semua produk dari database
func (r *postgresProductRepo) GetAll() ([]models.Product, error) {
	dbProducts, err := r.queries.GetAllProducts(context.Background())
	if err != nil {
		return nil, err
	}

	var products []models.Product
	for _, p := range dbProducts {
		products = append(products, models.Product{
			ID:    int(p.ID),
			Name:  p.Name,
			Price: p.Price,
			Stock: int(p.Stock),
		})
	}

	return products, nil
}

// GetByID mengambil satu produk berdasarkan ID
func (r *postgresProductRepo) GetByID(id int) (*models.Product, error) {
	p, err := r.queries.GetProductByID(context.Background(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}

	return &models.Product{
		ID:    int(p.ID),
		Name:  p.Name,
		Price: p.Price,
		Stock: int(p.Stock),
	}, nil
}

// Create menyimpan produk baru ke database
func (r *postgresProductRepo) Create(product *models.Product) error {
	p, err := r.queries.CreateProduct(context.Background(), db.CreateProductParams{
		Name:  product.Name,
		Price: product.Price,
		Stock: int32(product.Stock),
	})
	if err != nil {
		return err
	}

	// Update ID produk dengan ID yang baru di-generate oleh database
	product.ID = int(p.ID)
	return nil
}

// Update mengubah data produk yang sudah ada
func (r *postgresProductRepo) Update(product *models.Product) error {
	_, err := r.queries.UpdateProduct(context.Background(), db.UpdateProductParams{
		ID:    int32(product.ID),
		Name:  product.Name,
		Price: product.Price,
		Stock: int32(product.Stock),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("produk tidak ditemukan")
		}
		return err
	}

	return nil
}

// Delete menghapus produk berdasarkan ID
func (r *postgresProductRepo) Delete(id int) error {
	err := r.queries.DeleteProduct(context.Background(), int32(id))
	return err
}
