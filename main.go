package main

import (
	"net/http"
)

func login(w http.ResponseWriter, r *http.Request) {}

func logout(w http.ResponseWriter, r *http.Request) {}

func register(w http.ResponseWriter, r *http.Request) {}

func pinsList(w http.ResponseWriter, r *http.Request) {}

func runServer(addr string) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", login)
	mux.HandleFunc("POST /register", register)
	mux.HandleFunc("GET /logout", logout)
	mux.HandleFunc("GET /pinsList", pinsList)

	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}

	server.ListenAndServe()
}

func main() {
	runServer(":8080")
}
