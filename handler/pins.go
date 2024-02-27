package handler

import (
	"log"
	"net/http"
	"strconv"
)

func LogIfError(err error) {
	if err != nil {
		log.Print(err)
	}
}

func (handler *APIHandler) PinsList(w http.ResponseWriter, r *http.Request) {
	userId := r.URL.Query().Get("userId")
	switch userId {
	case "":
		pins, err := handler.connector.GetAllPins()
		if err != nil {
			log.Println(err)
			return
		}
		log.Print(pins)
	default:
		id, err := strconv.Atoi(userId)
		if err != nil {
			log.Println(err)
			return
		}
		pins, err := handler.connector.GetPinsOfUser(id)
		LogIfError(err)
		log.Println(pins)
	}
}
