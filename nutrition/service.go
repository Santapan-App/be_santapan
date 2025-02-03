package nutrition

import (
	"context"
	"santapan/domain"
)

type PostgresRepositoryQuery interface {
	GetByClassification(ctx context.Context, class string) (domain.Nutrition, error)
}

type Service struct {
	postgresRepoQuery PostgresRepositoryQuery
}

func NewService(pq PostgresRepositoryQuery) *Service {
	return &Service{
		postgresRepoQuery: pq,
	}
}

func (s *Service) GetByClassification(ctx context.Context, class string) (domain.Nutrition, error) {
	res, err := s.postgresRepoQuery.GetByClassification(ctx, class)
	if err != nil {
		return domain.Nutrition{}, err
	}
	return res, nil
}
