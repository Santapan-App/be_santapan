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

		// Log the scanned transaction for debugging purposes
		logrus.WithFields(logrus.Fields{
			"transactionID": transaction.ID,
			"userID":        transaction.UserID,
			"amount":        transaction.Amount,
			"status":        transaction.Status,
		}).Info("Scanned transaction")

		result = append(result, transaction)
	}

	// Check if there were any errors during the row iteration
	if err = rows.Err(); err != nil {
		logrus.Error("Error during row iteration:", err)
		return nil, err
	}

	// Log the total number of results fetched
	logrus.WithFields(logrus.Fields{
		"resultCount": len(result),
	}).Info("Fetched transactions")

	return result, nil
}

// Get By User ID
func (m *PostgresTransactionQueryRepository) GetByUserID(ctx context.Context, userID int64) (res []domain.Transaction, err error) {
	query := `SELECT *  FROM transaction  WHERE user_id = $1 AND (status = 'success' OR status = 'failed')`
	result, err := m.fetch(ctx, query, userID)

	print("result", result)
	if err != nil {
		return nil, err // Return nil for empty slice if there was an error
	}

	return result, nil // Return the full result slice
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
