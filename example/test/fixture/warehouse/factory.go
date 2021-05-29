package warehouse

import (
	"database/sql"

	"github.com/tungnt/goit/example/model"
	"github.com/tungnt/goit/fixture"
)

type warehouseFixtureFactory struct {
	db *sql.DB
}

func newWarehouseFixtureFactory(db *sql.DB) *warehouseFixtureFactory {
	return &warehouseFixtureFactory{db: db}
}

func (f *warehouseFixtureFactory) createWarehouse(name string, city fixture.ModelWithID) (*model.Warehouse, error) {
	warehouse := model.Warehouse{
		Name:   name,
		CityID: city.GetID(),
	}
	stmt, err := f.db.Prepare("INSERT INTO warehouse(name, city_id) values(?, ?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(warehouse.Name, warehouse.CityID)
	if err != nil {
		return nil, err
	}
	warehouse.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &warehouse, nil
}
