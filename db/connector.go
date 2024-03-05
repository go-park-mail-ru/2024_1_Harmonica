package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"harmonica/config"
)

type Connector struct {
	db *sqlx.DB
}

func NewConnector(conf config.DBConf) (*Connector, error) {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.DBname)
	db, err := sqlx.Open("postgres", psqlConn)
	return &Connector{db: db}, err
}

func (connector *Connector) Disconnect() error {
	return connector.db.Close()
}
