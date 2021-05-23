package fixture

import (
	"fmt"
	"sync"
)

const (
	DefaultTestID = "0"
)

var (
	once     sync.Once
	mutex    = &sync.Mutex{}
	instance *fixtureStore
)

type FixtureStore interface {
	HasReference(runID, reference string) bool
	Set(runID, reference string, fixture interface{}) error
	Get(runID, reference string) (interface{}, error)
	Reset()
}

type fixtureStore struct {
	fixtures map[string]map[string]interface{}
}

func NewFixtureStore() *fixtureStore {
	once.Do(func() {
		instance = &fixtureStore{
			fixtures: make(map[string]map[string]interface{}),
		}
	})

	return instance
}

func (s *fixtureStore) initStoreForTest(testID string) {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := s.fixtures[testID]; ok {
		return
	}
	s.fixtures[testID] = make(map[string]interface{})
}

func (s *fixtureStore) Set(testID, reference string, fixture interface{}) error {
	if _, ok := s.fixtures[testID]; !ok {
		s.initStoreForTest(testID)
	}
	if _, ok := s.fixtures[testID][reference]; ok {
		return fmt.Errorf(`Reference "%s" is set already in test "%s"`, reference, testID)
	}
	s.fixtures[testID][reference] = fixture

	return nil
}

func (s *fixtureStore) Get(testID, reference string) (interface{}, error) {
	if _, ok := s.fixtures[testID]; !ok {
		return nil, fmt.Errorf(`Reference "%s" is not set in test "%s" yet`, reference, testID)
	}
	fixture, ok := s.fixtures[testID][reference]
	if !ok {
		return nil, fmt.Errorf(`Reference "%s" is not set in test "%s" yet.`, reference, testID)
	}

	return fixture, nil
}

func (s *fixtureStore) HasReference(testID, reference string) bool {
	if _, ok := s.fixtures[testID]; !ok {
		return false
	}
	_, ok := s.fixtures[testID][reference]
	return ok
}

func (s *fixtureStore) Reset() {
	mutex.Lock()
	defer mutex.Unlock()
	s.fixtures = make(map[string]map[string]interface{})
}
