package scope

import (
	"fmt"

	"github.com/gobuffalo/nulls"

	"github.com/gofrs/uuid"
	"gorm.io/gorm"
)

var stubRegex = "^SELECT\\s+FROM (\"|'|`)?stubs(\"|'|`)?( AS (\"|'|`)stubs(\"|'|`))?\\s+"
var stubRegexAlt = "^SELECT\\s{1,}FROM (\"|'|`)?stubs(\"|'|`)?( AS (\"|'|`)stubs(\"|'|`))?\\s{1,}"

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

func ForOne() ScopeFunc {
	return ForLimit(1)
}

// ForFirst scopes a query for the first record in a table based on the `created_at` timestamp.
//
// The first record is the oldest record.
func ForFirst() ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Order("created_at ASC, id ASC").Limit(1)
	}
}

// ForLast scopes a query for the last record in a table based on the `created_at` timestamp.
//
// The last record is the newest record.
func ForLast() ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Order("created_at DESC, id DESC").Limit(1)
	}
}

func ForLimit(limit int) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Limit(limit)
	}
}

func ForID(id string) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		if len(id) > 0 {
			return q.Where("id = ?", uuid.Must(uuid.FromString(id)))
		}

		return q.Where("1 = 0")
	}
}

func ForIDWithTableName(id string, tablename string) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		if len(id) > 0 {
			return q.Where(fmt.Sprintf("%s.id = ?", tablename), uuid.Must(uuid.FromString(id)))
		}

		return q.Where("1 = 0")
	}
}

func ForIDForModel(id string, model interface{}) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		if len(id) > 0 {
			tableName := TableName(model)
			return q.Where(fmt.Sprintf("%s.id = ?", tableName), uuid.Must(uuid.FromString(id)))
		}

		return q.Where("1 = 0")
	}
}

func ForIDs(ids []uuid.UUID) ScopeFunc {
	return ForIDSet(NewIDSet(ids))
}

func ForIDSet(idSet IDSet) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		if len(idSet) > 0 {
			return q.Where("id in (?)", idSet.Keys())
		}

		return q.Where("1 = 0")
	}
}

func ForIDsWithTableName(ids []uuid.UUID, tablename string) ScopeFunc {
	return ForIDSetWithTableName(NewIDSet(ids), tablename)
}

func ForIDSetWithTableName(idSet IDSet, tablename string) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		if len(idSet) > 0 {
			return q.Where(fmt.Sprintf("%s.id in (?)", tablename), idSet.Keys())
		}

		return q.Where("1 = 0")
	}
}

func ForIDsForModel(ids []uuid.UUID, model interface{}) ScopeFunc {
	return ForIDSetForModel(NewIDSet(ids), model)
}

func ForIDSetForModel(idSet IDSet, model interface{}) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		if len(idSet) > 0 {
			tableName := TableName(model)
			return q.Where(fmt.Sprintf("%s.id in (?)", tableName), idSet.Keys())
		}

		return q.Where("1 = 0")
	}
}

func ForNotID(id string) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		if len(id) > 0 {
			return q.Where("id != ?", uuid.Must(uuid.FromString(id)))
		}

		return q.Where("1 = 0")
	}
}

func ForNotIDWithTableName(id string, tablename string) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		if len(id) > 0 {
			return q.Where(fmt.Sprintf("%s.id != ?", tablename), uuid.Must(uuid.FromString(id)))
		}

		return q.Where("1 = 0")
	}
}

func ForNotIDForNotModel(id string, model interface{}) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		if len(id) > 0 {
			tableName := TableName(model)
			return q.Where(fmt.Sprintf("%s.id != ?", tableName), uuid.Must(uuid.FromString(id)))
		}

		return q.Where("1 = 0")
	}
}

func ForNotIDs(ids []uuid.UUID) ScopeFunc {
	return ForNotIDSet(NewIDSet(ids))
}

func ForNotIDSet(idSet IDSet) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		if len(idSet) > 0 {
			return q.Where("id not in (?)", idSet.Keys())
		}

		return q.Where("1 = 0")
	}
}

func ForNotIDsWithTableName(ids []uuid.UUID, tablename string) ScopeFunc {
	return ForNotIDSetWithTableName(NewIDSet(ids), tablename)
}

func ForNotIDSetWithTableName(idSet IDSet, tablename string) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		if len(idSet) > 0 {
			return q.Where(fmt.Sprintf("%s.id not in (?)", tablename), idSet.Keys())
		}

		return q.Where("1 = 0")
	}
}

func ForNotIDsForModel(ids []uuid.UUID, model interface{}) ScopeFunc {
	return ForNotIDSetForModel(NewIDSet(ids), model)
}

func ForNotIDSetForModel(idSet IDSet, model interface{}) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		if len(idSet) > 0 {
			tableName := TableName(model)
			return q.Where(fmt.Sprintf("%s.id not in (?)", tableName), idSet.Keys())
		}

		return q.Where("1 = 0")
	}
}

func ForUuidID(uid uuid.UUID) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where("id = ?", uid)
	}
}

func ForNotUuidID(uid uuid.UUID) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where("id != ?", uid)
	}
}

func ForNullsUuidID(uid nulls.UUID) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		if uid.Valid {
			return q.Where("id = ?", uid.UUID)
		}

		return q.Where("1 = 0")
	}
}

func ForNullDeletedAt() ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where("deleted_at is null")
	}
}

func ForNullDeletedAtForModel(model interface{}) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		tableName := TableName(model)
		return q.Where(fmt.Sprintf("%s.deleted_at is null", tableName))
	}
}

func ForNotNullDeletedAt() ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where("deleted_at is not null")
	}
}

func ForNotNullDeletedAtForModel(model interface{}) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		tableName := TableName(model)
		return q.Where(fmt.Sprintf("%s.deleted_at is not null", tableName))
	}
}

// ForNull scopes for a supplied field being null.
func ForNull(field string) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where(fmt.Sprintf("%s is null", field))
	}
}

// ForNotNull scopes for a supplied field being not null.
func ForNotNull(field string) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		return q.Where(fmt.Sprintf("%s is not null", field))
	}
}

func PaginateFromParams(params PaginationParams) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		paginator := NewPaginatorFromParams(params)
		page := paginator.Page
		if page == 0 {
			page = 1
		}

		pageSize := paginator.PerPage
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return q.Offset(offset).Limit(pageSize)
	}
}

func Paginate(page int, perPage int) ScopeFunc {
	return func(q *gorm.DB) *gorm.DB {
		paginator := NewPaginator(page, perPage)
		page := paginator.Page
		if page == 0 {
			page = 1
		}

		pageSize := paginator.PerPage
		switch {
		case pageSize > 100:
			pageSize = 100
		case pageSize <= 0:
			pageSize = 10
		}

		offset := (page - 1) * pageSize
		return q.Offset(offset).Limit(pageSize)
	}
}
