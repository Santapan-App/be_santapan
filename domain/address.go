package domain

import "time"

// Address represents the structure of an address record in the database
type Address struct {
	ID        int64     `json:"id"`         // Unique identifier for the address
	UserID    int64     `json:"user_id"`    // The ID of the user who owns the address
	Label     string    `json:"label"`      // Label for the address (e.g., 'Home', 'Work')
	Address   string    `json:"address"`    // The actual address
	Name      string    `json:"name"`       // Name associated with the address
	Notes     string    `json:"notes"`      // Additional notes about the address
	Phone     string    `json:"phone"`      // Phone number associated with the address
	CreatedAt time.Time `json:"created_at"` // Timestamp when the address was created
	UpdatedAt time.Time `json:"updated_at"` // Timestamp when the address was last updated
}
