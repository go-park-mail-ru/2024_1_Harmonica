package main

import (
	"harmonica/config"
	"harmonica/internal/handler"
	"harmonica/internal/repository"
	"harmonica/internal/service"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	v3 "github.com/swaggest/swgui/v3"
)

func runServer(addr string) {
	conf := config.New()

	dbConn, err := repository.NewConnector(conf.DB)
	if err != nil {
		log.Print(err)
		return
	}
	defer dbConn.Disconnect()

	//useCase := slslls.NewUseCase / NewApp
	r := repository.NewRepository(dbConn)
	s := service.NewService(r)
	h := handler.NewAPIHandler(s)
	//handler := handler2.NewAPIHandler(dbConn)
	mux := http.NewServeMux()

	go handler.CleanupSessions()

	mux.HandleFunc("POST /api/v1/login", h.Login)
	mux.HandleFunc("POST /api/v1/register", h.Register)
	mux.HandleFunc("GET /api/v1/logout", h.Logout)
	mux.HandleFunc("GET /api/v1/is_auth", h.IsAuth)
	mux.HandleFunc("GET /api/v1/pins_list", h.PinsList)
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
