package product

import (
	"database/sql"

	"github.com/tungnt/goit/example/model"
)

type productFixtureFactory struct {
	db *sql.DB
}

func newProductFixtureFactory(db *sql.DB) *productFixtureFactory {
	return &productFixtureFactory{db: db}
}

func (f *productFixtureFactory) createProduct(name string, category, warehouse interface{}) (*model.Product, error) {
	product := model.Product{
		Name:        name,
		CategoryID:  category.(*model.Category).ID,
		WarehouseID: warehouse.(*model.Warehouse).ID,
	}
	stmt, err := f.db.Prepare("INSERT INTO product(name, category_id, warehouse_id) values(?, ?, ?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(product.Name, product.CategoryID, product.WarehouseID)
	if err != nil {
		return nil, err
	}
	product.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &product, nil
}
