package scope_test

import (
	"github.com/alphaflow/scope"
	"github.com/gofrs/uuid"
)

func (ss *ScopesSuite) TestNewIdSet() {
	a := uuid.Must(uuid.NewV4())
	b := uuid.Must(uuid.NewV4())
	c := uuid.Must(uuid.NewV4())

	ids := []uuid.UUID{
		a, a, b, c,
	}

	idSet := scope.NewIDSet(ids)

	ss.Equal(scope.IDSet{
		a: true,
		b: true,
		c: true,
	}, idSet)
}

func (ss *ScopesSuite) TestNewIdSet_Keys() {
	a := uuid.Must(uuid.NewV4())
	b := uuid.Must(uuid.NewV4())
	c := uuid.Must(uuid.NewV4())

	ids := []uuid.UUID{
		a, a, b, c,
	}

	idSet := scope.NewIDSet(ids)

	keys := idSet.Keys()

	ss.ElementsMatch([]uuid.UUID{a, b, c}, keys)
}
