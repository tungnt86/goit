package city

import (
	"database/sql"

	"github.com/tungnt/goit/fixture"
)

const (
	BerlinCityReference = "berlin-city"
)

type berlinCity struct {
	*fixture.BaseFixture
}

func NewBerlinCity(db *sql.DB, foxStore fixture.FoxStore) *berlinCity {
	cityFixture := &berlinCity{BaseFixture: fixture.NewBaseFixture(db, foxStore)}
	cityFixture.SetFoxFixture(cityFixture)
	cityFixture.SetReference(BerlinCityReference)

	return cityFixture
}

func (s *berlinCity) Create() (fixture.ModelWithID, error) {
	return newCityFixtureFactory(s.DB()).createCity("Smartphone")
}
