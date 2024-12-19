package commands

import (
	"context"
	"database/sql"
	"santapan/domain"
)

type PostgresAddressCommandRepository struct {
	Conn *sql.DB
}

func NewPostgresAddressCommandRepository(Conn *sql.DB) *PostgresAddressCommandRepository {
	return &PostgresAddressCommandRepository{Conn}
}

// Create method for inserting a new address, including label and notes
func (m *PostgresAddressCommandRepository) Create(ctx context.Context, a domain.Address) (domain.Address, error) {
	query := `INSERT INTO address (name, address, phone, user_id, label, notes) 
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := m.Conn.QueryRowContext(ctx, query, a.Name, a.Address, a.Phone, a.UserID, a.Label, a.Notes).Scan(&a.ID)
	if err != nil {
		return a, err
	}
	return a, nil
}

// Update method for modifying an existing address, including label and notes
func (m *PostgresAddressCommandRepository) Update(ctx context.Context, a domain.Address) (domain.Address, error) {
	query := `UPDATE address SET name=$1, address=$2, phone=$3, label=$4, notes=$5, updated_at=NOW() 
              WHERE id=$6`
	_, err := m.Conn.ExecContext(ctx, query, a.Name, a.Address, a.Phone, a.Label, a.Notes, a.ID)
	if err != nil {
		return a, err
	}
	return a, nil
}

// Delete method remains unchanged
func (m *PostgresAddressCommandRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM address WHERE id=$1`
	_, err := m.Conn.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
