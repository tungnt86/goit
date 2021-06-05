package database

import (
	"database/sql"
	"errors"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/tungnt/goit/must"
)

var (
	once     sync.Once
	mutex    = &sync.Mutex{}
	instance *provider
)

type Provider interface {
	DB() (*sql.DB, error)
	CleanUpDB(db *sql.DB, truncateStmt string) error
	SQLiteDB(dbFilePath string, initStmt string) (*sql.DB, error)
	CleanUpSQLite(dbSubDir string) error
}

type provider struct {
	cfg Config
}

func NewProvider() Provider {
	once.Do(func() {
		cfg, err := newConfig()
		must.NotFail(err)
		if !strings.Contains(cfg.Name, TestDatabasePrefix) {
			panic(errors.New(`Test database name must have prefix "test_" to differentiate from operational database`))
		}
		instance = &provider{
			cfg: cfg,
		}
	})

	return instance
}

/*
We do not need to initiate test database because usually we already have a test database
Migration for test database should be done manually and it is not in scope of integration test
*/
func (p *provider) DB() (*sql.DB, error) {
	mutex.Lock()
	defer mutex.Unlock()
	return connect(p.cfg)
}

func connect(cfg Config) (*sql.DB, error) {
	dbUrl, err := url.Parse(cfg.URL)
	if err != nil {
		return nil, err
	}
	dataSource, _ := dbUrl.Parse("/" + cfg.Name + "?sslmode=" + cfg.SSLMode)
	sqlDB, err := sql.Open("postgres", dataSource.String())
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(time.Duration(cfg.ConnMaxLifetime))

	return sqlDB, nil
}

func (p *provider) CleanUpDB(db *sql.DB, truncateStmt string) error {
	mutex.Lock()
	defer mutex.Unlock()
	_, err := db.Exec(truncateStmt)
	if err != nil {
		return err
	}

	return nil
}

/*
SQLite database is a file based database
Create a new SQLite database file and init tables from a dump file
*/
func (p *provider) SQLiteDB(dbFilePath string, initStmt string) (*sql.DB, error) {
	mutex.Lock()
	defer mutex.Unlock()
	os.Remove(dbFilePath)
	file, err := os.Create(dbFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	db, err := sql.Open("sqlite3", dbFilePath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(initStmt)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (p *provider) CleanUpSQLite(dbFilePath string) error {
	mutex.Lock()
	defer mutex.Unlock()
	return os.Remove(dbFilePath)
}
