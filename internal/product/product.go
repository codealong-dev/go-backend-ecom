package product

import (
	"errors"
	"time"
)

var (
	ErrNotFound   = errors.New("product not found")
	ErrOutOfStock = errors.New("not enough stock")
)

// entity, model
type Product struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Price     int       `json:"price"`
	Stock     int       `json:"stock"`
	CreatedAt time.Time `json:"created_at"`
}
