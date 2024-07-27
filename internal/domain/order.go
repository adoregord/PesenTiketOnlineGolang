package domain

type Order struct {
	ID            int      `json:"id,omitempty"`
	Date          string   `json:"date" validate:"Datetime"`
	Status        string   `json:"status"`
	PaymentMethod string   `json:"payment_method,omitempty" validate:"noblank"`
	User          User     `json:"user,omitempty" validate:"dive,min=2"`
	Event         Event    `json:"event,omitempty" validate:"dive"`
	EventTicket   []Ticket `json:"event_ticket,omitempty" validate:"dive"`
	TotalPrice    float64  `json:"total_price,omitempty" validate:"noblank"`
}
