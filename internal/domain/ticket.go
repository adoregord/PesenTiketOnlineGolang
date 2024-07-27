package domain

type Ticket struct {
	ID       int     `json:"id,omitempty"`
	Type     string  `json:"type" validate:"required,noblank"`
	Quantity int     `json:"quantity" validate:"required,gt=0,numeric"`
	Price    float64 `json:"price,omitempty" validate:"omitempty,gt=0,numeric"`
}
