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
	mux.HandleFunc("GET /api/v1/pins_list", handler.PinsList)
	mux.Handle("GET /api/v1/img/", http.StripPrefix("/api/v1/img/", http.FileServer(http.Dir("./static/img"))))
	server := http.Server{
		Addr:    addr,
		Handler: mux,
	}
	server.ListenAndServe()
}

func main() {
	runServer(":8080")
}
