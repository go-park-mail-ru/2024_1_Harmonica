package repository

import (
	"harmonica/config"

	_ "github.com/lib/pq"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Connector struct {
	s3 *minio.Client
}

func NewConnector(conf *config.Config) (*Connector, error) {
	s3, err := NewS3Connector(conf.Minio)
	if err != nil {
		return &Connector{}, err
	}

	return &Connector{s3: s3}, nil
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
