package goit

import (
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/stretchr/testify/suite"

	"github.com/tungnt/goit/fixture"
	"github.com/tungnt/goit/must"
)

var (
	mutex = &sync.Mutex{}
)

type it struct {
	suite.Suite
	config Config
	dbMap  map[string]*sql.DB
}

func (i *it) GetCurrentTestDB() (*sql.DB, error) {
	testName, err := i.getCurrentTestFunctionName()
	if err != nil {
		return nil, err
	}
	return i.getDB(testName)
}

func (i *it) getCurrentTestFunctionName() (string, error) {
	pc, _, _, ok := runtime.Caller(2)
	if !ok {
		return "", errors.New("could not get test function name")
	}
	testFunction := runtime.FuncForPC(pc).Name()
	lastPoint := strings.LastIndex(testFunction, ".")
	return testFunction[lastPoint+1:], nil
}

func (i *it) getDB(testName string) (*sql.DB, error) {
	db, ok := i.dbMap[testName]
	if !ok {
		return nil, fmt.Errorf(
			`database of test "%s" is not created yet. Is "%s" your test function name?`,
			testName,
			testName,
		)
	}

	return db, nil
}

func (i *it) NewFixtureStore() fixture.FoxStore {
	return fixture.NewFixtureStore()
}

func (i *it) GetFixture(foxStore fixture.FoxStore, reference string) (fixture.ModelWithID, error) {
	return foxStore.Get(reference)
}

func (i *it) SetupSuite() {
	cfg, err := newConfig()
	must.NotFail(err)
	i.config = cfg
}

func (i *it) AfterTest(suiteName, testName string) {
	db, err := i.getDB(testName)
	must.NotFail(err)
	err = db.Close()
	must.NotFail(err)
}

func (i *it) initDBMapIfNeeded() {
	mutex.Lock()
	defer mutex.Unlock()
	if i.dbMap != nil {
		return
	}

	i.dbMap = make(map[string]*sql.DB)
}

func (i *it) setDBIntoMap(testName string, db *sql.DB) error {
	i.initDBMapIfNeeded()
	_, ok := i.dbMap[testName]
	if ok {
		return fmt.Errorf(`database of test "%s" is created already. Is "%s" your test function name?`, testName, testName)
	}
	i.dbMap[testName] = db
	return nil
}

func (i *it) rootDirectory() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}
