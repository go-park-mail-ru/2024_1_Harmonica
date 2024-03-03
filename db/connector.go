package db

import (
	"fmt"

	"harmonica/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBConnector struct {
	db *sqlx.DB
}

func NewConnector(conf config.DBConf) (*DBConnector, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.DBname)
	db, err := sqlx.Open("postgres", psqlconn)
	return &DBConnector{db: db}, err
}

func (connector *DBConnector) Disconnect() error {
	return connector.db.Close()
}
