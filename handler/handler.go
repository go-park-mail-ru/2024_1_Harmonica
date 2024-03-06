package handler

import (
	"harmonica/db"
)

type APIHandler struct {
	Connector db.IConnector
}

func NewAPIHandler(dbConn db.IConnector) *APIHandler {
	return &APIHandler{Connector: dbConn}
}
