package main

import (
	image "harmonica/internal/microservices/image/proto"
	l "harmonica/internal/microservices/like/proto"
	like "harmonica/internal/microservices/like/server"
	"net"

	"harmonica/config"
	"harmonica/internal/microservices/like/server/repository"
	"harmonica/internal/microservices/like/server/service"
	"log"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func configureZapLogger() *zap.Logger {
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/like/harmonium.log",
		MaxSize:    1024, // MB
		MaxBackups: 10,
		MaxAge:     60, // days
		Compress:   true,
	})
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	core := zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), ws, zap.NewAtomicLevelAt(zap.InfoLevel))
	return zap.New(core, zap.AddCaller(), zap.AddStacktrace(zap.ErrorLevel))
}

func init() {
	if err := godotenv.Load("conf.env"); err != nil {
		log.Print("No conf.env file found")
	}
}

func main() {
	logger := configureZapLogger()
	defer logger.Sync()

	conf := config.New()
	connector, err := repository.NewConnector(conf)

	if err != nil {
		log.Print(err)
		return
	}
	defer connector.Disconnect()

	imageConn, err := grpc.Dial(config.GetEnv("IMAGE_MICROSERVICE_PORT", ":8003"), grpc.WithInsecure())
	if err != nil {
		log.Print(err)
		return
	}
	imageCli := image.NewImageClient(imageConn)

	r := repository.NewRepository(connector, logger, imageCli)
	s := service.NewService(r)

	lis, err := net.Listen("tcp", config.GetEnv("LIKE_MICROSERVICE_PORT", ":8004"))
	if err != nil {
		log.Print(err)
	}

	server := grpc.NewServer()
	l.RegisterLikeServer(server, like.NewLikeServer(s, logger))
	server.Serve(lis)
}
