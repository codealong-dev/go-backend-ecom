package user

import (
	"context"

	"fmt"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, u User) error
	ByID(ctx context.Context, id string) (User, error)
	List(ctx context.Context) ([]User, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Get(ctx context.Context, id string) (User, error) {
	return s.repo.ByID(ctx, id)
}

func (s *Service) List(ctx context.Context) ([]User, error) {
	return s.repo.List(ctx)
}

func (s *Service) Create(ctx context.Context, payload CreateUserRequest) (User, error) {
	if payload.Name == "" {
		return User{}, ErrNameRequired
	}

	newID := uuid.New().String()
	newUser := User{
		ID:   newID,
		Name: payload.Name,
	}

	if err := s.repo.Create(ctx, newUser); err != nil {
		return User{}, fmt.Errorf("create user: %w", err)
	}

	return newUser, nil
}
