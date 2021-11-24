package examples

import (
	"net/http"

	"github.com/alphaflow/scope"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/pop/v5"
)

// Placeholders for example.
var tx *pop.Connection
var r render.Engine
type ToDo struct {}  // Placeholder "Model"
type ToDos ToDo
type toDosResource struct{}

func registerToDosHandlers(app *buffalo.App) {
	tdr := &toDosResource{}

	// Example endpoints
	app.GET("/todos/filter_options", tdr.FilterOptions)
	app.GET("/todos/filter_columns", tdr.FilterColumns)
	app.GET("/todos/sort_columns", tdr.SortColumns)
	app.GET("/todos/aggregate", tdr.Aggregate)
	app.GET("/todos/grouped_aggregate", tdr.GroupedAggregate)
	app.GET("/todos", tdr.List)
}

func (tdr toDosResource) List(c buffalo.Context) error {
	toDos := &ToDos{}

	// Generic filtering
	filterScope, err := scope.ForFiltersFromParams(c, ToDo{}, c.Params())
	if err != nil {
		return c.Render(http.StatusBadRequest, r.Auto(c, err))
	}

	// Generic ordering
	orderScope, err := scope.ForSortFromParams(c, ToDo{}, c.Params())
	if err != nil {
		return c.Render(http.StatusBadRequest, r.Auto(c, err))
	}

	// Pagination
	paginateScope := scope.ForPaginateFromParams(c.Params())

	sc := scope.NewCollection(tx)
	sc.Push(filterScope, orderScope, paginateScope)
	if err := tx.Scope(sc.Flatten()).All(toDos); err != nil {
		return c.Render(http.StatusBadRequest, r.Auto(c, err))
	}

	return c.Render(http.StatusOK, r.Auto(c, toDos))
}

// Aggregate gets aggregate statistics
func (tdr toDosResource) Aggregate(c buffalo.Context) error {
	// Generic filtering in conjunction with aggregation is very powerful.
	filterScope, err := scope.ForFiltersFromParams(c, ToDo{}, c.Params())
	if err != nil {
		return c.Render(http.StatusBadRequest, r.Auto(c, err))
	}

	sc := scope.NewCollection(tx)
	sc.Push(filterScope)
	aggregate, err := scope.GetAggregationsFromParams(c, tx, &ToDos{}, c.Params(), sc)
	if err != nil {
		return c.Render(http.StatusBadRequest, r.Auto(c, err))
	}

	return c.Render(http.StatusOK, r.Auto(c, aggregate))
}

// GroupedAggregate gets grouped aggregate statistics
func (tdr toDosResource) GroupedAggregate(c buffalo.Context) error {
	// Generic filtering in conjunction with aggregation is very powerful.
	filterScope, err := scope.ForFiltersFromParams(c, ToDo{}, c.Params())
	if err != nil {
		return c.Render(http.StatusBadRequest, r.Auto(c, err))
	}

	sc := scope.NewCollection(tx)
	sc.Push(filterScope)

	groupedAggregates, err := scope.GetGroupedAggregationsFromParams(c, tx, &ToDos{}, c.Params(), sc)
	if err != nil {
		return c.Render(http.StatusBadRequest, r.Auto(c, err))
	}

	return c.Render(http.StatusOK, r.Auto(c, groupedAggregates))
}

// FilterOptions gets all values in the system for the supplied filter_column.
func (tdr toDosResource) FilterOptions(c buffalo.Context) error {
	sc := scope.NewCollection(tx)
	filterOptions, err := scope.GetFilterOptions(c, tx, &ToDos{}, c.Param("filter_column"), sc)
	if err != nil {
		return c.Render(http.StatusBadRequest, r.Auto(c, err))
	}

	return c.Render(http.StatusOK, r.Auto(c, filterOptions))
}

// FilterColumns gets all filterable columns.
func (tdr toDosResource) FilterColumns(c buffalo.Context) error {
	filterColumns, err := scope.GetAllFilterColumnNames(c, &ToDo{})
	if err != nil {
		return c.Render(http.StatusBadRequest, r.Auto(c, err))
	}

	return c.Render(http.StatusOK, r.Auto(c, filterColumns))
}

// SortColumns gets all sortable columns.
func (tdr toDosResource) SortColumns(c buffalo.Context) error {
	sortColumns, err := scope.GetAllSortColumnNames(c, &ToDo{})
	if err != nil {
		return c.Render(http.StatusBadRequest, r.Auto(c, err))
	}

	return c.Render(http.StatusOK, r.Auto(c, sortColumns))
}
