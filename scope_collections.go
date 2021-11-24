package scope

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"

	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

type Collection struct {
	tx     *pop.Connection
	scopes []pop.ScopeFunc
}

func NewCollection(tx ...*pop.Connection) *Collection {
	sc := &Collection{}
	sc.scopes = make([]pop.ScopeFunc, 0)

	if len(tx) > 0 && tx[0] != nil {
		sc.tx = tx[0]
	}

	return sc
}

func (sc *Collection) Get() []pop.ScopeFunc {
	return sc.scopes
}

func (sc *Collection) Flatten() pop.ScopeFunc {
	scopeFunc := func(q *pop.Query) *pop.Query {
		for _, sf := range sc.Dedupe().Get() {
			q = sf(q)
		}

		return q
	}

	return scopeFunc
}

func (sc *Collection) Dedupe() *Collection {
	if sc.tx == nil {
		panic(errors.Errorf("invalid tx value for Dedupe"))
	}

	sc.scopes = dedupeScopes(sc.tx, sc.Get()...)

	return sc
}

func (sc *Collection) Push(scopes ...pop.ScopeFunc) *Collection {
	sc.scopes = append(sc.scopes, scopes...)

	return sc
}

func dedupeScopes(tx *pop.Connection, scopes ...pop.ScopeFunc) []pop.ScopeFunc {
	type __stub__ struct{}
	dedupedScopeMap := make(map[string]pop.ScopeFunc)
	for _, s := range scopes {
		scopeQueryFunc := s(tx.Q())
		scopeQuerySQL, _ := scopeQueryFunc.ToSQL(&pop.Model{Value: __stub__{}})
		regex := regexp.MustCompile(`^SELECT\s{1,}FROM stubs AS stubs\s{1,}`)
		scopeQueryRaw := regex.ReplaceAllString(scopeQuerySQL, "")

		hash := md5HashForString(scopeQueryRaw)
		if _, ok := dedupedScopeMap[hash]; !ok {
			dedupedScopeMap[hash] = s
		}
	}
	sparseScopes := make([]pop.ScopeFunc, 0, len(dedupedScopeMap))
	for _, s := range dedupedScopeMap {
		sparseScopes = append(sparseScopes, s)
	}

	return sparseScopes
}

func md5HashForString(s string) string {
	hash := md5.New()
	hash.Write([]byte(s))

	return hex.EncodeToString(hash.Sum(nil))
}
