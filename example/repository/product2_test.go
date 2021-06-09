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

type ProductRepoTestSuite2 struct {
	goit.ITsqlite
}

func (s *ProductRepoTestSuite2) TestGetOne_NoError() {
	s.T().Parallel()
	db, err := s.GetCurrentTestDB()
	s.NoError(err)
	foxStore := s.NewFixtureStore()
	err = product.NewIphoneProduct(db, foxStore).Build()
	s.NoError(err)
	sportCategory, err := s.GetFixture(foxStore, category.HiTechCategoryReference)
	s.NoError(err)
	berlinWarehouse, err := s.GetFixture(foxStore, warehouse.BerlinWarehouseReference)
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

func (s *ProductRepoTestSuite2) TestGetOneInParallel_NoError() {
	s.T().Parallel()
	tests := []struct {
		name    string
		fixture func(db *sql.DB, foxStore fixture.FoxStore) *model.Product
		want    func(foxStore fixture.FoxStore) *model.Product
	}{
		{
			name: "test get one iphone 10",
			fixture: func(db *sql.DB, foxStore fixture.FoxStore) *model.Product {
				err := product.NewIphoneProduct(db, foxStore).Build()
				s.NoError(err)
				iphoneProduct, err := s.GetFixture(foxStore, product.IphoneProductReference)
				s.NoError(err)
				return iphoneProduct.(*model.Product)
			},
			want: func(foxStore fixture.FoxStore) *model.Product {
				iphoneProduct, err := s.GetFixture(foxStore, product.IphoneProductReference)
				s.NoError(err)
				hitechCategory, err := s.GetFixture(foxStore, category.HiTechCategoryReference)
				s.NoError(err)
				berlinWarehouse, err := s.GetFixture(foxStore, warehouse.BerlinWarehouseReference)
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
			fixture: func(db *sql.DB, foxStore fixture.FoxStore) *model.Product {
				err := product.NewTennisBallProduct(db, foxStore).Build()
				s.NoError(err)
				product, err := s.GetFixture(foxStore, product.TennisBallProductReference)
				s.NoError(err)
				return product.(*model.Product)
			},
			want: func(foxStore fixture.FoxStore) *model.Product {
				tennisProduct, err := s.GetFixture(foxStore, product.TennisBallProductReference)
				s.NoError(err)
				sportCategory, err := s.GetFixture(foxStore, category.SportCategoryReference)
				s.NoError(err)
				berlinWarehouse, err := s.GetFixture(foxStore, warehouse.BerlinWarehouseReference)
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

	db, err := s.GetCurrentTestDB()
	s.NoError(err)
	var wg sync.WaitGroup
	for id := range tests {
		wg.Add(1)
		test := tests[id]
		s.T().Run(test.name, func(t *testing.T) {
			t.Parallel()
			foxStore := s.NewFixtureStore()
			repo := &productRepo{db: db}
			product := test.fixture(db, foxStore)
			actualResult, err := repo.GetOne(context.Background(), product.ID)
			s.NoError(err)
			s.Equal(test.want(foxStore), actualResult)
			wg.Done()
		})
	}
	wg.Wait()
}
