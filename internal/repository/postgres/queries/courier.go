package queries

import (
	"context"
	"database/sql"
	"santapan/domain"

	"github.com/sirupsen/logrus"
)

type CourierRepository struct {
	Conn *sql.DB
}

// NewCourierRepository creates an instance of CourierRepository.
func NewPostgresCourierQueryRepository(conn *sql.DB) *CourierRepository {
	return &CourierRepository{conn}
}

func (m *CourierRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Courier, err error) {
	rows, err := m.Conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer func() {
		errRow := rows.Close()
		if errRow != nil {
			logrus.Error(errRow)
		}
	}()

	result = make([]domain.Courier, 0)
	for rows.Next() {
		t := domain.Courier{}
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Logo,
			&t.Price,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

// Fetch retrieves all couriers without pagination.
func (m *CourierRepository) Fetch(ctx context.Context) (res []domain.Courier, err error) {
	query := `SELECT id, name, logo, price, created_at, updated_at FROM couriers`
	res, err = m.fetch(ctx, query)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// FetchByID retrieves a courier by its ID.
func (m *CourierRepository) FetchByID(ctx context.Context, id int) (*domain.Courier, error) {
	query := `SELECT id, name, logo, price, created_at, updated_at FROM couriers WHERE id=$1`
	list, err := m.fetch(ctx, query, id)
	if err != nil {
		return nil, err
	}

	if len(list) > 0 {
		return &list[0], nil
	}

	return nil, nil
}
