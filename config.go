package goit

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	SQLiteDatabasePath     string `envconfig:"TEST_SQLITE_DATABASE_PATH" default:"/tmp/sqlite/testsuite"`
	SQLiteDatabaseInitFile string `envconfig:"TEST_SQLITE_DATABASE_INIT_FILE" default:"example/test/fixture/test_db_init_stmt.sql"`
	DatabaseTruncateFile   string `envconfig:"TEST_DATABASE_TRUNCATE_FILE" default:"example/test/fixture/test_db_truncate_stmt.sql"`
}

func newConfig() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if nil != err {
		return Config{}, err
	}

	return cfg, err
}
