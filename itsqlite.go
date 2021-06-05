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
	suiteID string
}

func (i *ITsqlite) SetupSuite() {
	i.it.SetupSuite()
	suiteID, err := i.getRandomString(5)
	must.NotFail(err)
	i.suiteID = suiteID
}

func (i *ITsqlite) BeforeTest(suiteName, testName string) {
	db, err := i.initSQLiteDatabase(suiteName, testName)
	must.NotFail(err)
	err = i.setConnectionIntoMap(testName, db)
	must.NotFail(err)
}

func (i *ITsqlite) AfterTest(suiteName, testName string) {
	i.it.AfterTest(suiteName, testName)
	dbFilePath := i.sqliteDatabaseFilePath(suiteName, testName)
	err := database.NewProvider().CleanUpSQLite(dbFilePath)
	must.NotFail(err)
}

func (i *ITsqlite) initSQLiteDatabase(suiteName, testName string) (*sql.DB, error) {
	dbFilePath := i.sqliteDatabaseFilePath(suiteName, testName)
	initStmt, err := ioutil.ReadFile(i.rootDirectory() + "/" + i.config.SQLiteDatabaseInitFile)
	if err != nil {
		return nil, err
	}
	return database.NewProvider().SQLiteDB(dbFilePath, string(initStmt))
}

func (i *ITsqlite) sqliteDatabaseFilePath(suiteName, testName string) string {
	dbFile := i.suiteID + "_" + suiteName + "_" + testName + ".db"
	dbPath := i.config.SQLiteDatabasePath + "/" + dbFile
	return dbPath
}

func (i *ITsqlite) getRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
