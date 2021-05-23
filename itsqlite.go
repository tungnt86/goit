package goit

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"io/ioutil"

	"github.com/tungnt/goit/database"
	"github.com/tungnt/goit/must"
)

type ITsqlite struct {
	it
	suiteName string
}

func (i *ITsqlite) BeforeTest(suiteName, testName string) {
	if i.suiteName == "" {
		i.suiteName = suiteName
	}
	db, err := i.initSQLiteDatabase(suiteName, testName)
	must.NotFail(err)
	err = i.setConnectionIntoMap(db)
	must.NotFail(err)
}

func (i *ITsqlite) initSQLiteDatabase(suiteName, testName string) (*sql.DB, error) {
	rootDir := i.rootDirectory()
	initStmt, err := ioutil.ReadFile(rootDir + "/" + i.config.SQLiteDatabaseInitFile)
	if err != nil {
		return nil, err
	}

	dbSubDir, dbFile, err := i.sqliteDatabasePath(suiteName, testName)
	if err != nil {
		return nil, err
	}
	return database.NewProvider().SQLiteDB(dbSubDir, dbFile, string(initStmt))
}

func (i *ITsqlite) sqliteDatabasePath(suiteName, testName string) (dbSubDir string, dbFile string, err error) {
	bytes := make([]byte, 5)
	_, err = rand.Read(bytes)
	if err != nil {
		return "", "", err
	}
	suffix := base64.URLEncoding.EncodeToString(bytes)
	dbFile = testName + "_" + suffix + ".db"
	dbSubDir = i.sqliteDatabaseSubDir(suiteName)

	return dbSubDir, dbFile, nil
}

func (i *ITsqlite) TearDownSuite() {
	i.it.TearDownSuite()
	dbSubDir := i.sqliteDatabaseSubDir(i.suiteName)
	err := database.NewProvider().CleanUpSQLite(dbSubDir)
	must.NotFail(err)
}

func (i *ITsqlite) sqliteDatabaseSubDir(suiteName string) string {
	return i.rootDirectory() + "/" + i.config.SQLiteDatabasePath + "/" + suiteName
}
