package main

import (
	"go.uber.org/zap"
	"harmonica/config"
	h "harmonica/internal/handler"
	"harmonica/internal/handler/middleware"
	r "harmonica/internal/repository"
	s "harmonica/internal/service"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	v3 "github.com/swaggest/swgui/v3"
)

func runServer(addr string) {
	logger := zap.Must(zap.NewProduction())

	conf := config.New()
	dbConn, err := r.NewConnector(conf.DB)
	if err != nil {
		log.Print(err)
		return
	}
	defer dbConn.Disconnect()

	repo := r.NewRepository(dbConn, logger)
	service := s.NewService(repo, logger)
	handler := h.NewAPIHandler(service, logger)

	mux := http.NewServeMux()

	go h.CleanupSessions()

	mux.HandleFunc("POST /api/v1/login", middleware.NotAuth(logger, handler.Login))
	mux.HandleFunc("GET /api/v1/logout", handler.Logout)
	mux.HandleFunc("POST /api/v1/users", middleware.NotAuth(logger, handler.Register))
	mux.HandleFunc("POST /api/v1/users/{user_id}", middleware.Auth(logger, handler.UpdateUser))
	mux.HandleFunc("GET /api/v1/is_auth", middleware.Auth(logger, handler.IsAuth))

	mux.HandleFunc("GET /api/v1/pins/created/{user_id}", handler.UserPins)
	mux.HandleFunc("GET /api/v1/pins", handler.Feed)
	mux.HandleFunc("POST /api/v1/pins", handler.CreatePin) // Обернуть в OnlyAuth
	mux.HandleFunc("GET /api/v1/pins/{pin_id}", handler.GetPin)
	mux.HandleFunc("POST /api/v1/pins/{pin_id}", handler.UpdatePin)   // Обернуть в OnlyAuth
	mux.HandleFunc("DELETE /api/v1/pins/{pin_id}", handler.DeletePin) // Обернуть в OnlyAuth

	mux.HandleFunc("POST /api/v1/pins/{pin_id}/like", handler.CreateLike)   // Обернуть в OnlyAuth
	mux.HandleFunc("DELETE /api/v1/pins/{pin_id}/like", handler.DeleteLike) // Обернуть в OnlyAuth
	mux.HandleFunc("GET /api/v1/likes/{pin_id}/users", handler.UsersLiked)

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
