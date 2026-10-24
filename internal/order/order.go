package order

import "time"
import "errors"

var (
	ErrCustomerRequired = errors.New("customer is mandatory")
	ErrNotFound         = errors.New("order not found")
	ErrProductNotFound  = errors.New("product not found")
	ErrNoItems          = errors.New("items are required")
	ErrCustomerNotFound = errors.New("customer does not exist")
	ErrInvalidQuantity  = errors.New("invalid quantity requested")
)

type Status string

const StatusPending Status = "pending"

// entity
type Order struct {
	ID         string      `json:"id"`
	Status     Status      `json:"status"`
	Items      []OrderItem `json:"items"`
	CustomerID string      `json:"customer_id"`
	CreatedAt  time.Time   `json:"created_at"`
}

type OrderItem struct {
	ID        string    `json:"id"`
	ProductID string    `json:"product_id"`
	OrderID   string    `json:"order_id"`
	Price     int       `json:"price"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
}

type OrderItemInput struct {
	ProductID string `json:"product_id"`
	Quantity  int    `json:"quantity"`
}

type PlaceOrderRequest struct {
	CustomerID string           `json:"customer_id"`
	Items      []OrderItemInput `json:"items"`
}
