package handler

import (
	"fmt"
	//"github.com/stretchr/testify/assert"
	"harmonica/db"
	"harmonica/db/mocks"
	"testing"
)

func TestLogin(t *testing.T) {
	mock := mocks.NewMethods(t)
	//или mock := new(mocks.Methods)
	conn := &db.Connector{Methods: mock}
	mock.On("GetUserByEmail", "test2@kkk.k").Return(db.User{}, nil).Once()

	user, err := conn.GetUserByEmail("test2@kkk.k")
	if err != nil {
		t.Errorf("Expected no error, got: %v", err)
	}

	fmt.Println(user)
	mock.AssertExpectations(t)
}
