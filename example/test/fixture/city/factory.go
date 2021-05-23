package city

import (
	"database/sql"

	"github.com/tungnt/goit/example/model"
)

type cityFixtureFactory struct {
	db *sql.DB
}

func newCityFixtureFactory(db *sql.DB) *cityFixtureFactory {
	return &cityFixtureFactory{db: db}
}

func (f *cityFixtureFactory) createCity(name string) (*model.City, error) {
	city := model.City{Name: name}
	stmt, err := f.db.Prepare("INSERT INTO city(name) values(?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(city.Name)
	if err != nil {
		return nil, err
	}
	city.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &city, nil
}
