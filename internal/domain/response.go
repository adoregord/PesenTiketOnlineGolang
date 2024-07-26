package domain

type Response struct {
	Message string `json:"message"`
	Status  any    `json:"status,omitempty"`
	Data    any    `json:"data,omitempty"`
}

