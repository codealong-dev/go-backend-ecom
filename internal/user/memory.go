package user

import (
	"context"
	"sort"
	"sync"
)

type MemoryRepository struct {
	mu    sync.RWMutex
	users map[string]User
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{users: make(map[string]User)}
}

func (r *MemoryRepository) Create(_ context.Context, u User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[u.ID] = u
	return nil
}

func (r *MemoryRepository) ByID(_ context.Context, id string) (User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	u, ok := r.users[id]
	if !ok {
		return User{}, ErrNotFound
	}
	return u, nil
}

func (r *MemoryRepository) List(_ context.Context) ([]User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]User, 0, len(r.users))
	for _, u := range r.users {
		out = append(out, u)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].ID < out[j].ID })
	return out, nil
}
