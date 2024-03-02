package main

import (
	"harmonica/db"
	h "harmonica/handler"
	"log"
	"net/http"
)

func runServer(addr string) {
	dbConn, err := db.NewConnector(Conf)
	if err != nil {
		log.Print(err)
		return
	}
	defer dbConn.Disconnect()
	handler := h.NewAPIHandler(dbConn)
	mux := http.NewServeMux()
	mux.HandleFunc("POST /login", handler.Login)
	mux.HandleFunc("POST /logout", handler.Logout) // POST или GET все-таки?
	mux.HandleFunc("POST /register", handler.Register)
	mux.HandleFunc("GET /pinsList", handler.PinsList)
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	server.ListenAndServe()
}

func main() {
	runServer(":8080")
}
