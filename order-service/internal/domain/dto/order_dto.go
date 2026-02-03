package dto

type OrderDTO struct {
	UserId string  `json:"user_id"`
	Amount float64 `json:"amount"`
}
