package item

import (
	"context"
	"santapan_transaction_service/domain"
)

// PostgresRepositoryQueries defines the methods for querying the category repository.
type PostgresRepositoryQueries interface {
	GetByCartID(ctx context.Context, cartID int64) ([]domain.CartItem, error)
	GetByID(ctx context.Context, id int64) (domain.CartItem, error)
}

// PostgresRepositoryCommand defines the methods for executing commands on the category repository.
type PostgresRepositoryCommand interface {
	Store(ctx context.Context, cart *domain.CartItem) error
	Update(ctx context.Context, cart *domain.CartItem) error // Explicit update method
	Delete(ctx context.Context, id int64) error
}

//go:generate mockery --name CategoryRepository
type Service struct {
	postgresRepoQuery   PostgresRepositoryQueries
	postgresRepoCommand PostgresRepositoryCommand
}

// NewService creates a new category service.
func NewService(postgresRepoQuery PostgresRepositoryQueries, postgresRepoCommand PostgresRepositoryCommand) *Service {
	return &Service{
		postgresRepoQuery:   postgresRepoQuery,
		postgresRepoCommand: postgresRepoCommand,
	}
}

// Fetch fetches the cart based on the user ID.
func (s *Service) GetByCartID(ctx context.Context, cartID int64) ([]domain.CartItem, error) {
	return s.postgresRepoQuery.GetByCartID(ctx, cartID)
}

func (s *Service) Store(ctx context.Context, cart *domain.CartItem) error {
	return s.postgresRepoCommand.Store(ctx, cart)
}

func (s *Service) Update(ctx context.Context, cart *domain.CartItem) error {
	return s.postgresRepoCommand.Update(ctx, cart)
}

// Delete
func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.postgresRepoCommand.Delete(ctx, id)
}

// GetByID
func (s *Service) GetByID(ctx context.Context, id int64) (domain.CartItem, error) {
	return s.postgresRepoQuery.GetByID(ctx, id)
}
