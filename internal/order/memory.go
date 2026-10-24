package order

import (
	"context"
	"sync"
)

type MemoryRepository struct {
	mu     sync.RWMutex
	orders map[string]Order
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{orders: make(map[string]Order)}
}

func (r *MemoryRepository) Create(ctx context.Context, o Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders[o.ID] = o
	return nil
}

func (r *MemoryRepository) ByID(ctx context.Context, id string) (Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if order, exists := r.orders[id]; exists {
		return order, nil
	}

	return Order{}, ErrNotFound
}
