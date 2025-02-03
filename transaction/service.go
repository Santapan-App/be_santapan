package transaction

import (
	"context"
	"santapan_transaction_service/domain"
)

// PostgresRepositoryQueries defines the methods for querying the category repository.
type PostgresRepositoryQueries interface {
	GetByUserID(ctx context.Context, userID int64) ([]domain.Transaction, error)
	GetOngoing(ctx context.Context, userID int64) ([]domain.Transaction, error)
	GetByID(ctx context.Context, transactionID int64) (domain.Transaction, error)
	Validate(ctx context.Context, userID int64, transactionID int64) error
}

// PostgresRepositoryCommand defines the methods for executing commands on the category repository.
type PostgresRepositoryCommand interface {
	Store(ctx context.Context, transaction *domain.Transaction) error
	Update(ctx context.Context, transaction *domain.Transaction) error // Explicit update method
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

// Fetch fetches the transaction based on the user ID.
func (s *Service) GetByUserID(ctx context.Context, userID int64) ([]domain.Transaction, error) {
	return s.postgresRepoQuery.GetByUserID(ctx, userID)
}

// GetByID fetches the transaction based on the transaction ID.
func (s *Service) GetByID(ctx context.Context, transactionID int64) (domain.Transaction, error) {
	return s.postgresRepoQuery.GetByID(ctx, transactionID)
}

// Store stores the transaction.
func (s *Service) Store(ctx context.Context, transaction *domain.Transaction) error {
	return s.postgresRepoCommand.Store(ctx, transaction)
}

// Update updates the transaction.
func (s *Service) Update(ctx context.Context, transaction *domain.Transaction) error {
	return s.postgresRepoCommand.Update(ctx, transaction)
}

// Validate validates the transaction.
func (s *Service) Validate(ctx context.Context, userID int64, transactionID int64) error {
	return s.postgresRepoQuery.Validate(ctx, userID, transactionID)
}

// GetOngoing fetches the ongoing transactions based on the user ID.
func (s *Service) GetOngoing(ctx context.Context, userID int64) ([]domain.Transaction, error) {
	return s.postgresRepoQuery.GetOngoing(ctx, userID)
}
