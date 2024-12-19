package cart

import (
	"context"
	"santapan_transaction_service/domain"
)

// PostgresRepositoryQueries defines the methods for querying the category repository.
type PostgresRepositoryQueries interface {
	GetByUserID(ctx context.Context, userID int64, status string) (domain.Cart, error)
	GetByID(ctx context.Context, cartID int64) (domain.Cart, error)
	Validate(ctx context.Context, userID int64, cartID int64) error
}

// PostgresRepositoryCommand defines the methods for executing commands on the category repository.
type PostgresRepositoryCommand interface {
	Store(ctx context.Context, cart *domain.Cart) error
	Update(ctx context.Context, cart *domain.Cart) error          // Explicit update method
	Delete(ctx context.Context, userID int64, cartID int64) error // New method for deleting a cart
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
func (s *Service) GetByUserID(ctx context.Context, userID int64, status string) (domain.Cart, error) {
	return s.postgresRepoQuery.GetByUserID(ctx, userID, status)
}

// Store stores the cart.
func (s *Service) Store(ctx context.Context, cart *domain.Cart) error {
	return s.postgresRepoCommand.Store(ctx, cart)
}

// Update updates the cart.
func (s *Service) Update(ctx context.Context, cart *domain.Cart) error {
	return s.postgresRepoCommand.Update(ctx, cart)
}

// Delete deletes the cart.
func (s *Service) Delete(ctx context.Context, userID int64, cartID int64) error {
	return s.postgresRepoCommand.Delete(ctx, userID, cartID)
}

// Validate validates the cart.
func (s *Service) Validate(ctx context.Context, userID int64, cartID int64) error {
	return s.postgresRepoQuery.Validate(ctx, userID, cartID)
}

// GetByID
func (s *Service) GetByID(ctx context.Context, cartID int64) (domain.Cart, error) {
	return s.postgresRepoQuery.GetByID(ctx, cartID)
}
