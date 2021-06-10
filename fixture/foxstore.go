package fixture

import (
	"fmt"
	"sync"
)

var (
	mutex = &sync.Mutex{}
)

type FoxStore interface {
	HasReference(reference string) bool
	Set(reference string, fixture ModelWithID) error
	Get(reference string) (ModelWithID, error)
}

type foxStore struct {
	fixtures map[string]ModelWithID
}

func NewFixtureStore() *foxStore {
	return &foxStore{
		fixtures: make(map[string]ModelWithID),
	}
}

func (s *foxStore) Set(reference string, fixture ModelWithID) error {
	mutex.Lock()
	defer mutex.Unlock()
	if s.HasReference(reference) {
		return fmt.Errorf(`fixture "%s" was created already`, reference)
	}
	s.fixtures[reference] = fixture

	return nil
}

func (s *foxStore) Get(reference string) (ModelWithID, error) {
	fixture, ok := s.fixtures[reference]
	if !ok {
		return nil, fmt.Errorf(`fixture "%s" is not created yet`, reference)
	}

	return fixture, nil
}

func (s *foxStore) HasReference(reference string) bool {
	_, ok := s.fixtures[reference]
	return ok
}
