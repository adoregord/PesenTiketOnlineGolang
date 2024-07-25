package domain

type Ticket struct {
	ID    int     `json:"id"`
	Type  string  `json:"type" validate:"noblank"`
	Stock int     `json:"stock" validate:"gt=0,numeric"`
	Price float64 `json:"price" validate:"gt=0,numeric"`
}
