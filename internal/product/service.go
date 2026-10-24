package product

import "context"

type Repository interface {
	Create(context.Context, Product) error
	ByID(context.Context, string) (Product, error)
	List(context.Context) ([]Product, error)
	TakeStock(context.Context, map[string]int) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) List(ctx context.Context) ([]Product, error) {
	return s.repo.List(ctx)
}
