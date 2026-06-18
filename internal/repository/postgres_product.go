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

func NewPostgresProductRepository(dbPool *pgxpool.Pool) ProductRepository {
	return &postgresProductRepo{
		queries: db.New(dbPool),
	}
}

func (r *postgresProductRepo) GetAll(ctx context.Context) ([]models.Product, error) {
	dbProducts, err := r.queries.GetAllProducts(ctx)
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

func (r *postgresProductRepo) GetByID(ctx context.Context, id int) (*models.Product, error) {
	p, err := r.queries.GetProductByID(ctx, int32(id))
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

func (r *postgresProductRepo) Create(ctx context.Context, product *models.Product) error {
	p, err := r.queries.CreateProduct(ctx, db.CreateProductParams{
		Name:  product.Name,
		Price: product.Price,
		Stock: int32(product.Stock),
	})
	if err != nil {
		return err
	}

	product.ID = int(p.ID)
	return nil
}

func (r *postgresProductRepo) Update(ctx context.Context, product *models.Product) error {
	_, err := r.queries.UpdateProduct(ctx, db.UpdateProductParams{
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

func (r *postgresProductRepo) Delete(ctx context.Context, id int) error {
	err := r.queries.DeleteProduct(ctx, int32(id))
	return err
}
