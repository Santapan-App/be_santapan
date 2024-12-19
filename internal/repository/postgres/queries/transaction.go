package queries

import (
	"context"
	"database/sql"
	"santapan_transaction_service/domain"

	"github.com/sirupsen/logrus"
)

type PostgresTransactionQueryRepository struct {
	conn *sql.DB
}

func NewPostgresTransactionQueryRepository(conn *sql.DB) *PostgresTransactionQueryRepository {
	return &PostgresTransactionQueryRepository{conn: conn}
}

func (m *PostgresTransactionQueryRepository) fetch(ctx context.Context, query string, args ...interface{}) (result []domain.Transaction, err error) {
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

	result = make([]domain.Transaction, 0)
	for rows.Next() {
		transaction := domain.Transaction{}
		err = rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.CartID,
			&transaction.PaymentID,
			&transaction.CourierID,
			&transaction.AddressID,
			&transaction.Amount,
			&transaction.Status,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)

		if err != nil {
			logrus.Error("Failed to scan row:", err)
			return nil, err
		}
		result = append(result, transaction)
	}

	return result, nil
}

// Get By User ID
func (m *PostgresTransactionQueryRepository) GetByUserID(ctx context.Context, userID int64) (res domain.Transaction, err error) {
	query := `SELECT * FROM transaction WHERE user_id = $1`
	result, err := m.fetch(ctx, query, userID)
	if err != nil {
		return domain.Transaction{}, err
	}

	if len(result) > 0 {
		return result[0], nil
	}

	return domain.Transaction{}, sql.ErrNoRows
}

// Get By ID
func (m *PostgresTransactionQueryRepository) GetByID(ctx context.Context, id int64) (res domain.Transaction, err error) {
	query := `SELECT * FROM transaction WHERE id = $1`
	result, err := m.fetch(ctx, query, id)
	if err != nil {
		return domain.Transaction{}, err
	}

	if len(result) > 0 {
		return result[0], nil
	}

	return domain.Transaction{}, sql.ErrNoRows
}

// Validate
func (m *PostgresTransactionQueryRepository) Validate(ctx context.Context, transactionID int64, userID int64) (err error) {
	query := `SELECT * FROM transaction WHERE id = $1 AND user_id = $2`
	result, err := m.fetch(ctx, query, transactionID, userID)
	if err != nil {
		return err
	}

	if len(result) == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// GetOngoing pending and processing status
func (m *PostgresTransactionQueryRepository) GetOngoing(ctx context.Context, userID int64) (res []domain.Transaction, err error) {
	query := `SELECT * FROM transaction WHERE user_id = $1 AND status IN ('pending', 'processing')`
	result, err := m.fetch(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	return result, nil
}
