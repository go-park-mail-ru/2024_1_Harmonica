package main

import (
	au "harmonica/internal/microservices/auth/proto"
	auth "harmonica/internal/microservices/auth/server"
	"net"

	"harmonica/config"
	"harmonica/internal/repository"
	"harmonica/internal/service"
	"log"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func configureZapLogger() *zap.Logger {
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/auth/harmonium.log",
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
	r := repository.NewRepository(connector, logger)
	s := service.NewService(r)

	go auth.CleanupSessions()

	lis, err := net.Listen("tcp", config.GetEnv("AUTH_MICRO_PORT", ":8001"))
	if err != nil {
		log.Print(err)
	}

	server := grpc.NewServer()
	au.RegisterAuthorizationServer(server, auth.NewAuthorizationServer(s, r))
	server.Serve(lis)
}
