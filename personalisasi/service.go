package personalisasi

import (
	"context"
	"santapan/domain"
)

// PostgresRepositoryCommand defines the methods for executing commands on the personalisasi repository.
type PostgresRepositoryCommand interface {
	InsertOrUpdate(ctx context.Context, p *domain.Personalisasi) (domain.Personalisasi, error)
}

// PersonalisasiRepository defines the methods for executing queries on the personalisasi repository.
type PostgresRepositoryQuery interface {
	GetByUserID(ctx context.Context, userID int64) (domain.Personalisasi, error)
}

//go:generate mockery --name PersonalisasiRepository
type Service struct {
	postgresRepoCommand PostgresRepositoryCommand
	postgresRepoQuery   PostgresRepositoryQuery
}

// NewService will create a new personalisasi service object with both query and command repositories.
func NewService(pc PostgresRepositoryCommand, pq PostgresRepositoryQuery) *Service {
	return &Service{
		postgresRepoCommand: pc,
		postgresRepoQuery:   pq,
	}
}

// InsertOrUpdate inserts or updates the personalisasi for a user.
func (s *Service) InsertOrUpdate(ctx context.Context, p domain.Personalisasi) (domain.Personalisasi, error) {
	// Insert or update the personalisasi record and return the result
	res, err := s.postgresRepoCommand.InsertOrUpdate(ctx, &p)
	if err != nil {
		return res, err
	}
	return res, nil
}

// GetByUserID retrieves the personalisasi data for a user.
func (s *Service) GetByUserID(ctx context.Context, userID int64) (domain.Personalisasi, error) {
	res, err := s.postgresRepoQuery.GetByUserID(ctx, userID)
	if err != nil {
		return domain.Personalisasi{}, err
	}
	return res, nil
}
