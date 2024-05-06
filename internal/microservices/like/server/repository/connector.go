package repository

import (
	"context"
	"fmt"
	"harmonica/config"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
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

func GetRequestId(ctx context.Context) string {
	if len(metadata.ValueFromIncomingContext(ctx, "request_id")) > 0 {
		return metadata.ValueFromIncomingContext(ctx, "request_id")[0]
	}
	return ""
}

func LogDBQuery(r *DBRepository, ctx context.Context, query string, duration time.Duration) {
	requestId := GetRequestId(ctx)
	r.logger.Info("DB query handled",
		zap.String("request_id", requestId),
		zap.String("query", query),
		zap.String("duration", duration.String()),
	)
}
