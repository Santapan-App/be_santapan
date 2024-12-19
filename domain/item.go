package domain

import "time"

// CartItem represents an item or bundle in a cart.
type CartItem struct {
	ID         int64     `json:"id"`          // Unique identifier for the cart item
	CartID     int64     `json:"cart_id"`     // Reference to the cart
	MenuID     *int64    `json:"menu_id"`     // Reference to the menu item (nullable)
	BundlingID *int64    `json:"bundling_id"` // Reference to the bundling (nullable)
	ImageUrl   string    `json:"image_url"`   // URL to the image of the item/bundle
	Name       string    `json:"name"`        // Name of the item/bundle
	Quantity   int       `json:"quantity"`    // Quantity of the item/bundle
	Price      float64   `json:"price"`       // Price per unit of the item/bundle
	CreatedAt  time.Time `json:"created_at"`  // Timestamp for when the item was added
	UpdatedAt  time.Time `json:"updated_at"`  // Timestamp for when the item was last updated
}
