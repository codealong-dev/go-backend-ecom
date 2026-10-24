package user

import (
	"context"
	"fmt"
)

var demoCustomers = []User{
	{ID: "9a8b7c6d-5e4f-4a3b-8c2d-000000000001", Name: "Ada"},
}

func Seed(ctx context.Context, repo Repository) error {
	for _, u := range demoCustomers {
		if err := repo.Create(ctx, u); err != nil {
			return fmt.Errorf("seed %s: %w", u.Name, err)
		}
	}
	return nil
}
