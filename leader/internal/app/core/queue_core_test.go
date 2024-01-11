package core_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type CoreTestSuite struct {
	suite.Suite
}

func TestCoreTestSuite(t *testing.T) {
	suite.Run(t, new(CoreTestSuite))
}
