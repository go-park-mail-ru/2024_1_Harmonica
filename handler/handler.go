package handler

import (
	"harmonica/db"
)

type APIHandler struct {
	connector *db.Connector
}

func NewAPIHandler() (*APIHandler, error) {
	conn, err := db.NewConnector()
	return &APIHandler{connector: conn}, err
}

func (handler *APIHandler) CloseAPIHandler() {
	handler.connector.Disconnect()
}
