package services

import (
	"context"
	"errors"
	"fmt"
	"go-api/internal/models"
	"go-api/internal/repository"
	"sync"
	"time"
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
	admins := []string{"budi@admin.com", "siti@admin.com", "joko@admin.com"}
	go func() {
		var wg sync.WaitGroup

		for _, admin := range admins {
			wg.Add(1)

			go func(adminEmail string) {
				defer wg.Done()

				fmt.Printf("Mulai mengirim email ke %s...\n", adminEmail)
				time.Sleep(2 * time.Second)
				fmt.Printf("✅ Terkirim ke %s!\n", adminEmail)
			}(admin) 
		}

		wg.Wait()
		fmt.Println("seluruh email notifikasi berhasil dikirim!")
	}()

	return s.repo.Create(ctx, product)
}


func (s *productService) Update(ctx context.Context, product *models.Product) error {
	return s.repo.Update(ctx, product)
}

func (s *productService) Delete(ctx context.Context, id int) error {
	return s.repo.Delete(ctx, id)
}

func SendEmailNotification(productName string, statusChan chan string) {
	fmt.Printf("Mulai mengirim email untuk produk: %s...\n", productName)
	
	time.Sleep(3 * time.Second) 
	
	statusChan <- fmt.Sprintf("Email sukses terkirim untuk produk: %s!", productName)
}

