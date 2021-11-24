package scope

import (
	"fmt"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
)

type IDSet map[uuid.UUID]bool

func NewIDSet(ids []uuid.UUID) IDSet {
	idSet := IDSet{}
	for i := range ids {
		id := ids[i]
		idSet[id] = true
	}

	return idSet
}

func (idSet IDSet) Keys() []uuid.UUID {
	i := 0
	ids := make([]uuid.UUID, len(idSet))
	for id := range idSet {
		ids[i] = id
		i++
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
	return ForIDSet(NewIDSet(ids))
}

func ForIDSet(idSet IDSet) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(idSet) > 0 {
			return q.Where("id in (?)", idSet.Keys())
		}

		return q.Where("1 = 0")
	}
}

func ForIDsWithTableName(ids []uuid.UUID, tablename string) pop.ScopeFunc {
	return ForIDSetWithTableName(NewIDSet(ids), tablename)
}

func ForIDSetWithTableName(idSet IDSet, tablename string) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(idSet) > 0 {
			return q.Where(fmt.Sprintf("%s.id in (?)", tablename), idSet.Keys())
		}

		return q.Where("1 = 0")
	}
}

func ForIDsForModel(ids []uuid.UUID, model interface{}) pop.ScopeFunc {
	return ForIDSetForModel(NewIDSet(ids), model)
}

func ForIDSetForModel(idSet IDSet, model interface{}) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(idSet) > 0 {
			tableNameAble := pop.Model{Value: model}
			return q.Where(fmt.Sprintf("%s.id in (?)", tableNameAble.TableName()), idSet.Keys())
		}

		return q.Where("1 = 0")
	}
}

func ForNotID(id string) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(id) > 0 {
			return q.Where("id != ?", uuid.Must(uuid.FromString(id)))
		}

		return q.Where("1 = 0")
	}
}

func ForNotIDWithTableName(id string, tablename string) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(id) > 0 {
			return q.Where(fmt.Sprintf("%s.id != ?", tablename), uuid.Must(uuid.FromString(id)))
		}

		return q.Where("1 = 0")
	}
}

func ForNotIDForNotModel(id string, model interface{}) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(id) > 0 {
			tableNameAble := pop.Model{Value: model}
			return q.Where(fmt.Sprintf("%s.id != ?", tableNameAble.TableName()), uuid.Must(uuid.FromString(id)))
		}

		return q.Where("1 = 0")
	}
}

func ForNotIDs(ids []uuid.UUID) pop.ScopeFunc {
	return ForNotIDSet(NewIDSet(ids))
}

func ForNotIDSet(idSet IDSet) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(idSet) > 0 {
			return q.Where("id not in (?)", idSet.Keys())
		}

		return q.Where("1 = 0")
	}
}

func ForNotIDsWithTableName(ids []uuid.UUID, tablename string) pop.ScopeFunc {
	return ForNotIDSetWithTableName(NewIDSet(ids), tablename)
}

func ForNotIDSetWithTableName(idSet IDSet, tablename string) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(idSet) > 0 {
			return q.Where(fmt.Sprintf("%s.id not in (?)", tablename), idSet.Keys())
		}

		return q.Where("1 = 0")
	}
}

func ForNotIDsForModel(ids []uuid.UUID, model interface{}) pop.ScopeFunc {
	return ForNotIDSetForModel(NewIDSet(ids), model)
}

func ForNotIDSetForModel(idSet IDSet, model interface{}) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		if len(idSet) > 0 {
			tableNameAble := pop.Model{Value: model}
			return q.Where(fmt.Sprintf("%s.id not in (?)", tableNameAble.TableName()), idSet.Keys())
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
