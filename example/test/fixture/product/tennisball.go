package product

import (
	"database/sql"

	"github.com/tungnt/goit/example/test/fixture/category"
	"github.com/tungnt/goit/example/test/fixture/warehouse"
	"github.com/tungnt/goit/fixture"
)

const (
	TennisBallProductReference = "tennis-ball-product"
)

type tennisBallProduct struct {
	*fixture.BaseFixture
}

func NewTennisBallProduct(db *sql.DB, foxStore fixture.FoxStore) *tennisBallProduct {
	tennisBallFixture := &tennisBallProduct{BaseFixture: fixture.NewBaseFixture(db, foxStore)}
	tennisBallFixture.SetFoxFixture(tennisBallFixture)
	tennisBallFixture.SetReference(TennisBallProductReference)
	tennisBallFixture.AddDependencies(
		category.NewSportCategory(db, foxStore),
		warehouse.NewBerlinWarehouse(db, foxStore),
	)

	return tennisBallFixture
}

func (s *tennisBallProduct) Create() (fixture.ModelWithID, error) {
	sportCategory, err := s.GetFixture(category.SportCategoryReference)
	if err != nil {
		return nil, err
	}
	berlinWarehouse, err := s.GetFixture(warehouse.BerlinWarehouseReference)
	if err != nil {
		return nil, err
	}

	return newProductFixtureFactory(s.DB()).createProduct("Tennis ball", sportCategory, berlinWarehouse)
}
