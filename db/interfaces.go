package db

import "harmonica/models"

// //go:generate mockery --name Methods

//go:generate mockgen -source=interfaces.go
type Methods interface {
	GetUserByEmail(email string) (User, error)
	GetUserById(id int64) (User, error)
	RegisterUser(user User) error
	GetPins(limit, offset int) (models.Pins, error)
}

//type MethodsStruct struct {
//	Methods
//}
//
//func NewMethodsStruct(i Methods) *MethodsStruct {
//	return &MethodsStruct{Methods: i}
//}
