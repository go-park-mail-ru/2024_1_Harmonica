package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"harmonica/config"
	"harmonica/internal/handler"
	"harmonica/internal/handler/middleware"
	"harmonica/internal/repository"
	"harmonica/internal/service"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	v3 "github.com/swaggest/swgui/v3"
)

func runServer(addr string) {
	//logger := zap.Must(zap.NewProduction())
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
	h := handler.NewAPIHandler(s, logger)

	mux := http.NewServeMux()

	go handler.CleanupSessions()

	configureUserRoutes(logger, h, mux)
	configurePinRoutes(logger, h, mux)
	configureBoardRoutes(logger, h, mux)

	mux.Handle("GET /docs/swagger.json", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))
	mux.Handle("GET /swagger/", v3.NewHandler("My API", "/docs/swagger.json", "/swagger"))
	mux.HandleFunc("GET /img/{image_name}", h.GetImage)
	mux.HandleFunc("GET /api/v1/CSAT", h.GetRatings)
	mux.HandleFunc("POST /api/v1/CSAT", h.CreateRatings)

	loggedMux := middleware.Logging(logger, mux)

	server := http.Server{
		Addr:    addr,
		Handler: middleware.CSRF(middleware.CORS(loggedMux)),
	}
	server.ListenAndServeTLS("/etc/letsencrypt/live/harmoniums.ru/fullchain.pem", "/etc/letsencrypt/live/harmoniums.ru/privkey.pem")
}

func configureZapLogger() *zap.Logger {
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/harmonium.log",
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

func configureUserRoutes(logger *zap.Logger, h *handler.APIHandler, mux *http.ServeMux) {
	authRoutes := map[string]http.HandlerFunc{
		"POST /api/v1/users/{user_id}": h.UpdateUser,
	}
	notAuthRoutes := map[string]http.HandlerFunc{
		"POST /api/v1/login": h.Login,
		"POST /api/v1/users": h.Register,
	}
	checkAuthRoutes := map[string]http.HandlerFunc{
		"GET /api/v1/is_auth":          h.IsAuth, // check it
		"GET /api/v1/logout":           h.Logout,
		"GET /api/v1/users/{nickname}": h.GetUser,
	}
	for pattern, f := range authRoutes {
		mux.HandleFunc(pattern, middleware.AuthRequired(logger, f))
	}
	for pattern, f := range notAuthRoutes {
		mux.HandleFunc(pattern, middleware.NoAuthRequired(logger, f))
	}
	for pattern, f := range checkAuthRoutes {
		mux.HandleFunc(pattern, middleware.CheckAuth(logger, f))
	}
}

func configurePinRoutes(logger *zap.Logger, h *handler.APIHandler, mux *http.ServeMux) {
	authRoutes := map[string]http.HandlerFunc{
		"POST /api/v1/pins":                 h.CreatePin,
		"POST /api/v1/pins/{pin_id}":        h.UpdatePin,
		"DELETE /api/v1/pins/{pin_id}":      h.DeletePin,
		"POST /api/v1/pins/{pin_id}/like":   h.CreateLike,
		"DELETE /api/v1/pins/{pin_id}/like": h.DeleteLike,
	}
	checkAuthRoutes := map[string]http.HandlerFunc{
		"GET /api/v1/pins/{pin_id}": h.GetPin,
	}
	publicRoutes := map[string]http.HandlerFunc{
		"GET /api/v1/pins":                    h.Feed,
		"GET /api/v1/pins/created/{nickname}": h.UserPins,
		"GET /api/v1/likes/{pin_id}/users":    h.UsersLiked,
	}
	for pattern, f := range authRoutes {
		mux.HandleFunc(pattern, middleware.AuthRequired(logger, f))
	}
	for pattern, f := range checkAuthRoutes {
		mux.HandleFunc(pattern, middleware.CheckAuth(logger, f))
	}
	for pattern, f := range publicRoutes {
		mux.HandleFunc(pattern, f)
	}
}

func configureBoardRoutes(logger *zap.Logger, h *handler.APIHandler, mux *http.ServeMux) {
	authRoutes := map[string]http.HandlerFunc{
		"POST /api/v1/boards":                            h.CreateBoard,
		"POST /api/v1/boards/{board_id}":                 h.UpdateBoard,
		"DELETE /api/v1/boards/{board_id}":               h.DeleteBoard,
		"POST /api/v1/boards/{board_id}/pins/{pin_id}":   h.AddPinToBoard,
		"DELETE /api/v1/boards/{board_id}/pins/{pin_id}": h.DeletePinFromBoard,
	}
	checkAuthRoutes := map[string]http.HandlerFunc{
		"GET /api/v1/boards/{board_id}":         h.GetBoard,
		"GET /api/v1/boards/created/{nickname}": h.UserBoards,
	}
	for pattern, f := range authRoutes {
		mux.HandleFunc(pattern, middleware.AuthRequired(logger, f))
	}
	for pattern, f := range checkAuthRoutes {
		mux.HandleFunc(pattern, middleware.CheckAuth(logger, f))
	}
}

func init() {
	if err := godotenv.Load("conf.env"); err != nil {
		log.Print("No conf.env file found")
	}
}

// @title			Harmonium backend API
// @version		1.0
// @description	This is API-docs of backend server of Harmonica team.
// @host			https://harmoniums.ru
// @BasePath		api/v1
func main() {
	runServer(":8080")
}
