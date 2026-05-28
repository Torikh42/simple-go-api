package repository

import (
	"context"
	"errors"
	"go-api/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// pgxpool.Pool adalah koneksi database yang bisa dipakai bersama secara aman
// oleh banyak goroutine (Connection Pool).
type postgresProductRepo struct {
	db *pgxpool.Pool
}

// Constructor: Terima koneksi database, kembalikan ProductRepository
func NewPostgresProductRepository(db *pgxpool.Pool) ProductRepository {
	return &postgresProductRepo{db: db}
}

func (r *postgresProductRepo) GetAll() ([]models.Product, error) {
	// context.Background() = "Jalankan ini tanpa batas waktu / timeout"
	// Nanti kita akan ganti ini dengan ctx dari Request (Modul 5: Go Context)
	rows, err := r.db.Query(context.Background(), "SELECT id, name, price, stock FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close() // Pastikan rows selalu ditutup setelah selesai dipakai

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *postgresProductRepo) GetByID(id int) (*models.Product, error) {
	var p models.Product
	err := r.db.QueryRow(
		context.Background(),
		"SELECT id, name, price, stock FROM products WHERE id = $1",
		id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)

	if err != nil {
		// pgx.ErrNoRows setara dengan "404 Not Found" di level database
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("produk tidak ditemukan")
		}
		return nil, err
	}
	return &p, nil
}

func (r *postgresProductRepo) Create(product *models.Product) error {
	// RETURNING id digunakan untuk mengambil ID yang di-generate oleh database
	// dan langsung memasukkannya ke dalam struct product kita
	err := r.db.QueryRow(
		context.Background(),
		"INSERT INTO products (name, price, stock) VALUES ($1, $2, $3) RETURNING id",
		product.Name, product.Price, product.Stock,
	).Scan(&product.ID)
	return err
}

func (r *postgresProductRepo) Update(product *models.Product) error {
	result, err := r.db.Exec(
		context.Background(),
		"UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4",
		product.Name, product.Price, product.Stock, product.ID,
	)
	if err != nil {
		return err
	}
	// Cek apakah ada baris yang benar-benar diupdate
	if result.RowsAffected() == 0 {
		return errors.New("produk tidak ditemukan")
	}
	return nil
}

func (r *postgresProductRepo) Delete(id int) error {
	result, err := r.db.Exec(
		context.Background(),
		"DELETE FROM products WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return errors.New("produk tidak ditemukan")
	}
	return nil
}
