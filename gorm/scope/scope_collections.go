package scope

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Collection struct {
	tx     *gorm.DB
	scopes []ScopeFunc
}

func NewCollection(tx ...*gorm.DB) *Collection {
	sc := &Collection{}
	sc.scopes = make([]ScopeFunc, 0)

	if len(tx) > 0 && tx[0] != nil {
		sc.tx = tx[0]
	}

	return sc
}

func (sc *Collection) Get() []ScopeFunc {
	return sc.scopes
}

func (sc *Collection) Flatten() ScopeFunc {
	scopeFunc := func(q *gorm.DB) *gorm.DB {
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

func (sc *Collection) Push(scopes ...ScopeFunc) *Collection {
	sc.scopes = append(sc.scopes, scopes...)

	return sc
}

func dedupeScopes(tx *gorm.DB, scopes ...ScopeFunc) []ScopeFunc {
	dedupedScopeMap := make(map[string]ScopeFunc)
	for _, s := range scopes {
		q := tx.Session(&gorm.Session{DryRun: true, NewDB: true, SkipHooks: true}).Model(__stub__{})
		q.Statement.SQL.Reset()
		scopeQueryFunc := s(q).Find(q.Statement.Model)
		scopeQuerySQL := scopeQueryFunc.Statement.SQL.String()
		scopeQuerySQL = strings.Replace(scopeQuerySQL, "SELECT *", "SELECT ", 1)
		regex := regexp.MustCompile(stubRegexAlt)
		scopeQueryRaw := regex.ReplaceAllString(scopeQuerySQL, "")

		hash := md5HashForString(scopeQueryRaw)
		if _, ok := dedupedScopeMap[hash]; !ok {
			dedupedScopeMap[hash] = s
		}
	}
	sparseScopes := make([]ScopeFunc, 0, len(dedupedScopeMap))
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
