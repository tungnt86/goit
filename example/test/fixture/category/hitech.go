package category

import (
	"database/sql"

	"github.com/tungnt/goit/fixture"
)

const (
	HiTechCategoryReference = "hitech-category"
)

type hiTechCategory struct {
	*fixture.BaseFixture
}

func NewHiTechCategory(db *sql.DB) *hiTechCategory {
	categoryFixture := &hiTechCategory{BaseFixture: fixture.NewBaseFixture(db)}
	categoryFixture.SetFoxFixture(categoryFixture)
	categoryFixture.SetReference(HiTechCategoryReference)

	return categoryFixture
}

func (s *hiTechCategory) Create() (interface{}, error) {
	return newCategoryFixtureFactory(s.DB()).createCategory("Smartphone")
}
