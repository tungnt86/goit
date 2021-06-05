package product

import (
	"database/sql"

	"github.com/tungnt/goit/example/test/fixture/category"
	"github.com/tungnt/goit/example/test/fixture/warehouse"
	"github.com/tungnt/goit/fixture"
)

const (
	IphoneProductReference = "iphone-product"
)

type IphoneProduct struct {
	*fixture.BaseFixture
}

func NewIphoneProduct(db *sql.DB, foxStore fixture.FoxStore) *IphoneProduct {
	iphoneProduct := &IphoneProduct{BaseFixture: fixture.NewBaseFixture(db, foxStore)}
	iphoneProduct.SetFoxFixture(iphoneProduct)
	iphoneProduct.SetReference(IphoneProductReference)
	iphoneProduct.AddDependencies(
		category.NewHiTechCategory(db, foxStore),
		warehouse.NewBerlinWarehouse(db, foxStore),
	)

	return iphoneProduct
}

func (s *IphoneProduct) Create() (fixture.ModelWithID, error) {
	hiTechCategory, err := s.GetFixture(category.HiTechCategoryReference)
	if err != nil {
		return nil, err
	}
	berlinWarehouse, err := s.GetFixture(warehouse.BerlinWarehouseReference)
	if err != nil {
		return nil, err
	}

	return newProductFixtureFactory(s.DB()).createProduct("Iphone 10", hiTechCategory, berlinWarehouse)
}
