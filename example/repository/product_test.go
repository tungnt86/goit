package repository

import (
	"context"
	"database/sql"
	"testing"

	"github.com/tungnt/goit"
	"github.com/tungnt/goit/example/model"
	"github.com/tungnt/goit/example/test/fixture/category"
	"github.com/tungnt/goit/example/test/fixture/product"
	"github.com/tungnt/goit/example/test/fixture/warehouse"
)

type ProductRepoTestSuite struct {
	goit.ITsqlite
	repo *productRepo
}

func (s *ProductRepoTestSuite) SetupTest() {
	s.ITsqlite.SetupTest()
	s.repo = &productRepo{db: s.DB()}
}

func (s *ProductRepoTestSuite) TestGetOne_NoError() {
	err := product.NewTennisBallProduct(s.DB()).Build()
	s.NoError(err)
	sportCategory, err := s.GetFixture(category.SportCategoryReference)
	s.NoError(err)
	berlinWarehouse, err := s.GetFixture(warehouse.BerlinWarehouseReference)
	s.NoError(err)
	expectedResult := &model.Product{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Tennis ball",
		CategoryID:  sportCategory.GetID(),
		WarehouseID: berlinWarehouse.GetID(),
	}
	actualResult, err := s.repo.GetOne(context.Background(), 1)
	s.NoError(err)
	s.Equal(expectedResult, actualResult)
}

func (s *ProductRepoTestSuite) TestGetOneInParallel_NoError() {
	// A lint rule requires to define parallel test here
	s.T().Parallel()
	// Call SetupTest here to init the test if the test is run in parallel.
	s.SetupTest()
	tests := []struct {
		name    string
		fixture func(testID string, db *sql.DB) *model.Product
		want    func(testID string) *model.Product
	}{
		{
			name: "test get one iphone 10",
			fixture: func(testID string, db *sql.DB) *model.Product {
				err := product.NewIphoneProduct(db).Build(testID)
				s.NoError(err)
				iphoneProduct, err := s.GetFixture(product.IphoneProductReference, testID)
				s.NoError(err)
				return iphoneProduct.(*model.Product)
			},
			want: func(testID string) *model.Product {
				iphoneProduct, err := s.GetFixture(product.IphoneProductReference, testID)
				s.NoError(err)
				hitechCategory, err := s.GetFixture(category.HiTechCategoryReference, testID)
				s.NoError(err)
				berlinWarehouse, err := s.GetFixture(warehouse.BerlinWarehouseReference, testID)
				s.NoError(err)
				return &model.Product{
					BaseModel:   model.BaseModel{ID: iphoneProduct.GetID()},
					Name:        "Iphone 10",
					CategoryID:  hitechCategory.GetID(),
					WarehouseID: berlinWarehouse.GetID(),
				}
			},
		},
		{
			name: "test get one tennis ball",
			fixture: func(testID string, db *sql.DB) *model.Product {
				err := product.NewTennisBallProduct(db).Build(testID)
				s.NoError(err)
				product, err := s.GetFixture(product.TennisBallProductReference, testID)
				s.NoError(err)
				return product.(*model.Product)
			},
			want: func(testID string) *model.Product {
				tennisProduct, err := s.GetFixture(product.TennisBallProductReference, testID)
				s.NoError(err)
				sportCategory, err := s.GetFixture(category.SportCategoryReference, testID)
				s.NoError(err)
				berlinWarehouse, err := s.GetFixture(warehouse.BerlinWarehouseReference, testID)
				s.NoError(err)
				return &model.Product{
					BaseModel:   model.BaseModel{ID: tennisProduct.GetID()},
					Name:        "Tennis ball",
					CategoryID:  sportCategory.GetID(),
					WarehouseID: berlinWarehouse.GetID(),
				}
			},
		},
	}

	for id := range tests {
		db := s.DB()
		test := tests[id]
		s.T().Run(test.name, func(t *testing.T) {
			t.Parallel()
			product := test.fixture(test.name, db)
			actualResult, err := s.repo.GetOne(context.Background(), product.ID)
			s.NoError(err)
			s.Equal(test.want(test.name), actualResult)
		})
	}
}
