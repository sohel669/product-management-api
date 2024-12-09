package models

import "time"

type Product struct {
	ID             int       `json:"id" db:"id"`
	UserID         int       `json:"user_id" db:"user_id"`
	ProductName    string    `json:"product_name" db:"product_name"`
	ProductDescription string `json:"product_description" db:"product_description"`
	ProductImages  []string  `json:"product_images" db:"product_images"`
	ProductPrice   float64   `json:"product_price" db:"product_price"`
	CompressedImages []string `json:"compressed_images" db:"compressed_images"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
