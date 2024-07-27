package domain

type User struct {
	ID      int     `json:"id,omitempty"`
	Name    string  `json:"name" validate:"noblank,min=2"`
	Balance float64 `json:"balance,omitempty" validate:"gt=0,numeric"`
}
