package model

// Order represents an client order
type Order struct {
	ID           string      `json:"id,omitempty" bson:"_id"`
	CustomerID   string      `json:"customer_id" bson:"customer_id"`
	Status       string      `json:"status" bson:"status"`
	CreatedOn    int64       `json:"created_on,omitempty" bson:"created_on,omitempty"`
	RestaurantID string      `json:"restaurant_id" bson:"restaurant_id" validate:"nonzero"`
	OrderItems   []OrderItem `json:"order_items,omitempty" bson:"order_items,omitempty" validate:"nonzero, min=1"`
}

// OrderItem represents items in an order
type OrderItem struct {
	ProductCode string  `json:"product_code" bson:"product_code"`
	Name        string  `json:"name" bson:"name"`
	UnitPrice   float32 `json:"unit_price" bson:"unit_price"`
	Quantity    int32   `json:"quantity" bson:"quantity"`
}
