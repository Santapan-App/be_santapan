package domain

import (
	"time"
)

// Cart represents a user's shopping cart.
type Cart struct {
	ID         int64     `json:"id"`          // Unique identifier for the cart
	UserID     int64     `json:"user_id"`     // Reference to the user owning the cart
	Status     string    `json:"status"`      // Status of the cart (e.g. "active", "completed")
	TotalPrice float64   `json:"total_price"` // Total price of items in the cart
	CreatedAt  time.Time `json:"created_at"`  // Timestamp for when the cart was created
	UpdatedAt  time.Time `json:"updated_at"`  // Timestamp for when the cart was last updated
}
