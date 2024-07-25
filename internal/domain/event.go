package domain

type Event struct {
	ID          int      `json:"id"`
	Name        string   `json:"name" validate:"noblank,min=2"`
	Date        string   `json:"date" validate:"Datetime"`
	Description string   `json:"description" validate:"noblank"`
	Location    string   `json:"location" validate:"noblank"`
	Ticket      []Ticket `json:"ticket" validate:"dive"`
}
