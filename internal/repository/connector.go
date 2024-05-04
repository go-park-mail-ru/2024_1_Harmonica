package repository

import (
	"fmt"
	"harmonica/config"

	image "harmonica/internal/microservices/image/proto"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Connector struct {
	db *sqlx.DB
	s3 image.ImageClient
}

func NewConnector(conf *config.Config, imageClient image.ImageClient) (*Connector, error) {
	db, err := NewDBConnector(conf.DB)
	if err != nil {
		return &Connector{}, err
	}
	return &Connector{db: db, s3: imageClient}, nil
}

func NewDBConnector(conf config.DBConf) (*sqlx.DB, error) {
	psqlConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		conf.Host, conf.Port, conf.User, conf.Password, conf.DBname)
	return sqlx.Open("postgres", psqlConn)
}

func NewS3Connector(conf config.MinioConf) (*minio.Client, error) {
	endpoint := conf.Endpoint
	accessKeyID := conf.AccessKeyID
	secretAccessKey := conf.SecretAccessKey
	useSSL := conf.UseSSL
	return minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
}

func (connector *Connector) Disconnect() error {
	return connector.db.Close()
}
