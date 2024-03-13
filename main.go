package main

import (
	"harmonica/config"
	"harmonica/db"
	h "harmonica/handler"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
	v3 "github.com/swaggest/swgui/v3"
)

func runServer(addr string) {
	conf := config.New()

	dbConn, err := db.NewConnector(conf.DB)
	if err != nil {
		log.Print(err)
		return
	}
	defer dbConn.Disconnect()
	handler := h.NewAPIHandler(dbConn)
	mux := http.NewServeMux()

	go h.CleanupSessions()

	mux.HandleFunc("POST /api/v1/login", handler.Login)
	mux.HandleFunc("POST /api/v1/register", handler.Register)
	mux.HandleFunc("GET /api/v1/logout", handler.Logout)
	mux.HandleFunc("GET /api/v1/is_auth", handler.IsAuth)
	mux.HandleFunc("GET /api/v1/pins_list", handler.PinsList)
	mux.Handle("GET /docs/swagger.json", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))
	mux.Handle("GET /swagger/", v3.NewHandler("My API", "/docs/swagger.json", "/swagger"))
	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"https://127.0.0.1", "https://85.192.35.36"},
		AllowCredentials:   true,
		AllowedMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:     []string{"*"},
		OptionsPassthrough: false,
	})

	server := http.Server{
		Addr:    addr,
		Handler: c.Handler(mux),
	}
	server.ListenAndServeTLS("cert.pem", "key.pem")
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
