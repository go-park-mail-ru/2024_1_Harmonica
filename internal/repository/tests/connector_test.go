package test_repository

import (
	"github.com/golang/mock/gomock"
	"harmonica/config"
	"harmonica/internal/repository"
	mock_proto "harmonica/mocks/microservices/image/proto"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_NewConnector(t *testing.T) {
	conf := &config.Config{
		DB: config.DBConf{
			Host:     "localhost",
			Port:     5432,
			User:     "user",
			Password: "password",
			DBname:   "dbname",
		},
	}
	ctrl := gomock.NewController(t)
	imageClient := mock_proto.NewMockImageClient(ctrl)
	conn, err := repository.NewConnector(conf, imageClient)
	assert.NotNil(t, conn)
	assert.NoError(t, err)
}

func TestRepository_NewDBConnector(t *testing.T) {
	conf := config.DBConf{
		Host:     "localhost",
		Port:     5432,
		User:     "user",
		Password: "password",
		DBname:   "dbname",
	}
	db, err := repository.NewDBConnector(conf)
	assert.NotNil(t, db)
	assert.NoError(t, err)
}

func TestRepository_NewS3Connector(t *testing.T) {
	conf := config.MinioConf{
		Endpoint:        "localhost:9000",
		AccessKeyID:     "accessKey",
		SecretAccessKey: "secretKey",
		UseSSL:          false,
	}
	s3Client, err := repository.NewS3Connector(conf)
	assert.NotNil(t, s3Client)
	assert.NoError(t, err)
}

func TestRepository_Disconnect(t *testing.T) {
	conf := &config.Config{
		DB: config.DBConf{
			Host:     "localhost",
			Port:     5432,
			User:     "user",
			Password: "password",
			DBname:   "dbname",
		},
	}
	ctrl := gomock.NewController(t)
	imageClient := mock_proto.NewMockImageClient(ctrl)
	conn, _ := repository.NewConnector(conf, imageClient)
	err := conn.Disconnect()
	assert.NoError(t, err)
}
