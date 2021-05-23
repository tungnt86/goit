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

func (s *ProductRepoTestSuite) BeforeTest(suiteName, testName string) {
	s.ITsqlite.BeforeTest(suiteName, testName)
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
		ID:          1,
		Name:        "Tennis ball",
		CategoryID:  sportCategory.(*model.Category).ID,
		WarehouseID: berlinWarehouse.(*model.Warehouse).ID,
	}
	actualResult, err := s.repo.GetOne(context.Background(), 1)
	s.NoError(err)
	s.Equal(expectedResult, actualResult)
}

func (s *ProductRepoTestSuite) TestGetOneInParallel_NoError() {
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
				product, err := s.GetFixture(product.IphoneProductReference, testID)
				s.NoError(err)
				return product.(*model.Product)
			},
			want: func(testID string) *model.Product {
				product, err := s.GetFixture(product.IphoneProductReference, testID)
				s.NoError(err)
				hitechCategory, err := s.GetFixture(category.HiTechCategoryReference, testID)
				s.NoError(err)
				berlinWarehouse, err := s.GetFixture(warehouse.BerlinWarehouseReference, testID)
				s.NoError(err)
				return &model.Product{
					ID:          product.(*model.Product).ID,
					Name:        "Iphone 10",
					CategoryID:  hitechCategory.(*model.Category).ID,
					WarehouseID: berlinWarehouse.(*model.Warehouse).ID,
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
				product, err := s.GetFixture(product.TennisBallProductReference, testID)
				s.NoError(err)
				sportCategory, err := s.GetFixture(category.SportCategoryReference, testID)
				s.NoError(err)
				berlinWarehouse, err := s.GetFixture(warehouse.BerlinWarehouseReference, testID)
				s.NoError(err)
				return &model.Product{
					ID:          product.(*model.Product).ID,
					Name:        "Tennis ball",
					CategoryID:  sportCategory.(*model.Category).ID,
					WarehouseID: berlinWarehouse.(*model.Warehouse).ID,
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
