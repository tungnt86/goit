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
	i.it.SetupTest()
	db, err := i.initSQLiteDatabase(i.suiteID, i.T().Name())
	must.NotFail(err)
	err = i.setConnectionIntoMap(db)
	must.NotFail(err)
}

func (i *ITsqlite) initSQLiteDatabase(suiteID, testName string) (*sql.DB, error) {
	initStmt, err := ioutil.ReadFile(i.rootDirectory() + "/" + i.config.SQLiteDatabaseInitFile)
	if err != nil {
		return nil, err
	}

	fileName := strings.Replace(testName, testNamePrefix, "", 1)
	dbSubDir, dbFile, err := i.sqliteDatabasePath(suiteID, fileName)
	if err != nil {
		return nil, err
	}
	return database.NewProvider().SQLiteDB(dbSubDir, dbFile, string(initStmt))
}

func (i *ITsqlite) sqliteDatabasePath(suiteID, testName string) (dbSubDir string, dbFile string, err error) {
	suffix, err := i.getRandomString(5)
	if err != nil {
		return "", "", err
	}
	dbFile = testName + "_" + suffix + ".db"
	dbSubDir = i.sqliteDatabaseSubDir(suiteID)

	return dbSubDir, dbFile, nil
}

func (i *ITsqlite) TearDownSuite() {
	i.it.TearDownSuite()
	dbSubDir := i.sqliteDatabaseSubDir(i.suiteID)
	err := database.NewProvider().CleanUpSQLite(dbSubDir)
	must.NotFail(err)
}

func (i *ITsqlite) sqliteDatabaseSubDir(suiteID string) string {
	return i.config.SQLiteDatabasePath + "/" + suiteID
}

func (i *ITsqlite) getRandomString(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
