package main

import (
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
	conf := config.New()
	dbConn, err := r.NewConnector(conf.DB)
	if err != nil {
		log.Print(err)
		return
	}
	defer dbConn.Disconnect()

	repo := r.NewRepository(dbConn)
	service := s.NewService(repo)
	handler := h.NewAPIHandler(service)
	//handler := handler2.NewAPIHandler(dbConn) // было
	mux := http.NewServeMux()

	go h.CleanupSessions()

	mux.HandleFunc("POST /api/v1/login", handler.Login)
	mux.HandleFunc("GET /api/v1/logout", handler.Logout)
	mux.HandleFunc("POST /api/v1/users", handler.Register)
	mux.HandleFunc("POST /api/v1/users/{user_id}", handler.UpdateUser)
	mux.HandleFunc("GET /api/v1/is_auth", handler.IsAuth)

	mux.HandleFunc("GET /api/v1/pins/created/{user_id}", handler.UserPins)
	mux.HandleFunc("GET /api/v1/pins", handler.Feed)
	mux.HandleFunc("POST /api/v1/pins", handler.CreatePin) // Обернуть в OnlyAuth
	mux.HandleFunc("GET /api/v1/pins/{pin_id}", handler.GetPin)
	mux.HandleFunc("POST /api/v1/pins/{pin_id}", handler.UpdatePin)   // Обернуть в OnlyAuth
	mux.HandleFunc("DELETE /api/v1/pins/{pin_id}", handler.DeletePin) // Обернуть в OnlyAuth

	mux.HandleFunc("GET /api/v1/likes/pins", handler.LikedPins)              // Обернуть в OnlyAuth
	mux.HandleFunc("POST /api/v1/likes/{pin_id}/like", handler.CreateLike)   // Обернуть в OnlyAuth
	mux.HandleFunc("DELETE /api/v1/likes/{pin_id}/like", handler.DeleteLike) // Обернуть в OnlyAuth
	mux.HandleFunc("GET /api/v1/likes/{pin_id}/count", handler.LikesCount)
	mux.HandleFunc("GET /api/v1/likes/{pin_id}/users", handler.UsersLiked)

	mux.Handle("GET /img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img"))))
	mux.Handle("GET /docs/swagger.json", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))
	mux.Handle("GET /swagger/", v3.NewHandler("My API", "/docs/swagger.json", "/swagger"))

	server := http.Server{
		Addr:    addr,
		Handler: middleware.CORS(mux),
	}
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
