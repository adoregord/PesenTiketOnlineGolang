package domain

// for making new ticket
type Ticket struct {
	ID       int     `json:"id,omitempty" validate:"omitempty"`
	Type     string  `json:"type" validate:"noblank"`
	Quantity int     `json:"quantity" validate:"required,gt=0,numeric"`
	Price    float64 `json:"price,omitempty" validate:"gt=0,numeric"`
}
