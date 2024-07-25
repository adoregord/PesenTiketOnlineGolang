package domain

type Event struct {
	Id          int      `json:"id"`
	Name        string   `json:"name" validate:"noblank"`
	Date        string   `json:"date" validate:"datetime"`
	Description string   `json:"description" validate:"noblank"`
	Location    string   `json:"location" validate:"noblank"`
	Ticket      []Ticket `json:"ticket" validate:"dive"`
}
