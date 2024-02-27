package handler

import (
	"harmonica/db"
)

type APIHandler struct {
	connector *db.DBConnector
}

func NewAPIHandler(dbConn *db.DBConnector) *APIHandler {
	return &APIHandler{connector: dbConn}
}
