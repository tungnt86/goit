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

func NewIphoneProduct(db *sql.DB) *IphoneProduct {
	hitechCategory := category.NewHiTechCategory(db)
	berlinWarehouse := warehouse.NewBerlinWarehouse(db)
	iphoneProduct := &IphoneProduct{BaseFixture: fixture.NewBaseFixture(db)}
	iphoneProduct.SetFoxFixture(iphoneProduct)
	iphoneProduct.SetReference(IphoneProductReference)
	iphoneProduct.AddDependencies(hitechCategory, berlinWarehouse)

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
