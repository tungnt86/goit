package warehouse

import (
	"database/sql"

	"github.com/tungnt/goit/example/test/fixture/city"
	"github.com/tungnt/goit/fixture"
)

const (
	BerlinWarehouseReference = "berlin-warehouse"
)

type berlinWarehouse struct {
	*fixture.BaseFixture
}

func NewBerlinWarehouse(db *sql.DB, foxStore fixture.FoxStore) *berlinWarehouse {
	warehouseFixture := &berlinWarehouse{BaseFixture: fixture.NewBaseFixture(db, foxStore)}
	warehouseFixture.SetFoxFixture(warehouseFixture)
	warehouseFixture.SetReference(BerlinWarehouseReference)
	warehouseFixture.AddDependencies(city.NewBerlinCity(db, foxStore))

	return warehouseFixture
}

func (s *berlinWarehouse) Create() (fixture.ModelWithID, error) {
	berlinCity, err := s.GetFixture(city.BerlinCityReference)
	if err != nil {
		return nil, err
	}
	return newWarehouseFixtureFactory(s.DB()).createWarehouse("Berlin", berlinCity)
}
