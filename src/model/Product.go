package model

type Product struct {
	ProductId   string  `json:"product-id"`
	ProductName string  `json:"product-name"`
	Price       float64 `json:"price"`
	Tag         string  `json:"tag"`
	Company     string  `json:"company"`
	Stock       int64   `json:"stock"`
}
