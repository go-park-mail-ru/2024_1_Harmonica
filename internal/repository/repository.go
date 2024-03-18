package repository

import "harmonica/internal/entity"

////go:generate mockery --name IConnector
//type IConnector interface {
//	GetUserByEmail(email string) (entity.User, error)
//	GetUserByNickname(nickname string) (entity.User, error)
//	GetUserById(id int64) (entity.User, error)
//	RegisterUser(user entity.User) error
//	GetPins(limit, offset int) (entity.Pins, error)
//}

type User interface {
	GetUserByEmail(email string) (entity.User, error)
	GetUserByNickname(nickname string) (entity.User, error)
	GetUserById(id int64) (entity.User, error)
	RegisterUser(user entity.User) error
}

type Pin interface {
	GetPins(limit, offset int) (entity.Pins, error)
}

type Repository struct {
	User
	Pin
}

func NewRepository(connector *Connector) *Repository {
	return &Repository{
		User: NewUserDB(connector.db),
		Pin:  NewPinDB(connector.db),
	}
}
