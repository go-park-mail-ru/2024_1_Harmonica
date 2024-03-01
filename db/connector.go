package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBConf struct {
	Host     string
	Port     int
	User     string
	Password string
	Dbname   string
}

type DBConnector struct {
	db *sqlx.DB
}

func NewConnector(conf DBConf) (*DBConnector, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Dbname)
	db, err := sqlx.Open("postgres", psqlconn)
	return &DBConnector{db: db}, err
}

func (handler *DBConnector) Disconnect() error {
	return handler.db.Close()
}
