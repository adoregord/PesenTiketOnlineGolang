package domain

type User struct {
	ID      int     `json:"id,omitempty"`
	Name    string  `json:"name" validate:"noblank"`
	Balance float64 `json:"balance" validate:"gt=0,numeric"`
}
