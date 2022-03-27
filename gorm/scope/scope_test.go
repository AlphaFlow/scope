package scope_test

import (
	"testing"

	"github.com/alphaflow/scope/gorm/suite"
)

type ScopesSuite struct {
	*suite.Model
}

func Test_ScopesSuite(t *testing.T) {
	model := suite.NewModel()

	ss := &ScopesSuite{
		Model: model,
	}

	suite.Run(t, ss)
}
