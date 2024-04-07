package main

import (
	"harmonica/config"
	"harmonica/internal/handler"
	"harmonica/internal/handler/middleware"
	"harmonica/internal/repository"
	"harmonica/internal/service"
	"log"
	"net/http"

	"go.uber.org/zap"

	"github.com/joho/godotenv"
	v3 "github.com/swaggest/swgui/v3"
)

func runServer(addr string) {
	logger := zap.Must(zap.NewProduction())

	conf := config.New()
	connector, err := repository.NewConnector(conf)
	if err != nil {
		log.Print(err)
		return
	}
	defer connector.Disconnect()
	r := repository.NewRepository(connector)
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

	server := http.Server{
		Addr:    addr,
		Handler: middleware.CORS(mux),
	}
	server.ListenAndServeTLS("/etc/letsencrypt/live/harmoniums.ru/fullchain.pem", "/etc/letsencrypt/live/harmoniums.ru/privkey.pem")
}

func configureUserRoutes(logger *zap.Logger, h *handler.APIHandler, mux *http.ServeMux) {
	authRoutes := map[string]http.HandlerFunc{
		"POST /api/v1/users/{user_id}": h.UpdateUser,
		"GET /api/v1/is_auth":          h.IsAuth,
	}
	notAuthRoutes := map[string]http.HandlerFunc{
		"POST /api/v1/login": h.Login,
		"POST /api/v1/users": h.Register,
	}
	publicRoutes := map[string]http.HandlerFunc{
		"GET /api/v1/logout": h.Logout,
	}
	for pattern, f := range authRoutes {
		mux.HandleFunc(pattern, middleware.Auth(logger, f))
	}
	for pattern, f := range notAuthRoutes {
		mux.HandleFunc(pattern, middleware.NotAuth(logger, f))
	}
	for pattern, f := range publicRoutes {
		mux.HandleFunc(pattern, f)
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
		mux.HandleFunc(pattern, middleware.Auth(logger, f))
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
		"GET /api/v1/boards/{board_id}":         h.GetBoard,   // с пагинацией
		"GET /api/v1/boards/created/{nickname}": h.UserBoards, // с пагинацией
	}
	for pattern, f := range authRoutes {
		mux.HandleFunc(pattern, middleware.Auth(logger, f))
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

//mux.HandleFunc("POST /api/v1/login", middleware.NotAuth(logger, handler.Login))
//mux.HandleFunc("GET /api/v1/logout", handler.Logout)
//mux.HandleFunc("POST /api/v1/users", middleware.NotAuth(logger, handler.Register))
//mux.HandleFunc("POST /api/v1/users/{user_id}", middleware.Auth(logger, handler.UpdateUser))
//mux.HandleFunc("GET /api/v1/is_auth", middleware.Auth(logger, handler.IsAuth))
//
//mux.HandleFunc("GET /api/v1/pins/created/{user_id}", handler.UserPins)
//mux.HandleFunc("GET /api/v1/pins", handler.Feed)
//mux.HandleFunc("POST /api/v1/pins", handler.CreatePin) // Обернуть в OnlyAuth
//mux.HandleFunc("GET /api/v1/pins/{pin_id}", handler.GetPin)
//mux.HandleFunc("POST /api/v1/pins/{pin_id}", handler.UpdatePin)   // Обернуть в OnlyAuth
//mux.HandleFunc("DELETE /api/v1/pins/{pin_id}", handler.DeletePin) // Обернуть в OnlyAuth
//
//mux.HandleFunc("POST /api/v1/pins/{pin_id}/like", handler.CreateLike)   // Обернуть в OnlyAuth
//mux.HandleFunc("DELETE /api/v1/pins/{pin_id}/like", handler.DeleteLike) // Обернуть в OnlyAuth
//mux.HandleFunc("GET /api/v1/likes/{pin_id}/users", handler.UsersLiked)
