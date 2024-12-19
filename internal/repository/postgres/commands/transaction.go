package commands

import (
	"context"
	"database/sql"
	"fmt"

	"santapan_transaction_service/domain"
)

// PostgresTransactionCommandRepository struct
type PostgresTransactionCommandRepository struct {
	Conn *sql.DB
}

// NewPostgresTransactionCommandRepository creates a new instance of PostgresTransactionCommandRepository
func NewPostgresTransactionCommandRepository(conn *sql.DB) *PostgresTransactionCommandRepository {
	return &PostgresTransactionCommandRepository{Conn: conn}
}

// Store stores the transaction
func (r *PostgresTransactionCommandRepository) Store(ctx context.Context, transaction *domain.Transaction) (err error) {
	query := `INSERT INTO transaction (user_id, cart_id, payment_id, courier_id, address_id, amount, status, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, transaction.UserID, transaction.CartID, transaction.PaymentID, transaction.CourierID, transaction.AddressID, transaction.Amount, transaction.Status, transaction.CreatedAt, transaction.UpdatedAt).Scan(&transaction.ID)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

// Update the transaction
func (r *PostgresTransactionCommandRepository) Update(ctx context.Context, transaction *domain.Transaction) (err error) {
	query := `UPDATE transaction SET status=$1, updated_at=CURRENT_TIMESTAMP WHERE id=$2`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, transaction.Status, transaction.ID)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	affect, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if affect != 1 {
		return fmt.Errorf("unexpected number of affected rows: %d", affect)
	}

	return nil
}
