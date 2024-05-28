package main

import (
	"context"
	"fmt"
	"harmonica/config"
	"log"

	"harmonica/internal/entity"
	"harmonica/internal/repository"

	"flag"

	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

var EmailFlag = flag.String("e", "", "User email")
var PasswordFlag = flag.String("p", "", "User password")
var NicknameFlag = flag.String("n", "", "User nickname")

func main() {
	flag.Parse()

	conf := config.New()
	logger := config.ConfigureZapLogger("addUser")
	ctx := context.WithValue(context.Background(), "request_id", "test_request_id")
	connector, err := repository.NewConnector(conf, nil)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer connector.Disconnect()
	r := repository.NewRepository(connector, logger)

	userToCreate := entity.User{
		Email:    *EmailFlag,
		Nickname: *NicknameFlag,
		Password: *PasswordFlag,
	}

	fmt.Println("Createing user: {", userToCreate.Email, userToCreate.Nickname, userToCreate.Password, "}")

	hashPassword, err := bcrypt.GenerateFromPassword([]byte(userToCreate.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err.Error())
	}
	userToCreate.Password = string(hashPassword)

	result := r.RegisterUser(ctx, userToCreate)
	if result != nil {
		log.Fatal(result.Error())
	}
	log.Println("User created.")
}

func init() {
	if err := godotenv.Load("conf.env"); err != nil {
		log.Print("No conf.env file found")
	}
}
