package handler

import (
	"harmonica/db"
)

type APIHandler struct {
	connector *db.Connector
}

func NewAPIHandler(dbConn *db.Connector) *APIHandler {
	return &APIHandler{connector: dbConn}
}

//type Handler struct {
//	Auth db.Auth
//	Pins db.Pins
//}

//func NewHandler(auth db.Auth, pins db.Pins) *Handler {
//	return &Handler{
//		Auth: auth,
//		Pins: pins,
//	}
//}
