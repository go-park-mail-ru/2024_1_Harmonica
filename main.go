package main

import (
	"harmonica/config"
	"harmonica/db"
	h "harmonica/handler"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
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
	mux.HandleFunc("POST /api/v1/login", handler.Login)
	mux.HandleFunc("POST /api/v1/register", handler.Register)
	mux.HandleFunc("GET /api/v1/logout", handler.Logout)
	mux.HandleFunc("GET /api/v1/is_auth", handler.IsAuth)
	mux.HandleFunc("GET /api/v1/pins_list", handler.PinsList)
	mux.Handle("GET /img/", http.StripPrefix("/img/", http.FileServer(http.Dir("./static/img"))))

	c := cors.New(cors.Options{
		AllowedOrigins:     []string{"http://localhost:8000", "http://89.111.174.111:8000"},
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

func main() {
	runServer(":8080")
}
