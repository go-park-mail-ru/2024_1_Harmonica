// Code generated by mockery v2.42.0. DO NOT EDIT.

package mocks

import (
	db "harmonica/db"
	models "harmonica/models"

	mock "github.com/stretchr/testify/mock"
)

// Methods is an autogenerated mock type for the Methods type
type Methods struct {
	mock.Mock
}

// GetPins provides a mock function with given fields: limit, offset
func (_m *Methods) GetPins(limit int, offset int) (models.Pins, error) {
	ret := _m.Called(limit, offset)

	if len(ret) == 0 {
		panic("no return value specified for GetPins")
	}

	var r0 models.Pins
	var r1 error
	if rf, ok := ret.Get(0).(func(int, int) (models.Pins, error)); ok {
		return rf(limit, offset)
	}
	if rf, ok := ret.Get(0).(func(int, int) models.Pins); ok {
		r0 = rf(limit, offset)
	} else {
		r0 = ret.Get(0).(models.Pins)
	}

	if rf, ok := ret.Get(1).(func(int, int) error); ok {
		r1 = rf(limit, offset)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserByEmail provides a mock function with given fields: email
func (_m *Methods) GetUserByEmail(email string) (db.User, error) {
	ret := _m.Called(email)

	if len(ret) == 0 {
		panic("no return value specified for GetUserByEmail")
	}

	var r0 db.User
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (db.User, error)); ok {
		return rf(email)
	}
	if rf, ok := ret.Get(0).(func(string) db.User); ok {
		r0 = rf(email)
	} else {
		r0 = ret.Get(0).(db.User)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUserById provides a mock function with given fields: id
func (_m *Methods) GetUserById(id int64) (db.User, error) {
	ret := _m.Called(id)

	if len(ret) == 0 {
		panic("no return value specified for GetUserById")
	}

	var r0 db.User
	var r1 error
	if rf, ok := ret.Get(0).(func(int64) (db.User, error)); ok {
		return rf(id)
	}
	if rf, ok := ret.Get(0).(func(int64) db.User); ok {
		r0 = rf(id)
	} else {
		r0 = ret.Get(0).(db.User)
	}

	if rf, ok := ret.Get(1).(func(int64) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RegisterUser provides a mock function with given fields: user
func (_m *Methods) RegisterUser(user db.User) error {
	ret := _m.Called(user)

	if len(ret) == 0 {
		panic("no return value specified for RegisterUser")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(db.User) error); ok {
		r0 = rf(user)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewMethods creates a new instance of Methods. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMethods(t interface {
	mock.TestingT
	Cleanup(func())
}) *Methods {
	mock := &Methods{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}