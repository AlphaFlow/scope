package scope

import (
	"reflect"
	"strconv"

	"github.com/gobuffalo/x/defaults"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ScopeFunc func(q *gorm.DB) *gorm.DB
type Tabler interface {
	TableName() string
}

type __stub__ struct{}

func (s __stub__) TableName() string {
	return "stubs"
}
func TableName(entity interface{}) string {
	if n, ok := entity.(Tabler); ok {
		return n.TableName()
	}

	t := reflect.TypeOf(entity)

	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.Slice, reflect.Array:
		el := t.Elem()
		if el.Kind() == reflect.Ptr {
			el = el.Elem()
		}

		// validates if the elem of slice or array implements TableNameAble interface.
		var tableNameAble *Tabler
		if el.Implements(reflect.TypeOf(tableNameAble).Elem()) {
			v := reflect.New(el)
			out := v.MethodByName("TableName").Call([]reflect.Value{})
			return out[0].String()
		}

		strategy := schema.NamingStrategy{}
		return strategy.TableName(el.Name())
	}

	strategy := schema.NamingStrategy{}
	return strategy.TableName(t.Name())
}

// Paginator is a type used to represent the pagination of records
// from the database.
type Paginator struct {
	// Current page you're on
	Page int `json:"page"`
	// Number of results you want per page
	PerPage int `json:"per_page"`
	// Page * PerPage (ex: 2 * 20, Offset == 40)
	Offset int `json:"offset"`
	// Total potential records matching the query
	TotalEntriesSize int `json:"total_entries_size"`
	// Total records returns, will be <= PerPage
	CurrentEntriesSize int `json:"current_entries_size"`
	// Total pages
	TotalPages int `json:"total_pages"`
}

// PaginatorPerPageDefault is the amount of results per page
var PaginatorPerPageDefault = 20

// PaginatorPageKey is the query parameter holding the current page index
var PaginatorPageKey = "page"

// PaginatorPerPageKey is the query parameter holding the amount of results per page
// to override the default one
var PaginatorPerPageKey = "per_page"

// PaginationParams is a parameters provider interface to get the pagination params from
type PaginationParams interface {
	Get(key string) string
}

// NewPaginator returns a new `Paginator` value with the appropriate
// defaults set.
func NewPaginator(page int, perPage int) *Paginator {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	p := &Paginator{Page: page, PerPage: perPage}
	p.Offset = (page - 1) * p.PerPage
	return p
}

func NewPaginatorFromParams(params PaginationParams) *Paginator {
	page := defaults.String(params.Get(PaginatorPageKey), "1")

	perPage := defaults.String(params.Get(PaginatorPerPageKey), strconv.Itoa(PaginatorPerPageDefault))

	p, err := strconv.Atoi(page)
	if err != nil {
		p = 1
	}

	pp, err := strconv.Atoi(perPage)
	if err != nil {
		pp = PaginatorPerPageDefault
	}
	return NewPaginator(p, pp)
}
