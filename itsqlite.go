package goit

import (
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"io/ioutil"
	"strings"

	"github.com/tungnt/goit/database"
	"github.com/tungnt/goit/must"
)

const testNamePrefix = "TestSuite/"

type ITsqlite struct {
	it
	suiteID string
}

func (i *ITsqlite) SetupSuite() {
	i.it.SetupSuite()
	suiteID, err := i.getRandomString(10)
	must.NotFail(err)
	i.suiteID = suiteID
}

func (i *ITsqlite) SetupTest() {
	db, err := i.initSQLiteDatabase()
	must.NotFail(err)
	err = i.setConnectionIntoMap(db)
	must.NotFail(err)
}

func (i *ITsqlite) initSQLiteDatabase() (*sql.DB, error) {
	dbFilePath := i.sqliteDatabaseFilePath()
	initStmt, err := ioutil.ReadFile(i.rootDirectory() + "/" + i.config.SQLiteDatabaseInitFile)
	if err != nil {
		return nil, err
	}
	return database.NewProvider().SQLiteDB(dbFilePath, string(initStmt))
}

func (i *ITsqlite) sqliteDatabaseFilePath() string {
	testName := strings.Replace(i.T().Name(), testNamePrefix, "", 1)
	dbFile := i.suiteID + "_" + testName + ".db"
	dbPath := i.config.SQLiteDatabasePath + "/" + dbFile
	return dbPath
}

func (i *ITsqlite) TearDownTest() {
	i.it.TearDownTest()
	dbFilePath := i.sqliteDatabaseFilePath()
	err := database.NewProvider().CleanUpSQLite(dbFilePath)
	must.NotFail(err)
}

func (i *ITsqlite) getRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
