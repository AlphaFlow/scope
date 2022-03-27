package scope

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/alphaflow/scope/util"
)

const FailQuery = "1=0"

/*************** Generic Filtering and Sorting ****************/
// Filter operands that can be used within ForFiltersFromParams
var filterTypes = map[string]string{
	"EQ":   "=",
	"NE":   "!=",
	"LT":   "<",
	"GT":   ">",
	"LTE":  "<=",
	"GTE":  ">=",
	"NU":   "is null",
	"NN":   "is not null",
	"LK":   "like",
	"ILK":  "ilike",
	"NLK":  "not like",
	"NILK": "not ilike",
	"DF":   "is distinct from",
	"NDF":  "is not distinct from",
	"IN":   "in",
	"NIN":  "not in",
}

// Filter logics that can be used within ForFiltersFromParams
var filterLogics = map[string]string{
	"AND": "AND",
	"OR":  "OR",
}

// Sort directions that can be used within ForSortFromParams
var sortDirections = map[string]string{
	"ASC":  "ASC",
	"DESC": "DESC",
}

// ForFiltersFromParams filters a model based on the provided filter params.
func ForFiltersFromParams(ctx context.Context, model interface{}, params buffalo.ParamValues) (pop.ScopeFunc, error) {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Struct {
		return nil, errors.New("struct expected")
	}
	modelPtr := reflect.New(reflect.TypeOf(model)).Interface()

	filterSeparator := getFilterSeparator(params)
	filterArgsSeparator := getFilterArgsSeparator(params)

	columns := make([]string, 0)
	if !util.IsBlank(params.Get("filter_columns")) {
		columns = strings.Split(params.Get("filter_columns"), filterSeparator)
	}
	types := make([]string, 0)
	if !util.IsBlank(params.Get("filter_types")) {
		types = strings.Split(params.Get("filter_types"), filterSeparator)
	}
	leftParens := make([]string, 0)
	if !util.IsBlank(params.Get("filter_left_parens")) {
		leftParens = strings.Split(params.Get("filter_left_parens"), filterSeparator)
	}
	rightParens := make([]string, 0)
	if !util.IsBlank(params.Get("filter_right_parens")) {
		rightParens = strings.Split(params.Get("filter_right_parens"), filterSeparator)
	}
	logic := make([]string, 0)
	if !util.IsBlank(params.Get("filter_logic")) {
		logic = strings.Split(params.Get("filter_logic"), filterSeparator)
	}

	// filter_values can be empty and still be valid (for example, X = "")
	values := strings.Split(params.Get("filter_values"), filterSeparator)

	// If nothing is specified, this is a no-op.
	if len(columns) == 0 && len(types) == 0 && len(values) == 1 && len(logic) == 0 && len(leftParens) == 0 && len(rightParens) == 0 {
		return func(q *pop.Query) *pop.Query {
			return q
		}, nil
	}

	if len(columns) != len(types) || len(columns) != len(values) || len(columns) != len(logic)+1 || len(leftParens) != len(rightParens) {
		// We must have the same number of all filtering params.  We must have 1 more column than logical operators.
		return nil, errors.New("missing or mismatched filter parameters")
	}

	// Check for custom filter fields, and handle appropriately.
	columnMap := make(map[string]string, 0)
	filterColumns, err := GetAllFilterColumns(ctx, modelPtr)
	if err != nil {
		return nil, err
	}

	for _, column := range filterColumns {
		columnMap[column.Name] = column.Statement
	}

	clauses := make([]string, len(columns))
	args := make([]interface{}, 0)
	argsPerClause := make([]int, len(columns))
	for i, col := range columns {
		// Find the correct operator for this filter.
		op, ok := filterTypes[strings.ToUpper(types[i])]
		if !ok {
			return nil, errors.Errorf("invalid filter type: %v", types[i])
		}

		// If this column is filterable, build this clause.
		if stmt, ok := columnMap[col]; ok {
			clauses[i] = fmt.Sprintf("%s %s", stmt, op)
		} else {
			return nil, errors.Errorf("invalid filter field: %v", col)
		}

		// Add the arg to the list if this operator takes args.
		if filterOperatorHasOneArg(types[i]) {
			args = append(args, values[i])
			argsPerClause[i] = 1
		} else if filterOperatorHasArgs(types[i]) {
			if util.IsBlank(values[i]) {
				argsPerClause[i] = 0
			} else {

				separatedValues := strings.Split(values[i], filterArgsSeparator)
				for _, arg := range separatedValues {
					args = append(args, arg)
				}
				argsPerClause[i] = len(separatedValues)
			}
		}
	}

	// Convert parenthesis into appropriate strings.
	leftParenIndicies := make(map[int]string, 0)
	rightParenIndicies := make(map[int]string, 0)
	for _, i := range leftParens {
		index, err := strconv.Atoi(i)
		if index > len(columns)-1 || err != nil {
			return nil, errors.Errorf("invalid filter parentheses: %v", i)
		}
		leftParenIndicies[index] = leftParenIndicies[index] + "("
	}

	for _, i := range rightParens {
		index, err := strconv.Atoi(i)
		if index > len(columns)-1 || err != nil {
			return nil, errors.Errorf("invalid filter parentheses: %v", i)
		}
		rightParenIndicies[index] = rightParenIndicies[index] + ")"
	}

	// Apply Logic, starting with the first clause.
	queryString := buildFilterClause(clauses[0], types[0], argsPerClause[0], leftParenIndicies[0], rightParenIndicies[0])
	for i, l := range logic {
		logic, ok := filterLogics[strings.ToUpper(l)]
		if !ok {
			return nil, errors.Errorf("invalid filter logic: %v", logic[i])
		}

		clauseWithArgs := buildFilterClause(clauses[i+1], types[i+1], argsPerClause[i+1], leftParenIndicies[i+1], rightParenIndicies[i+1])
		queryString = fmt.Sprintf("%s %s %s", queryString, logic, clauseWithArgs)
	}

	// Wrap our query string in parens, so its always evaluated as 1 expression and cannot conflict with other scopes.
	queryString = fmt.Sprintf("(%s)", queryString)

	return func(q *pop.Query) *pop.Query {
		return q.Where(queryString, args...)
	}, nil
}

func buildFilterClause(clause, operator string, argsPerClause int, leftParen, rightParen string) string {
	if filterOperatorHasNoArgs(operator) {
		return fmt.Sprintf("%s%s%s", leftParen, clause, rightParen)
	} else if filterOperatorHasOneArg(operator) {
		return fmt.Sprintf("%s%s ?%s", leftParen, clause, rightParen)
	}

	if argsPerClause == 0 {
		return fmt.Sprintf("%s%s%s", leftParen, FailQuery, rightParen)
	}

	// We add an extra space to the last argument in the query, to circumvent https://github.com/gobuffalo/pop/issues/610
	// Luckily, the pop code only replaces the exact string "(?)", so by returning "( ?)" our args are supplied properly.
	return fmt.Sprintf("%s%s (%v ?)%s", leftParen, clause, strings.Repeat("?, ", argsPerClause-1), rightParen)
}

func filterOperatorHasNoArgs(operator string) bool {
	return strings.ToUpper(operator) == "NN" || strings.ToUpper(operator) == "NU"
}

func filterOperatorHasOneArg(operator string) bool {
	return !(filterOperatorHasNoArgs(operator) || filterOperatorHasArgs(operator))
}

func filterOperatorHasArgs(operator string) bool {
	return strings.ToUpper(operator) == "IN" || strings.ToUpper(operator) == "NIN"
}

// ForOrder is a generic scope function for ordering.
func ForOrder(orderClauses ...string) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		for _, clause := range orderClauses {
			q.Order(clause)
		}
		return q
	}
}

// ForSortFromParams orders a query based on the provided query params.
func ForSortFromParams(ctx context.Context, model interface{}, params buffalo.ParamValues) (pop.ScopeFunc, error) {
	v := reflect.ValueOf(model)
	if v.Kind() != reflect.Struct {
		return nil, errors.New("struct expected")
	}
	modelPtr := reflect.New(reflect.TypeOf(model)).Interface()

	filterSeparator := getFilterSeparator(params)

	columns := make([]string, 0)
	if !util.IsBlank(params.Get("sort_columns")) {
		columns = strings.Split(params.Get("sort_columns"), filterSeparator)
	}
	directions := make([]string, 0)
	if !util.IsBlank(params.Get("sort_directions")) {
		directions = strings.Split(params.Get("sort_directions"), filterSeparator)
	}

	// If nothing is specified, this is a no-op.
	if len(columns) == 0 && len(directions) == 0 {
		return func(q *pop.Query) *pop.Query {
			return q
		}, nil
	}

	if len(columns) != len(directions) {
		// We must have the same number of all sorting params.
		return nil, errors.New("missing or mismatched sort parameters")
	}

	// Check for custom sort fields, and handle appropriately.
	columnMap := make(map[string]string, 0)
	sortColumns, err := GetAllSortColumns(ctx, modelPtr)
	if err != nil {
		return nil, err
	}

	for _, column := range sortColumns {
		columnMap[column.Name] = column.Statement
	}

	clauses := make([]string, len(columns))
	for i, col := range columns {
		// Find the correct operator for this filter.
		op, ok := sortDirections[strings.ToUpper(directions[i])]
		if !ok {
			return nil, errors.New(fmt.Sprintf("invalid sort direction: %v", directions[i]))
		}

		// If this column is sortable, build this clause.
		if stmt, ok := columnMap[col]; ok {
			clauses[i] = fmt.Sprintf("%s %s", stmt, op)
		} else {
			return nil, errors.Errorf("invalid sort field: %v", col)
		}
	}

	return ForOrder(clauses...), nil
}

// ForPaginateFromParams paginates a query based on a list of parameters, generally c.Params()
func ForPaginateFromParams(params buffalo.ParamValues) pop.ScopeFunc {
	return func(q *pop.Query) *pop.Query {
		return q.PaginateFromParams(params)
	}
}

// getFilterSeparator gets the filterSeparator token. The parameter filter_separator can be used to separate the filter
// columns, etc, if | is not suitable.
func getFilterSeparator(params buffalo.ParamValues) string {
	filterSeparator := "|"
	if !util.IsBlank(params.Get("filter_separator")) {
		filterSeparator = params.Get("filter_separator")
	}

	return filterSeparator
}

// getFilterArgsSeparator gets the filterSeparator token. The parameter filter_args_separator can be used to separate
// the filter args for operators with many args, if , is not suitable.
func getFilterArgsSeparator(params buffalo.ParamValues) string {
	filterArgsSeparator := ","
	if !util.IsBlank(params.Get("filter_args_separator")) {
		filterArgsSeparator = params.Get("filter_args_separator")
	}

	return filterArgsSeparator
}
