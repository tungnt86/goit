package fixture

import (
	"database/sql"
	"errors"
)

type FoxFixture interface {
	DB() *sql.DB
	Reference() (string, error)
	SetReference(reference string)
	Create() (ModelWithID, error)
}

type FoxBuilder interface {
	AddDependencies(dependencies ...FoxBuilder)
	Build() error
}

type BaseFixture struct {
	FoxFixture
	db           *sql.DB
	foxStore     FoxStore
	reference    string
	dependencies []FoxBuilder
}

func NewBaseFixture(db *sql.DB, foxStore FoxStore) *BaseFixture {
	return &BaseFixture{
		db:       db,
		foxStore: foxStore,
	}
}

func (b *BaseFixture) SetFoxFixture(foxFixture FoxFixture) {
	b.FoxFixture = foxFixture
}

func (b *BaseFixture) DB() *sql.DB {
	return b.db
}

func (b *BaseFixture) Reference() (string, error) {
	if b.reference == "" {
		return "", errors.New("reference of this fixture is not defined yet")
	}

	return b.reference, nil
}

func (b *BaseFixture) SetReference(reference string) {
	b.reference = reference
}

func (b *BaseFixture) AddDependencies(dependencies ...FoxBuilder) {
	b.dependencies = append(b.dependencies, dependencies...)
}

func (b *BaseFixture) BuildDependencies() error {
	for _, dependency := range b.dependencies {
		err := dependency.Build()
		if nil != err {
			return err
		}
	}

	return nil
}

func (b *BaseFixture) Build() error {
	_, err := b.GetFixture(b.reference)
	if err == nil {
		return nil
	}
	err = b.BuildDependencies()
	if nil != err {
		return err
	}
	fixture, err := b.Create()
	if nil != err {
		return err
	}

	return b.foxStore.Set(b.reference, fixture)
}

func (b *BaseFixture) GetFixture(reference string) (ModelWithID, error) {
	fixture, err := b.foxStore.Get(reference)
	if err != nil {
		return nil, err
	}

	return fixture, nil
}
