package goit

import (
	"database/sql"
	"fmt"
	"path/filepath"
	"runtime"
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
	config      Config
	connections map[string]*sql.DB
	foxStore    fixture.FoxStore
}

func (i *it) DB(testName string) (*sql.DB, error) {
	return i.getConnection(testName)
}

func (i *it) FixtureStore() fixture.FoxStore {
	return i.foxStore
}

func (i *it) SetupSuite() {
	cfg, err := newConfig()
	must.NotFail(err)
	i.config = cfg
	i.foxStore = fixture.NewFixtureStore()
}

func (i *it) AfterTest(suiteName, testName string) {
	db, err := i.getConnection(testName)
	must.NotFail(err)
	err = db.Close()
	must.NotFail(err)
}

func (i *it) initConnectionMapIfNeed() {
	mutex.Lock()
	defer mutex.Unlock()
	if i.connections != nil {
		return
	}

	i.connections = make(map[string]*sql.DB)
}

func (i *it) setConnectionIntoMap(testName string, db *sql.DB) error {
	i.initConnectionMapIfNeed()
	_, ok := i.connections[testName]
	if ok {
		return fmt.Errorf("Connection map key conflicts (%s)", testName)
	}
	i.connections[testName] = db
	return nil
}

func (i *it) getConnection(testName string) (*sql.DB, error) {
	db, ok := i.connections[testName]
	if !ok {
		return nil, fmt.Errorf(
			`Connection of test "%s" is not set yet. Is "%s" your test function name?`,
			testName,
			testName,
		)
	}

	return db, nil
}

func (i *it) GetFixture(reference string, testID ...string) (fixture.ModelWithID, error) {
	id := fixture.DefaultTestID
	if len(testID) > 0 {
		id = testID[0]
	}
	return i.foxStore.Get(id, reference)
}

func (i *it) rootDirectory() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}
