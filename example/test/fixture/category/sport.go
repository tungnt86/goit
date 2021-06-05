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

func NewSportCategory(db *sql.DB, foxStore fixture.FoxStore) *sportCategory {
	categoryFixture := &sportCategory{BaseFixture: fixture.NewBaseFixture(db, foxStore)}
	categoryFixture.SetFoxFixture(categoryFixture)
	categoryFixture.SetReference(SportCategoryReference)

	return categoryFixture
}

func (s *sportCategory) Create() (fixture.ModelWithID, error) {
	return newCategoryFixtureFactory(s.DB()).createCategory("Tennis")
}
