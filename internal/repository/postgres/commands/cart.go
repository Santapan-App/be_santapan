package commands

import (
	"context"
	"database/sql"
	"fmt"
	"santapan_transaction_service/domain"
)

type PostgresCartCommandRepository struct {
	Conn *sql.DB
}

func NewPostgresCartCommandRepository(conn *sql.DB) *PostgresCartCommandRepository {
	return &PostgresCartCommandRepository{Conn: conn}
}

// Store: Tambahkan status
func (r *PostgresCartCommandRepository) Store(ctx context.Context, cart *domain.Cart) (err error) {
	query := `INSERT INTO cart (user_id, total_price, status, created_at, updated_at) 
              VALUES ($1, $2, $3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP) RETURNING id`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, cart.UserID, cart.TotalPrice, cart.Status).Scan(&cart.ID)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

// Update: Tambahkan status
func (r *PostgresCartCommandRepository) Update(ctx context.Context, cart *domain.Cart) (err error) {
	query := `UPDATE cart SET total_price=$1, status=$2, updated_at=CURRENT_TIMESTAMP WHERE id=$3`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, cart.TotalPrice, cart.Status, cart.ID)
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

// Delete Cart: Tambahkan logika untuk status jika diperlukan
func (r *PostgresCartCommandRepository) Delete(ctx context.Context, userID int64, cartID int64) (err error) {
	// Pilihan 1: Hapus langsung
	query := `DELETE FROM cart WHERE user_id=$1 AND id=$2`

	// Pilihan 2: Update status sebelum menghapus (gunakan salah satu jika diperlukan)
	// query := `UPDATE cart SET status='deleted', updated_at=CURRENT_TIMESTAMP WHERE user_id=$1 AND id=$2`

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, userID, cartID)
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
