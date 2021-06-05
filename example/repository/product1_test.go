package repository

import (
	"context"
	"database/sql"
	"sync"
	"testing"

	"github.com/tungnt/goit"
	"github.com/tungnt/goit/example/model"
	"github.com/tungnt/goit/example/test/fixture/category"
	"github.com/tungnt/goit/example/test/fixture/product"
	"github.com/tungnt/goit/example/test/fixture/warehouse"
	"github.com/tungnt/goit/fixture"
)

type ProductRepoTestSuite1 struct {
	goit.ITsqlite
}

func (s *ProductRepoTestSuite1) TestGetOne_NoError() {
	db := s.DB()
	testName := s.T().Name()
	s.T().Parallel()
	err := product.NewTennisBallProduct(db, s.FixtureStore()).Build(testName)
	s.NoError(err)
	sportCategory, err := s.GetFixture(category.SportCategoryReference, testName)
	s.NoError(err)
	berlinWarehouse, err := s.GetFixture(warehouse.BerlinWarehouseReference, testName)
	s.NoError(err)
	expectedResult := &model.Product{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Tennis ball",
		CategoryID:  sportCategory.GetID(),
		WarehouseID: berlinWarehouse.GetID(),
	}
	repo := &productRepo{db: db}
	actualResult, err := repo.GetOne(context.Background(), 1)
	s.NoError(err)
	s.Equal(expectedResult, actualResult)
}

func (s *ProductRepoTestSuite1) TestGetOneAgain_NoError() {
	db := s.DB()
	testName := s.T().Name()
	s.T().Parallel()
	err := product.NewTennisBallProduct(db, s.FixtureStore()).Build(testName)
	s.NoError(err)
	sportCategory, err := s.GetFixture(category.SportCategoryReference, testName)
	s.NoError(err)
	berlinWarehouse, err := s.GetFixture(warehouse.BerlinWarehouseReference, testName)
	s.NoError(err)
	expectedResult := &model.Product{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Tennis ball",
		CategoryID:  sportCategory.GetID(),
		WarehouseID: berlinWarehouse.GetID(),
	}
	repo := &productRepo{db: db}
	actualResult, err := repo.GetOne(context.Background(), 1)
	s.NoError(err)
	s.Equal(expectedResult, actualResult)
}

func (s *ProductRepoTestSuite1) TestGetOneIphone_NoError() {
	db := s.DB()
	testName := s.T().Name()
	s.T().Parallel()
	err := product.NewIphoneProduct(db, s.FixtureStore()).Build(testName)
	s.NoError(err)
	sportCategory, err := s.GetFixture(category.HiTechCategoryReference, testName)
	s.NoError(err)
	berlinWarehouse, err := s.GetFixture(warehouse.BerlinWarehouseReference, testName)
	s.NoError(err)
	expectedResult := &model.Product{
		BaseModel:   model.BaseModel{ID: 1},
		Name:        "Iphone 10",
		CategoryID:  sportCategory.GetID(),
		WarehouseID: berlinWarehouse.GetID(),
	}
	repo := &productRepo{db: db}
	actualResult, err := repo.GetOne(context.Background(), 1)
	s.NoError(err)
	s.Equal(expectedResult, actualResult)
}

func (s *ProductRepoTestSuite1) TestGetOneInParallel_NoError() {
	db := s.DB()
	s.T().Parallel()
	tests := []struct {
		name    string
		fixture func(testID string, db *sql.DB, foxStore fixture.FoxStore) *model.Product
		want    func(testID string) *model.Product
	}{
		{
			name: "test get one iphone 10",
			fixture: func(testID string, db *sql.DB, foxStore fixture.FoxStore) *model.Product {
				err := product.NewIphoneProduct(db, foxStore).Build(testID)
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
			fixture: func(testID string, db *sql.DB, foxStore fixture.FoxStore) *model.Product {
				err := product.NewTennisBallProduct(db, foxStore).Build(testID)
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

	var wg sync.WaitGroup
	for id := range tests {
		wg.Add(1)
		test := tests[id]
		s.T().Run(test.name, func(t *testing.T) {
			t.Parallel()
			repo := &productRepo{db: db}
			product := test.fixture(test.name, db, s.FixtureStore())
			actualResult, err := repo.GetOne(context.Background(), product.ID)
			s.NoError(err)
			s.Equal(test.want(test.name), actualResult)
			wg.Done()
		})
	}
	wg.Wait()
}
