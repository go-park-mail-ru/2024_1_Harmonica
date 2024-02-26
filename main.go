package main

import (
	pq "harmonica/db/init/postgres"
	"harmonica/pins"
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {}

func logout(w http.ResponseWriter, r *http.Request) {}

func register(w http.ResponseWriter, r *http.Request) {}

func runServer(addr string) {
	handler, _ := pq.NewAPIHandler() // Надо что-то делать, если мы не смогли в connect к БД
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", login)
	mux.HandleFunc("POST /register", register)
	mux.HandleFunc("GET /logout", logout)
	mux.HandleFunc("GET /pinsList", pins.PinsList(handler))

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	server.ListenAndServe()
}

func main() {
	runServer(":8080")
}
