package product

import (
	"context"
	"fmt"
	"time"
)

var catalog = []Product{
	{ID: "7c1e4d2a-3b5f-4e8a-9c6d-000000000001", Name: "Wireless Mouse", Price: 2499, Stock: 25},
	{ID: "7c1e4d2a-3b5f-4e8a-9c6d-000000000002", Name: "Mechanical Keyboard", Price: 8999, Stock: 3},
	{ID: "7c1e4d2a-3b5f-4e8a-9c6d-000000000003", Name: "27-inch Monitor", Price: 22900, Stock: 5},
	{ID: "7c1e4d2a-3b5f-4e8a-9c6d-000000000004", Name: "USB-C Hub", Price: 3999, Stock: 0},
}

// Challenge: Create a seed script to insert the following
// products when we start our API.
func Seed(ctx context.Context, repo Repository) error {
	for _, p := range catalog {
		p.CreatedAt = time.Now().UTC()
		if err := repo.Create(ctx, p); err != nil {
			return fmt.Errorf("seed %s: %w", p.Name, err)
		}
	}

	return nil
}
