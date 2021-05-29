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

func (i *it) DB() *sql.DB {
	return i.connections[i.T().Name()]
}

func (i *it) SetupSuite() {
	cfg, err := newConfig()
	must.NotFail(err)
	i.config = cfg
}

func (i *it) SetupTest() {
	i.initFixtureStore()
}

func (i *it) initConnectionMapIfNeed() {
	mutex.Lock()
	defer mutex.Unlock()
	if i.connections != nil {
		return
	}

	i.connections = make(map[string]*sql.DB)
}

func (i *it) setConnectionIntoMap(db *sql.DB) error {
	i.initConnectionMapIfNeed()
	_, ok := i.connections[i.T().Name()]
	if ok {
		return fmt.Errorf("Connection map key conflicts (%s)", i.T().Name())
	}
	i.connections[i.T().Name()] = db
	return nil
}

func (i *it) getConnection() (*sql.DB, error) {
	db, ok := i.connections[i.T().Name()]
	if !ok {
		return nil, fmt.Errorf("Connection of test (%s) is not set yet", i.T().Name())
	}

	return db, nil
}

func (i *it) TearDownSuite() {
	for _, db := range i.connections {
		err := db.Close()
		must.NotFail(err)
	}
	i.connections = make(map[string]*sql.DB)
}

func (i *it) TearDownTest() {
	i.foxStore.Reset()
}

func (i *it) initFixtureStore() {
	i.foxStore = fixture.NewFixtureStore()
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
