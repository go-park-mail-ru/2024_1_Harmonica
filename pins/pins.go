package pins

import (
	pq "harmonica/db/init/postgres"
	"log"
	"net/http"
	"strconv"
)

func LogIfError(err error) {
	if err != nil {
		log.Print(err)
	}
}

func PinsList(handler *pq.APIHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.URL.Query().Get("userId")
		switch userId {
		case "":
			pins, err := handler.GetAllPins()
			if err != nil {
				log.Println(err)
				return
			}
			log.Print(pins)
		default:
			id, err := strconv.Atoi(userId)
			if err != nil {
				log.Println(err)
			}
			pins, err := handler.GetPinsOfUser(id)
			LogIfError(err)
			log.Println(pins)
		}
	}
}
