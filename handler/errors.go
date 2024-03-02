package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

var (
	ErrReadingRequestBody = errors.New("error reading request body")
	ErrDB                 = errors.New("db error")
)

var HttpStatusByErr = map[error]int{
	ErrReadingRequestBody: 400,
	ErrDB:                 500,
}

func SetHttpError(w http.ResponseWriter, serverErr error, localErr error, status int) {
	fmtMessage := ""
	switch {
	case localErr == nil:
		fmtMessage = fmt.Sprintf("ERROR %s", serverErr.Error())
	default:
		fmtMessage = fmt.Sprintf("ERROR %s: %s", serverErr.Error(), localErr.Error())
	}
	log.Print(fmtMessage)
	w.WriteHeader(HttpStatusByErr[serverErr])
	response, _ := json.Marshal(struct {
		Message string `json:"message"`
	}{Message: serverErr.Error()})
	w.Write(response)
}
