package models

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name" validate:"required"`
	Price int64  `json:"price" validate:"required,gt=0"`
	Stock int    `json:"stock" validate:"required,gte=0"`
}
