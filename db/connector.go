package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DBConf = struct {
	host     string
	port     int
	user     string
	password string
	dbname   string
}{
	host:     "localhost",
	port:     5432,
	user:     "postgres",
	password: "postgres",
	dbname:   "pinterest",
}

type Connector struct {
	db *sqlx.DB
}

func NewConnector() (*Connector, error) {
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DBConf.host, DBConf.port, DBConf.user, DBConf.password, DBConf.dbname)
	db, err := sqlx.Open("postgres", psqlconn)
	return &Connector{db: db}, err
}

func (handler *Connector) Disconnect() error {
	return handler.db.Close()
}
