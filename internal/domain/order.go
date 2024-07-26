package domain

type Order struct {
	ID            int      `json:"id"`
	Date          string   `json:"date" validate:"Datetime"`
	Status        string   `json:"status" validate:"noblank"`
	PaymentMethod string   `json:"payment_method" validate:"noblank"`
	User          User     `json:"user" validate:"dive,min=2"`
	Event         Event    `json:"event" validate:"dive"`
	EventTicket   []Ticket `json:"event_ticket" validate:"dive"`
	TotalPrice    float64  `json:"total_price" validate:"noblank"`
}
