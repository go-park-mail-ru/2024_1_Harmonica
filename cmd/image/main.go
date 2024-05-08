package main

import (
	im "harmonica/internal/microservices/image/proto"
	image "harmonica/internal/microservices/image/server"
	"net"

	"harmonica/config"
	"harmonica/internal/microservices/image/server/repository"
	"harmonica/internal/microservices/image/server/service"
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
	logger := config.ConfigureZapLogger("image")
	defer logger.Sync()

	conf := config.New()
	connector, err := repository.NewConnector(conf)

	if err != nil {
		log.Print(err)
		return
	}

	r := repository.NewRepository(connector, logger)
	s := service.NewService(r)

	lis, err := net.Listen("tcp", config.GetEnv("IMAGE_MICROSERVICE_PORT", ":8003"))
	if err != nil {
		log.Print(err)
	}

	server := grpc.NewServer()
	im.RegisterImageServer(server, image.NewImageServer(s, logger))
	server.Serve(lis)
}
