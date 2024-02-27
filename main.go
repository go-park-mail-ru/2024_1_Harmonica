package main

import (
	h "harmonica/handler"
	"log"
	"net/http"
)

func runServer(addr string) {
	handler, err := h.NewAPIHandler()
	if err != nil {
		log.Print(err)
		return
	}
	defer handler.CloseAPIHandler()

	mux := http.NewServeMux()
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
