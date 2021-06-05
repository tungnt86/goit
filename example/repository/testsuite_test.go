package repository

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestSuite(t *testing.T) {
	suite.Run(t, new(ProductRepoTestSuite1))
	suite.Run(t, new(ProductRepoTestSuite2))
}
