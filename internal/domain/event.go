package domain

type Event struct {
	ID          int      `json:"id,omitempty"`
	Name        string   `json:"name" validate:"required,noblank,min=2"`
	Date        string   `json:"date" validate:"required,Datetime"`
	Description string   `json:"description" validate:"required,noblank"`
	Location    string   `json:"location" validate:"required,noblank"`
	Ticket      []Ticket `json:"ticket,omitempty" validate:"dive"`
}
