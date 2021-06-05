package fixture

import (
	"fmt"
	"sync"
)

const (
	DefaultTestID = "0"
)

var (
	mutex = &sync.Mutex{}
)

type FoxStore interface {
	HasReference(runID, reference string) bool
	Set(runID, reference string, fixture ModelWithID) error
	Get(runID, reference string) (ModelWithID, error)
}

type foxStore struct {
	fixtures map[string]map[string]ModelWithID
}

func NewFixtureStore() *foxStore {
	return &foxStore{
		fixtures: make(map[string]map[string]ModelWithID),
	}
}

func (s *foxStore) Set(testID, reference string, fixture ModelWithID) error {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := s.fixtures[testID]; !ok {
		s.fixtures[testID] = make(map[string]ModelWithID)
	}
	if s.HasReference(testID, reference) {
		return fmt.Errorf(`Reference "%s" is set already in test "%s"`, reference, testID)
	}
	s.fixtures[testID][reference] = fixture

	return nil
}

func (s *foxStore) Get(testID, reference string) (ModelWithID, error) {
	if _, ok := s.fixtures[testID]; !ok {
		return nil, fmt.Errorf(`Reference "%s" is not set in test "%s" yet`, reference, testID)
	}
	fixture, ok := s.fixtures[testID][reference]
	if !ok {
		return nil, fmt.Errorf(`Reference "%s" is not set in test "%s" yet.`, reference, testID)
	}

	return fixture, nil
}

func (s *foxStore) HasReference(testID, reference string) bool {
	if _, ok := s.fixtures[testID]; !ok {
		return false
	}
	_, ok := s.fixtures[testID][reference]
	return ok
}
