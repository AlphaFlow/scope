package scope_test

import (
	"github.com/gofrs/uuid"

	"github.com/alphaflow/api-core/gorm/scope"
)

func (ss *ScopesSuite) TestScopeCollection_dedupe() {
	collection := scope.NewCollection(ss.DB)

	scope1 := scope.ForID(uuid.Nil.String())
	scope2 := scope.ForID(uuid.Nil.String())
	scope3 := scope.ForOrder("id ASC")
	scope4 := scope.ForOrder("id ASC")
	collection.Push(scope1, scope2, scope3, scope4)
	deduped := collection.Dedupe()
	normalized_scopes := deduped.Get()
	ss.Equal(2, len(normalized_scopes))
}
