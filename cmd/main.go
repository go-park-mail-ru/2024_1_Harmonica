package main

import (
	"github.com/rs/cors"
	"harmonica/config"
	h "harmonica/internal/handler"
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
	mux.HandleFunc("POST /api/v1/register", handler.Register)
	mux.HandleFunc("GET /api/v1/logout", handler.Logout)
	mux.HandleFunc("GET /api/v1/is_auth", handler.IsAuth)
	mux.HandleFunc("GET /api/v1/pins_list", handler.PinsList)
	mux.HandleFunc("POST /api/v1/update_user", handler.UpdateUser)
	mux.Handle("GET /img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img"))))
	mux.Handle("GET /docs/swagger.json", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))
	mux.Handle("GET /swagger/", v3.NewHandler("My API", "/docs/swagger.json", "/swagger"))
	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"http://localhost:8000", "http://85.192.35.36:8000"},
		AllowCredentials:   true,
		AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:     []string{"*"},
		OptionsPassthrough: false,
	})
	server := http.Server{
		Addr:    addr,
		Handler: c.Handler(mux),
	}

	//server := http.Server{
	//	Addr:    addr,
	//	Handler: handler.CorsMiddleware(mux),
	//}

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

// @host		http://85.192.35.36:8080
// @BasePath	api/v1
func main() {
	runServer(":8080")
}
