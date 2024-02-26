package postgres

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type APIHandler struct {
	db *sqlx.DB
}

func NewAPIHandler() (*APIHandler, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DBConf.host, DBConf.port, DBConf.user, DBConf.password, DBConf.dbname)
	db, err := sqlx.Open("postgres", psqlconn)
	return &APIHandler{db: db}, err
}

func (handler *APIHandler) Disconnect() error {
	return handler.db.Close()
}
