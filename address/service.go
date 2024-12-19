package address

import (
	"context"
	"santapan/domain"
)

// PostgresRepositoryQueries defines the methods for querying the category repository.
type PostgresRepositoryQueries interface {
	GetByUserID(ctx context.Context, userID int64) ([]domain.Address, error)
	GetByID(ctx context.Context, id int64) (domain.Address, error)
}

// PostgresRepositoryCommand defines the methods for executing commands on the category repository.
type PostgresRepositoryCommand interface {
	Create(ctx context.Context, a domain.Address) (domain.Address, error)
	Update(ctx context.Context, a domain.Address) (domain.Address, error)
	Delete(ctx context.Context, id int64) error
}

//go:generate mockery --name CategoryRepository
type Service struct {
	postgresRepoQuery   PostgresRepositoryQueries
	postgresRepoCommand PostgresRepositoryCommand
}

// NewService will create a new category service object.
func NewService(pq PostgresRepositoryQueries, pc PostgresRepositoryCommand) *Service {
	return &Service{
		postgresRepoQuery:   pq,
		postgresRepoCommand: pc,
	}
}

// GetByUserID retrieves an address by user ID.
func (c *Service) GetByUserID(ctx context.Context, userID int64) ([]domain.Address, error) {
	res, err := c.postgresRepoQuery.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// GetByID retrieves an address by the given ID.
func (c *Service) GetByID(ctx context.Context, id int64) (domain.Address, error) {
	res, err := c.postgresRepoQuery.GetByID(ctx, id)
	if err != nil {
		return res, err
	}
	return res, nil
}

// Create creates a new address.
func (c *Service) Create(ctx context.Context, a domain.Address) (domain.Address, error) {
	res, err := c.postgresRepoCommand.Create(ctx, a)
	if err != nil {
		return res, err
	}
	return res, nil
}

// Update updates an address.
func (c *Service) Update(ctx context.Context, a domain.Address) (domain.Address, error) {
	res, err := c.postgresRepoCommand.Update(ctx, a)
	if err != nil {
		return res, err
	}
	return res, nil
}

// Delete deletes an address.
func (c *Service) Delete(ctx context.Context, id int64) error {
	err := c.postgresRepoCommand.Delete(ctx, id)
	if err != nil {
		return err
	}
	return nil
}
