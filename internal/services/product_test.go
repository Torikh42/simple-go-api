package services_test

import (
	"context"
	"testing"

	"go-api/internal/repository"
	"go-api/internal/services"
)

func TestProductService_GetByID(t *testing.T) {
	mockRepo := repository.NewMemoryProductRepository()
	service := services.NewProductService(mockRepo)
	ctx := context.Background()

	tests := []struct {
		name        string
		productID   int
		expectError bool
	}{
		{
			name:        "ID valid dan produk ditemukan",
			productID:   1,
			expectError: false,
		},
		{
			name:        "ID tidak valid karena negatif",
			productID:   -5,
			expectError: true,
		},
		{
			name:        "ID valid tapi produk tidak ada di database",
			productID:   999,
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := service.GetByID(ctx, tt.productID)

			if tt.expectError {
				if err == nil {
					t.Errorf("Diharapkan gagal, tapi malah sukses mendapatkan produk: %v", product)
				}
			} else {
				if err != nil {
					t.Errorf("Diharapkan sukses, tapi mendapat error: %v", err)
				}
				if product == nil {
					t.Error("Produk tidak boleh kosong jika sukses")
				}
			}
		})
	}
}
