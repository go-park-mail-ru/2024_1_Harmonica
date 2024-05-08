package main

import (
	au "harmonica/internal/microservices/auth/proto"
	auth "harmonica/internal/microservices/auth/server"
	"net"

	"harmonica/config"
	"harmonica/internal/microservices/auth/server/repository"
	"harmonica/internal/microservices/auth/server/service"
	"log"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func init() {
	if err := godotenv.Load("conf.env"); err != nil {
		log.Print("No conf.env file found")
	}
}

func main() {
	logger := config.ConfigureZapLogger("auth")
	defer logger.Sync()

	conf := config.New()
	connector, err := repository.NewConnector(conf)

	if err != nil {
		log.Print(err)
		return
	}
	defer connector.Disconnect()
	r := repository.NewRepository(connector, logger)
	s := service.NewService(r)

	go auth.CleanupSessions()

	lis, err := net.Listen("tcp", config.GetEnv("AUTH_MICROSERVICE_PORT", ":8002"))
	if err != nil {
		log.Print(err)
	}

	server := grpc.NewServer()
	au.RegisterAuthorizationServer(server, auth.NewAuthorizationServer(s))
	server.Serve(lis)
}
