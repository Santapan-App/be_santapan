package commands

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"santapan_transaction_service/domain"
)

type PostgresItemCommandRepository struct {
	Conn *sql.DB
}

func NewPostgresItemCommandRepository(conn *sql.DB) *PostgresItemCommandRepository {
	return &PostgresItemCommandRepository{Conn: conn}
}

func (r *PostgresItemCommandRepository) Store(ctx context.Context, item *domain.CartItem) error {
	if r.Conn == nil {
		return fmt.Errorf("database connection is nil")
	}

	if item == nil {
		return fmt.Errorf("item cannot be nil")
	}

	if item.MenuID == nil && item.BundlingID == nil {
		return fmt.Errorf("either MenuID or BundlingID must be provided")
	}

	query := `
        INSERT INTO cart_item (cart_id, menu_id, bundling_id, image_url, name, quantity, price, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
        RETURNING id
    `

	var menuID, bundlingID sql.NullInt64
	if item.MenuID != nil {
		menuID = sql.NullInt64{Int64: *item.MenuID, Valid: true}
	}
	if item.BundlingID != nil {
		bundlingID = sql.NullInt64{Int64: *item.BundlingID, Valid: true}
	}

	stmt, err := r.Conn.PrepareContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}
	defer stmt.Close()

	err = stmt.QueryRowContext(ctx, item.CartID, menuID, bundlingID, item.ImageUrl, item.Name, item.Quantity, item.Price).Scan(&item.ID)
	if err != nil {
		log.Printf("Error executing statement: %v", err)
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	log.Printf("Inserted cart item with ID: %d", item.ID)
	return nil
}

// Update modifies an existing item in the cart
func (r *PostgresItemCommandRepository) Update(ctx context.Context, item *domain.CartItem) error {
	query := `
		UPDATE cart_item
		SET quantity = $1, price = $2, updated_at = CURRENT_TIMESTAMP
		WHERE id = $3
	`

	_, err := r.Conn.ExecContext(ctx, query, item.Quantity, item.Price, item.ID)
	if err != nil {
		return fmt.Errorf("failed to update item: %w", err)
	}

	return nil
}

// Delete removes an item from the cart based on the item's ID
func (r *PostgresItemCommandRepository) Delete(ctx context.Context, itemID int64) error {
	query := `
		DELETE FROM cart_item
		WHERE id = $1
	`

	_, err := r.Conn.ExecContext(ctx, query, itemID)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	return nil
}
