package order

import (
	"context"

	"codealong/internal/product"
	"codealong/internal/user"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Repository interface {
	Create(context.Context, Order) error
	ByID(context.Context, string) (Order, error)
}

type Service struct {
	repo     Repository
	products product.Repository
	users    user.Repository
}

func NewService(repo Repository, products product.Repository, users user.Repository) *Service {
	return &Service{
		repo:     repo,
		products: products,
		users:    users,
	}
}

func (s *Service) PlaceOrder(ctx context.Context, req PlaceOrderRequest) (Order, error) {
	if req.CustomerID == "" {
		return Order{}, ErrCustomerRequired
	}

	if len(req.Items) == 0 {
		return Order{}, ErrNoItems
	}

	if _, err := s.users.ByID(ctx, req.CustomerID); err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return Order{}, ErrCustomerNotFound
		}
		return Order{}, err
	}

	items := []OrderItem{}
	quantities := map[string]int{}
	newID := uuid.New().String()

	for _, item := range req.Items {
		if item.Quantity < 1 {
			return Order{}, fmt.Errorf("%w: product %s", ErrInvalidQuantity, item.ProductID)
		}

		p, err := s.products.ByID(ctx, item.ProductID)
		if err != nil {
			if errors.Is(err, product.ErrNotFound) {
				return Order{}, fmt.Errorf("%w: %s", ErrProductNotFound, item.ProductID)
			}
			return Order{}, err
		}

		items = append(items, OrderItem{
			ID:        uuid.New().String(),
			OrderID:   newID,
			ProductID: p.ID,
			Price:     p.Price,
			Quantity:  item.Quantity,
			CreatedAt: time.Now().UTC(),
		})
		quantities[p.ID] += item.Quantity
	}

	// Reserves stock before creating an order
	if err := s.products.TakeStock(ctx, quantities); err != nil {
		return Order{}, err
	}

	order := Order{
		ID:         newID,
		CustomerID: req.CustomerID,
		Status:     StatusPending,
		Items:      items,
		CreatedAt:  time.Now().UTC(),
	}

	if err := s.repo.Create(ctx, order); err != nil {
		return Order{}, err
	}

	return order, nil
}

func (s *Service) Get(ctx context.Context, id string) (Order, error) {
	order, err := s.repo.ByID(ctx, id)
	if err != nil {
		return Order{}, err
	}

	return order, nil
}
