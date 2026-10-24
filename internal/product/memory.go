package product

import (
	"context"
	"fmt"
	"sort"
	"sync"
)

type MemoryRepository struct {
	mu       sync.RWMutex
	products map[string]Product
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{products: make(map[string]Product)}
}

func (r *MemoryRepository) Create(_ context.Context, p Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.products[p.ID] = p
	return nil
}

func (r *MemoryRepository) ByID(_ context.Context, id string) (Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.products[id]
	if !ok {
		return Product{}, ErrNotFound
	}
	return p, nil
}

func (r *MemoryRepository) List(_ context.Context) ([]Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]Product, 0, len(r.products))
	for _, p := range r.products {
		out = append(out, p)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}

func (r *MemoryRepository) TakeStock(_ context.Context, quantities map[string]int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for id, quantity := range quantities {
		p, ok := r.products[id]
		if !ok {
			return ErrNotFound
		}
		if p.Stock < quantity {
			return fmt.Errorf("%w: %s (%d left)", ErrOutOfStock, p.Name, p.Stock)
		}
	}

	for id, quantity := range quantities {
		p := r.products[id]
		p.Stock -= quantity
		r.products[id] = p
	}

	return nil
}
