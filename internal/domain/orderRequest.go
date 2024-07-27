package domain

type OrderRequest struct {
	UserID  int      `json:"userid" validate:"required,numeric"`
	EventID int      `json:"eventid" validate:"required,numeric"`
	Ticket  []Ticket `json:"ticket" validate:"required,min=1,dive"`
}
