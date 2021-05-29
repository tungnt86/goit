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

func NewBerlinWarehouse(db *sql.DB) *berlinWarehouse {
	berlinCity := city.NewBerlinCity(db)
	warehouseFixture := &berlinWarehouse{BaseFixture: fixture.NewBaseFixture(db)}
	warehouseFixture.SetFoxFixture(warehouseFixture)
	warehouseFixture.SetReference(BerlinWarehouseReference)
	warehouseFixture.AddDependencies(berlinCity)

	return warehouseFixture
}

func (s *berlinWarehouse) Create() (fixture.ModelWithID, error) {
	berlinCity, err := s.GetFixture(city.BerlinCityReference)
	if err != nil {
		return nil, err
	}
	return newWarehouseFixtureFactory(s.DB()).createWarehouse("Berlin", berlinCity)
}
