package main

import (
	"go.uber.org/zap"
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
	logger := zap.Must(zap.NewProduction())

	conf := config.New()
	dbConn, err := repository.NewConnector(conf.DB)
	if err != nil {
		log.Print(err)
		return
	}
	defer dbConn.Disconnect()

	r := repository.NewRepository(dbConn)
	s := service.NewService(r)
	h := handler.NewAPIHandler(s, logger)

	mux := http.NewServeMux()

	go handler.CleanupSessions()

	configureUserRoutes(logger, h, mux)
	configurePinRoutes(logger, h, mux)

	mux.Handle("GET /img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img"))))
	mux.Handle("GET /docs/swagger.json", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))
	mux.Handle("GET /swagger/", v3.NewHandler("My API", "/docs/swagger.json", "/swagger"))

	server := http.Server{
		Addr:    addr,
		Handler: middleware.CORS(mux),
	}
	//server.ListenAndServeTLS("cert.pem", "key.pem")
	server.ListenAndServe()
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
	checkAuthRoutes := map[string]http.HandlerFunc{
		"GET /api/v1/users/{nickname}": h.GetUser,
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
	for pattern, f := range checkAuthRoutes {
		mux.HandleFunc(pattern, middleware.CheckAuth(logger, f))
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
	publicRoutes := map[string]http.HandlerFunc{
		"GET /api/v1/pins":                   h.Feed,
		"GET /api/v1/pins/{pin_id}":          h.GetPin,
		"GET /api/v1/pins/created/{user_id}": h.UserPins,
		"GET /api/v1/likes/{pin_id}/users":   h.UsersLiked,
	}
	for pattern, f := range authRoutes {
		mux.HandleFunc(pattern, middleware.Auth(logger, f))
	}
	for pattern, f := range publicRoutes {
		mux.HandleFunc(pattern, f)
	}
}

func init() {
	if err := godotenv.Load("conf.env"); err != nil {
		log.Print("No conf.env file found")
	}
}

//	@title			Harmonium backend API
//	@version		1.0
//	@description	This is API-docs of backend server of Harmonica team.

// @host		https://85.192.35.36:8080
// @BasePath	api/v1
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
