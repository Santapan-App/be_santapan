package queries

import (
	"context"
	"database/sql"
	"santapan_transaction_service/domain"

	"github.com/sirupsen/logrus"
)

type PostgresCartQueryRepository struct {
	conn *sql.DB
}

func NewPostgresCartQueryRepository(conn *sql.DB) *PostgresCartQueryRepository {
	return &PostgresCartQueryRepository{conn: conn}
}

func (m *PostgresCartQueryRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Cart, err error) {
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

	result = make([]domain.Cart, 0)
	for rows.Next() {
		cart := domain.Cart{}
		err = rows.Scan(
			&cart.ID,
			&cart.UserID,
			&cart.TotalPrice,
			&cart.Status,
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

// Get By User ID
func (m *PostgresCartQueryRepository) GetByUserID(ctx context.Context, userID int64, status string) (res domain.Cart, err error) {
	query := `SELECT * FROM cart WHERE user_id = $1 AND status = $2`
	result, err := m.fetch(ctx, query, userID, status)

	logrus.WithFields(logrus.Fields{
		"query":  query,
		"args":   []interface{}{userID, status},
		"result": result,
	}).Info("Result")

	if err != nil {
		return domain.Cart{}, err
	}

	if len(result) > 0 {
		res = result[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}

// Validate
func (m *PostgresCartQueryRepository) Validate(ctx context.Context, cartID int64, userID int64) (err error) {
	query := `SELECT * FROM cart WHERE id = $1 AND user_id = $2`
	result, err := m.fetch(ctx, query, cartID, userID)
	if err != nil {
		return err
	}

	if len(result) == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// Get By UserID
func (m *PostgresCartQueryRepository) GetByID(ctx context.Context, cartID int64) (res domain.Cart, err error) {
	query := `SELECT * FROM cart WHERE id = $1`
	list, err := m.fetch(ctx, query, cartID)

	if err != nil {
		return domain.Cart{}, err
	}

	if len(list) > 0 {
		res = list[0]
	} else {
		return res, domain.ErrNotFound
	}

	return
}
