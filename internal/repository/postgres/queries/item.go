package queries

import (
	"context"
	"database/sql"
	"santapan_transaction_service/domain"

	"github.com/sirupsen/logrus"
)

type PostgresItemQueryRepository struct {
	conn *sql.DB
}

func NewPostgresItemQueryRepository(conn *sql.DB) *PostgresItemQueryRepository {
	return &PostgresItemQueryRepository{conn: conn}
}

func (m *PostgresItemQueryRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.CartItem, err error) {
	rows, err := m.conn.QueryContext(ctx, query, args...)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"query": query,
			"args":  args,
		}).Error("Failed to execute query")
		return nil, err
	}

	defer func() {
		if errRow := rows.Close(); errRow != nil {
			logrus.Error("Failed to close rows:", errRow)
		}
	}()

	result = make([]domain.CartItem, 0)
	for rows.Next() {
		cart := domain.CartItem{}
		err = rows.Scan(
			&cart.ID,
			&cart.CartID,
			&cart.MenuID,
			&cart.BundlingID,
			&cart.ImageUrl,
			&cart.Name,
			&cart.Quantity,
			&cart.Price,
			&cart.CreatedAt,
			&cart.UpdatedAt,
		)

		if err != nil {
			logrus.Error("Failed to scan row:", err)
			return nil, err
		}
		result = append(result, cart)
	}

	return result, nil
}

// Get By Cart ID
func (m *PostgresItemQueryRepository) GetByCartID(ctx context.Context, cartID int64) (res []domain.CartItem, err error) {
	query := `SELECT * FROM cart_item WHERE cart_id = $1`
	result, err := m.fetch(ctx, query, cartID)

	if err != nil {
		return nil, err
	}

	return result, nil
}

// Get By ID
func (m *PostgresItemQueryRepository) GetByID(ctx context.Context, id int64) (res domain.CartItem, err error) {
	query := `SELECT * FROM cart_item WHERE id = $1`
	result, err := m.fetch(ctx, query, id)

	if err != nil {
		return domain.CartItem{}, err
	}

	if len(result) > 0 {
		res = result[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}
