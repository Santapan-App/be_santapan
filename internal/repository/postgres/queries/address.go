package queries

import (
	"context"
	"database/sql"
	"santapan/domain"

	"github.com/sirupsen/logrus"
)

type PostgresAddressQueryRepository struct {
	Conn *sql.DB
}

func NewPostgresAddressQueryRepository(Conn *sql.DB) *PostgresAddressQueryRepository {
	return &PostgresAddressQueryRepository{Conn}
}

func (m *PostgresAddressQueryRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Address, err error) {
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

	result = make([]domain.Address, 0)
	for rows.Next() {
		t := domain.Address{}
		err = rows.Scan(
			&t.ID,
			&t.UserID,
			&t.Label,
			&t.Address,
			&t.Name,
			&t.Notes,
			&t.Phone,
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

// Get By User ID retrieves an address by user ID.
func (m *PostgresAddressQueryRepository) GetByUserID(ctx context.Context, userID int64) (res []domain.Address, err error) {
	query := `SELECT id, user_id, label, address, name, notes, phone, created_at, updated_at
			  FROM address WHERE user_id=$1`

	res, err = m.fetch(ctx, query, userID)

	if err != nil {
		return nil, err
	}

	return
}

// GetByID retrieves an address by the given ID.
func (m *PostgresAddressQueryRepository) GetByID(ctx context.Context, id int64) (res domain.Address, err error) {
	query := `SELECT id, user_id, label, address, name, notes, phone, created_at, updated_at
			  FROM address WHERE id=$1`

	list, err := m.fetch(ctx, query, id)

	if err != nil {
		return domain.Address{}, err
	}

	if len(list) == 0 {
		return domain.Address{}, domain.ErrNotFound
	}

	res = list[0]
	return
}
