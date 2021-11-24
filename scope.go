package scope

import (
	"fmt"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
)

type IDSet map[uuid.UUID]bool

func (idSet IDSet) Keys() []uuid.UUID {
	ids := make([]uuid.UUID, len(idSet))
	for i := range idSet {
		ids = append(ids, i)
	}

	return ids
}

func ForOne() pop.ScopeFunc {
	return ForLimit(1)
}

// ForFirst scopes a query for the first record in a table based on the `created_at` timestamp.
//
// The first record is the oldest record.
func ForFirst() pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.Order("created_at ASC, id ASC").Limit(1)
	}
}

// ForLast scopes a query for the last record in a table based on the `created_at` timestamp.
//
// The last record is the newest record.
func ForLast() pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.Order("created_at DESC, id DESC").Limit(1)
	}
}

func ForLimit(limit int) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.Limit(limit)
	}
}

func ForID(id string) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(id) > 0 {
			return q.Where("id = ?", uuid.Must(uuid.FromString(id)))
		}

		return q.Where("1 = 0")
	}
}

func ForIDWithTableName(id string, tablename string) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(id) > 0 {
			return q.Where(fmt.Sprintf("%s.id = ?", tablename), uuid.Must(uuid.FromString(id)))
		}

		return q.Where("1 = 0")
	}
}

func ForIDForModel(id string, model interface{}) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(id) > 0 {
			tableNameAble := pop.Model{Value: model}
			return q.Where(fmt.Sprintf("%s.id = ?", tableNameAble.TableName()), uuid.Must(uuid.FromString(id)))
		}

		return q.Where("1 = 0")
	}
}

func ForIDs(ids []uuid.UUID) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(ids) > 0 {
			return q.Where("id in (?)", ids)
		}

		return q.Where("1 = 0")
	}
}

func ForIDSet(idSet IDSet) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(idSet) > 0 {
			return q.Where("id in (?)", idSet.Keys())
		}

		return q.Where("1 = 0")
	}
}

func ForUuidID(uid uuid.UUID) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.Where("id = ?", uid)
	}
}

func ForNotUuidID(uid uuid.UUID) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.Where("id != ?", uid)
	}
}

func ForNullsUuidID(uid nulls.UUID) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if uid.Valid {
			return q.Where("id = ?", uid.UUID)
		}

		return q.Where("1 = 0")
	}
}

func ForNullDeletedAt() pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.Where("deleted_at is null")
	}
}

func ForNullDeletedAtForModel(model interface{}) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		tableNameAble := pop.Model{Value: model}
		return q.Where(fmt.Sprintf("%s.deleted_at is null", tableNameAble.TableName()))
	}
}

func ForNotNullDeletedAt() pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.Where("deleted_at is not null")
	}
}
func ForNotNullDeletedAtForModel(model interface{}) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		tableNameAble := pop.Model{Value: model}
		return q.Where(fmt.Sprintf("%s.deleted_at is not null", tableNameAble.TableName()))
	}
}

// ForNull scopes for a supplied field being null.
func ForNull(field string) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.Where(fmt.Sprintf("%s is null", field))
	}
}

// ForNotNull scopes for a supplied field being not null.
func ForNotNull(field string) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.Where(fmt.Sprintf("%s is not null", field))
	}
}
