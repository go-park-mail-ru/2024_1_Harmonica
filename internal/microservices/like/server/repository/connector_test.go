package repository

import (
	"github.com/stretchr/testify/assert"
	"harmonica/config"
	"testing"
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
	conn, err := NewConnector(conf)
	assert.NotNil(t, conn)
	assert.NoError(t, err)
	err = conn.Disconnect()
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
	db, err := NewDBConnector(conf)
	assert.NotNil(t, db)
	assert.NoError(t, err)
}
