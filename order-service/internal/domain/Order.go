package domain

import "time"

type Order struct {
	Id       string      `json:"id" gorm:"primarykey"`
	UserId   string      `json:"user_id"`
	Amount   float64     `json:"amount"`
	Status   OrderStatus `json:"status"`
	CreateAt time.Time   `json:"create_at" gorm:"autoCreateTime"`
	UpdateAt time.Time   `json:"update_at" gorm:"autoCreateTime"`
}
