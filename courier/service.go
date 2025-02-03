package courier

import (
	"context"
	"santapan/domain"
)

// PostgresRepositoryQueries defines the methods for querying the category repository.
type PostgresRepositoryQueries interface {
	Fetch(ctx context.Context) ([]domain.Courier, error)
	FetchByID(ctx context.Context, id int) (*domain.Courier, error)
}

// Service defines the methods for executing commands on the category repository.
type Service struct {
	postgresRepoQuery PostgresRepositoryQueries
}

// NewService will create a new category service object.
func NewService(pq PostgresRepositoryQueries) *Service {
	return &Service{
		postgresRepoQuery: pq,
	}
}

// Fetch retrieves all couriers.
func (c *Service) Fetch(ctx context.Context) ([]domain.Courier, error) {
	res, err := c.postgresRepoQuery.Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// FetchByID retrieves a courier by its ID.
func (c *Service) FetchByID(ctx context.Context, id int) (*domain.Courier, error) {
	return c.postgresRepoQuery.FetchByID(ctx, id)
}
