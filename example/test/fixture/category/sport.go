package category

import (
	"database/sql"

	"github.com/tungnt/goit/fixture"
)

const (
	SportCategoryReference = "sport-category"
)

type sportCategory struct {
	*fixture.BaseFixture
}

func NewSportCategory(db *sql.DB) *sportCategory {
	categoryFixture := &sportCategory{BaseFixture: fixture.NewBaseFixture(db)}
	categoryFixture.SetFoxFixture(categoryFixture)
	categoryFixture.SetReference(SportCategoryReference)

	return categoryFixture
}

func (s *sportCategory) Create() (interface{}, error) {
	return newCategoryFixtureFactory(s.DB()).createCategory("Tennis")
}
