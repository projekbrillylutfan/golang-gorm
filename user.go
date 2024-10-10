package belajar_golang_gorm

import "time"

type User struct {
	Id   string
	Name string
	Password string
	CreatedAt time.Time
	UpdatedAt time.Time
}