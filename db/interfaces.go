package db

import (
	"harmonica/models"
)

//go:generate mockery --name Methods
type IConnector interface {
	GetUserByEmail(email string) (User, error)
	GetUserById(id int64) (User, error)
	RegisterUser(user User) error
	GetPins(limit, offset int) (models.Pins, error)
}