package model

import "time"

type User struct {
	Id       string    `gorm:"primarykey"`
	Username string    `json:"username" binding:"required"`
	Password string    `json:"password" binding:"required"`
	Email    string    `json:"email" binding:"required,email"`
	CreateAt time.Time `json:"create_at" gorm:"autoCreateTime"`
	UpdateAt time.Time `json:"update_at" gorm:"autoUpdateTime"`
}
