package category

import (
	"database/sql"

	"github.com/tungnt/goit/example/model"
)

type categoryFixtureFactory struct {
	db *sql.DB
}

func newCategoryFixtureFactory(db *sql.DB) *categoryFixtureFactory {
	return &categoryFixtureFactory{db: db}
}

func (f *categoryFixtureFactory) createCategory(name string) (*model.Category, error) {
	category := model.Category{Name: name}
	stmt, err := f.db.Prepare("INSERT INTO category(name) values(?)")
	if err != nil {
		return nil, err
	}
	res, err := stmt.Exec(category.Name)
	if err != nil {
		return nil, err
	}
	category.ID, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &category, nil
}
