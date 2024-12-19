package menu

import (
	"context"
	"santapan/domain"
)

// PostgresRepositoryQueries defines the methods for querying the category repository.
type PostgresRepositoryQueries interface {
	GetByID(ctx context.Context, id int64) (domain.Menu, error)
	GetByCategoryID(ctx context.Context, categoryID int64) ([]domain.Menu, error)
	Fetch(ctx context.Context, cursor string, num int64) (res []domain.Menu, nextCursor string, err error)
}

// PostgresRepositoryCommand defines the methods for executing commands on the category repository.
type PostgresRepositoryCommand interface {
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

// Fetch retrieves all categories.
func (c *Service) Fetch(ctx context.Context, cursor string, num int64) (res []domain.Menu, nextCursor string, err error) {
	res, nextCursor, err = c.postgresRepoQuery.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}
	return
}

// GetByID retrieves a category by its ID.
func (c *Service) GetByID(ctx context.Context, id int64) (domain.Menu, error) {
	menu, err := c.postgresRepoQuery.GetByID(ctx, id)
	if err != nil {
		return domain.Menu{}, err
	}
	return menu, nil
}

// GetByCategoryID retrieves a category by its category ID.
func (c *Service) GetByCategoryID(ctx context.Context, categoryID int64) ([]domain.Menu, error) {
	menu, err := c.postgresRepoQuery.GetByCategoryID(ctx, categoryID)
	if err != nil {
		return nil, err
	}
	return menu, nil
}
