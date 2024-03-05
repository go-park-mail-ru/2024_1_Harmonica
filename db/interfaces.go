package db

import "harmonica/models"

//go:generate mockery --name Authorization
type Authorization interface {
	GetUserByEmail(email string) (User, error)
	GetUserById(id int64) (User, error)
	RegisterUser(user User) error
}

//go:generate mockery --name Pins
type Pins interface {
	GetPins(limit, offset int) (models.Pins, error)
}
