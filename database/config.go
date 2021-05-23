package database

import (
	"github.com/kelseyhightower/envconfig"
)

const (
	TestDatabasePrefix = "test_"
)

type Config struct {
	URL             string `envconfig:"TEST_DATABASE_URL" default:"postgres://goit:goit@goit-postgres:5432"`
	Name            string `envconfig:"TEST_DATABASE_NAME" default:"test_goit"`
	SSLMode         string `envconfig:"TEST_DATABASE_SSLMODE" default:"disable"`
	MaxIdleConns    int    `envconfig:"TEST_DATABASE_MAX_IDLE_CONNS" default:"50"`
	MaxOpenConns    int    `envconfig:"TEST_DATABASE_MAX_OPEN_CONNS" default:"100"`
	ConnMaxLifetime int    `envconfig:"TEST_DATABASE_CONN_MAX_LIFETIME" default:"0"`
}

func newConfig() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if nil != err {
		return Config{}, err
	}

	return cfg, err
}
