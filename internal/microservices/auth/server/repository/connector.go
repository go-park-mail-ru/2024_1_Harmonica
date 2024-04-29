package repository

import (
	"context"
	"fmt"
	"harmonica/config"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type Connector struct {
	db *sqlx.DB
}

func NewConnector(conf *config.Config) (*Connector, error) {
	db, err := NewDBConnector(conf.DB)
	if err != nil {
		return &Connector{}, err
	}

	return &Connector{db: db}, nil
}

func NewDBConnector(conf config.DBConf) (*sqlx.DB, error) {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.DBname)
	return sqlx.Open("postgres", psqlConn)
}

func (connector *Connector) Disconnect() error {
	return connector.db.Close()
}

func LogDBQuery(r *DBRepository, ctx context.Context, query string, duration time.Duration) {
	requestId := ctx.Value("request_id").(string)
	r.logger.Info("DB query handled",
		zap.String("request_id", requestId),
		zap.String("query", query),
		zap.String("duration", duration.String()),
	)
}
