package scope_test

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"testing"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/alphaflow/scope/gorm/scope"
)

type TestModel struct {
	ID        uuid.UUID  `json:"id" db:"id" gorm:"primaryKey;column:id"`
	Nuid      nulls.UUID `json:"null_id" db:"db_null_id" gorm:"column:db_null_id"`
	NotInDb   int        `json:"not_in_db" db:"-" gorm:"-"`
	NotInJson int        `json:"-" db:"not_in_json" gorm:"column:not_in_json"`
}

func (t TestModel) GetCustomFilters(ctx context.Context) scope.CustomColumns {
	customFilter := scope.CustomColumn{
		Name:       "custom_filter",
		Statement:  `(SELECT '1234')`,
		ResultType: reflect.TypeOf("1234"),
	}
	return scope.CustomColumns{customFilter}
}

func (t TestModel) GetCustomSorts(ctx context.Context) scope.CustomColumns {
	customSort := scope.CustomColumn{
		Name:      "custom_sort",
		Statement: `(SELECT '1234')`,
	}
	return scope.CustomColumns{customSort}
}

func (ss *ScopesSuite) TestForFiltersFromParams() {
	tm := TestModel{}
	q := ss.DB.Session(&gorm.Session{DryRun: true}).Model(tm)
	q.Statement.SQL.Reset()
	scopeQueryFunc := q.Find(&tm)
	baseQuery := scopeQueryFunc.Statement.SQL.String()

	testCases := []struct {
		Name          string
		Params        map[string][]string
		ExpectErr     bool
		ExpectedQuery string
		ExpectedArgs  []string
	}{
		{
			Name: "No Params",
			Params: map[string][]string{
				"filter_columns": {""},
				"filter_types":   {""},
				"filter_values":  {""},
			},
			ExpectErr:     false,
			ExpectedQuery: baseQuery,
			ExpectedArgs:  []string{},
		},
		{
			Name: "Missing Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {""},
				"filter_values":  {"test"},
			},
			ExpectErr: true,
		},
		{
			Name: "Excess Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {""},
				"filter_values":  {"test"},
				"filter_logic":   {"and"},
			},
			ExpectErr: true,
		},
		{
			Name: "Column Is Dash - JSON",
			Params: map[string][]string{
				"filter_columns": {"-"},
				"filter_types":   {"eq"},
				"filter_values":  {"test"},
			},
			ExpectErr: true,
		},
		{
			Name: "Column Is Dash - DB",
			Params: map[string][]string{
				"filter_columns": {"not_in_db"},
				"filter_types":   {"eq"},
				"filter_values":  {"test"},
			},
			ExpectErr: true,
		},
		{
			Name: "Invalid parens",
			Params: map[string][]string{
				"filter_columns":      {"id"},
				"filter_types":        {"eq"},
				"filter_values":       {"test"},
				"filter_left_parens":  {"1"},
				"filter_right_parens": {"1"},
			},
			ExpectErr: true,
		},
		{
			Name: "Imbalanced parens",
			Params: map[string][]string{
				"filter_columns":     {"id"},
				"filter_types":       {"eq"},
				"filter_values":      {"test"},
				"filter_left_parens": {"0"},
			},
			ExpectErr: true,
		},
		{
			Name: "Equal Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"eq"},
				"filter_values":  {"test"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id = $1)", baseQuery),
			ExpectedArgs:  []string{"test"},
		},
		{
			Name: "Equal Operator Parens",
			Params: map[string][]string{
				"filter_columns":      {"id"},
				"filter_types":        {"eq"},
				"filter_values":       {"test"},
				"filter_left_parens":  {"0"},
				"filter_right_parens": {"0"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE ((test_models.id = $1))", baseQuery),
			ExpectedArgs:  []string{"test"},
		},
		{
			Name: "Equal Operator Custom Column",
			Params: map[string][]string{
				"filter_columns": {"custom_filter"},
				"filter_types":   {"eq"},
				"filter_values":  {"test"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE ((SELECT '1234') = $1)", baseQuery),
			ExpectedArgs:  []string{"test"},
		},
		{
			Name: "Not Equal Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"ne"},
				"filter_values":  {"test"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id != $1)", baseQuery),
			ExpectedArgs:  []string{"test"},
		},
		{
			Name: "Less Than Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"lt"},
				"filter_values":  {"test"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id < $1)", baseQuery),
			ExpectedArgs:  []string{"test"},
		},
		{
			Name: "Greater Than Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"gt"},
				"filter_values":  {"test"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id > $1)", baseQuery),
			ExpectedArgs:  []string{"test"},
		},
		{
			Name: "Is Distinct From Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"df"},
				"filter_values":  {"test"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id is distinct from $1)", baseQuery),
			ExpectedArgs:  []string{"test"},
		},
		{
			Name: "Is Not Distinct From Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"ndf"},
				"filter_values":  {"test"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id is not distinct from $1)", baseQuery),
			ExpectedArgs:  []string{"test"},
		},
		{
			Name: "Is Null Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"nu"},
				"filter_values":  {"test"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id is null)", baseQuery),
			ExpectedArgs:  []string{},
		},
		{
			Name: "Is Null Operator Parens",
			Params: map[string][]string{
				"filter_columns":      {"id"},
				"filter_types":        {"nu"},
				"filter_values":       {"test"},
				"filter_left_parens":  {"0"},
				"filter_right_parens": {"0"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE ((test_models.id is null))", baseQuery),
			ExpectedArgs:  []string{},
		},
		{
			Name: "Not Null Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"nn"},
				"filter_values":  {"test"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id is not null)", baseQuery),
			ExpectedArgs:  []string{},
		},
		{
			Name: "Like Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"lk"},
				"filter_values":  {"%test%"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id like $1)", baseQuery),
			ExpectedArgs:  []string{"%test%"},
		},
		{
			Name: "Ilike Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"ilk"},
				"filter_values":  {"%test%"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id ilike $1)", baseQuery),
			ExpectedArgs:  []string{"%test%"},
		},
		{
			Name: "In Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"in"},
				"filter_values":  {"test,test,2,3"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id in ($1, $2, $3,  $4))", baseQuery),
			ExpectedArgs:  []string{"test", "test", "2", "3"},
		},
		{
			Name: "In Operator Parens",
			Params: map[string][]string{
				"filter_columns":      {"id"},
				"filter_types":        {"in"},
				"filter_values":       {"test,test,2,3"},
				"filter_left_parens":  {"0"},
				"filter_right_parens": {"0"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE ((test_models.id in ($1, $2, $3,  $4)))", baseQuery),
			ExpectedArgs:  []string{"test", "test", "2", "3"},
		},
		{
			Name: "In Operator, 0 arg",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"in"},
				"filter_values":  {""},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (%s)", baseQuery, scope.FailQuery),
			ExpectedArgs:  []string{},
		},
		{
			Name: "In Operator, 1 arg",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"in"},
				"filter_values":  {"test"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id in ( $1))", baseQuery),
			ExpectedArgs:  []string{"test"},
		},
		{
			Name: "Not Like Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"nlk"},
				"filter_values":  {"%test%"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id not like $1)", baseQuery),
			ExpectedArgs:  []string{"%test%"},
		},
		{
			Name: "Not Ilike Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"nilk"},
				"filter_values":  {"%test%"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id not ilike $1)", baseQuery),
			ExpectedArgs:  []string{"%test%"},
		},
		{
			Name: "Not In Operator",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"nin"},
				"filter_values":  {"test,test,2,3"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id not in ($1, $2, $3,  $4))", baseQuery),
			ExpectedArgs:  []string{"test", "test", "2", "3"},
		},
		{
			Name: "Not In Operator, 0 arg",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"nin"},
				"filter_values":  {""},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (%s)", baseQuery, scope.FailQuery),
			ExpectedArgs:  []string{},
		},
		{
			Name: "Not In Operator, 1 arg",
			Params: map[string][]string{
				"filter_columns": {"id"},
				"filter_types":   {"nin"},
				"filter_values":  {"test"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id not in ( $1))", baseQuery),
			ExpectedArgs:  []string{"test"},
		},
		{
			Name: "In Operator with separator",
			Params: map[string][]string{
				"filter_columns":        {"id"},
				"filter_types":          {"in"},
				"filter_values":         {"test$test$2$3"},
				"filter_args_separator": {"$"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id in ($1, $2, $3,  $4))", baseQuery),
			ExpectedArgs:  []string{"test", "test", "2", "3"},
		},
		{
			Name: "Multiple Filters, IN",
			Params: map[string][]string{
				"filter_columns": {"id|null_id"},
				"filter_types":   {"in|nin"},
				"filter_values":  {"test|test2"},
				"filter_logic":   {"and"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id in ( $1) AND test_models.db_null_id not in ( $2))", baseQuery),
			ExpectedArgs:  []string{"test", "test2"},
		},
		{
			Name: "Multiple Filters, IN, multiple args",
			Params: map[string][]string{
				"filter_columns": {"id|null_id"},
				"filter_types":   {"in|nin"},
				"filter_values":  {"test,test2|test3"},
				"filter_logic":   {"and"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id in ($1,  $2) AND test_models.db_null_id not in ( $3))", baseQuery),
			ExpectedArgs:  []string{"test", "test2", "test3"},
		},
		{
			Name: "Multiple Filters, AND",
			Params: map[string][]string{
				"filter_columns": {"id|null_id"},
				"filter_types":   {"eq|ne"},
				"filter_values":  {"test|test2"},
				"filter_logic":   {"and"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id = $1 AND test_models.db_null_id != $2)", baseQuery),
			ExpectedArgs:  []string{"test", "test2"},
		},
		{
			Name: "Multiple Filters, AND, custom separator",
			Params: map[string][]string{
				"filter_columns":   {"id$null_id"},
				"filter_types":     {"eq$ne"},
				"filter_values":    {"test$test2"},
				"filter_logic":     {"and"},
				"filter_separator": {"$"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id = $1 AND test_models.db_null_id != $2)", baseQuery),
			ExpectedArgs:  []string{"test", "test2"},
		},
		{
			Name: "Multiple Filters, OR",
			Params: map[string][]string{
				"filter_columns": {"id|null_id"},
				"filter_types":   {"eq|ne"},
				"filter_values":  {"test|test2"},
				"filter_logic":   {"or"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id = $1 OR test_models.db_null_id != $2)", baseQuery),
			ExpectedArgs:  []string{"test", "test2"},
		},
		{
			Name: "Multiple Filters AND and OR",
			Params: map[string][]string{
				"filter_columns": {"id|null_id|id"},
				"filter_types":   {"eq|ne|ne"},
				"filter_values":  {"test|test2|test3"},
				"filter_logic":   {"or|and"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.id = $1 OR test_models.db_null_id != $2 AND test_models.id != $3)", baseQuery),
			ExpectedArgs:  []string{"test", "test2", "test3"},
		},
		{
			Name: "Multiple Filters AND and OR, PARENS",
			Params: map[string][]string{
				"filter_columns":      {"id|null_id|id"},
				"filter_types":        {"eq|ne|ne"},
				"filter_values":       {"test|test2|test3"},
				"filter_logic":        {"or|and"},
				"filter_left_parens":  {"0|1"},
				"filter_right_parens": {"2|2"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE ((test_models.id = $1 OR (test_models.db_null_id != $2 AND test_models.id != $3)))", baseQuery),
			ExpectedArgs:  []string{"test", "test2", "test3"},
		},
		{
			Name: "Multiple Filters Will Null Operator",
			Params: map[string][]string{
				"filter_columns": {"null_id|id"},
				"filter_types":   {"nu|eq"},
				"filter_values":  {"test|test2"},
				"filter_logic":   {"and"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s WHERE (test_models.db_null_id is null AND test_models.id = $1)", baseQuery),
			ExpectedArgs:  []string{"test2"},
		},
		{
			Name: "Multiple Filters Mismatched Fields",
			Params: map[string][]string{
				"filter_columns": {"null_id|id"},
				"filter_types":   {"nu,eq"},
				"filter_values":  {"test|test2"},
			},
			ExpectErr: true,
		},
		{
			Name: "Multiple Filters Mismatched Fields 2",
			Params: map[string][]string{
				"filter_columns": {"null_id"},
				"filter_types":   {"nu"},
				"filter_values":  {"test|test2"},
			},
			ExpectErr: true,
		},
	}

	for _, testCase := range testCases {
		ss.T().Run(testCase.Name, func(t *testing.T) {
			s, err := scope.ForFiltersFromParams(context.Background(), tm, url.Values(testCase.Params))
			if testCase.ExpectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			q := ss.DB.Session(&gorm.Session{DryRun: true}).Model(tm)
			q.Statement.SQL.Reset()
			scopeQueryFunc := q.Scopes(s).Find(&tm)

			args := scopeQueryFunc.Statement.Vars
			query := scopeQueryFunc.Statement.SQL.String()
			assert.Equal(t, testCase.ExpectedQuery, query)
			assert.Equal(t, len(testCase.ExpectedArgs), len(args))

			for i, arg := range testCase.ExpectedArgs {
				assert.Equal(t, arg, args[i])
			}
		})
	}
}

func (ss *ScopesSuite) TestForOrderFromParams() {
	tm := TestModel{}
	q := ss.DB.Session(&gorm.Session{DryRun: true}).Model(tm)
	q.Statement.SQL.Reset()
	scopeQueryFunc := q.Find(&tm)
	baseQuery := scopeQueryFunc.Statement.SQL.String()

	testCases := []struct {
		Name          string
		Params        map[string][]string
		ExpectErr     bool
		ExpectedQuery string
	}{
		{
			Name: "No Params",
			Params: map[string][]string{
				"sort_columns":    {""},
				"sort_directions": {""},
			},
			ExpectErr:     false,
			ExpectedQuery: baseQuery,
		},
		{
			Name: "Missing Operator",
			Params: map[string][]string{
				"sort_columns":    {"id"},
				"sort_directions": {""},
			},
			ExpectErr: true,
		},
		{
			Name: "Column Is Dash JSON",
			Params: map[string][]string{
				"sort_columns":    {"-"},
				"sort_directions": {"asc"},
			},
			ExpectErr: true,
		},
		{
			Name: "Column Is Dash DB",
			Params: map[string][]string{
				"sort_columns":    {"not_in_db"},
				"sort_directions": {"asc"},
			},
			ExpectErr: true,
		},
		{
			Name: "ASC Direction",
			Params: map[string][]string{
				"sort_columns":    {"id"},
				"sort_directions": {"asc"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s ORDER BY test_models.id ASC", baseQuery),
		},
		{
			Name: "ASC Direction Custom Column",
			Params: map[string][]string{
				"sort_columns":    {"custom_sort"},
				"sort_directions": {"asc"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s ORDER BY (SELECT '1234') ASC", baseQuery),
		},
		{
			Name: "DESC Direction",
			Params: map[string][]string{
				"sort_columns":    {"id"},
				"sort_directions": {"desc"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s ORDER BY test_models.id DESC", baseQuery),
		},
		{
			Name: "Multiple Sorts",
			Params: map[string][]string{
				"sort_columns":    {"id|null_id"},
				"sort_directions": {"asc|desc"},
			},
			ExpectErr:     false,
			ExpectedQuery: fmt.Sprintf("%s ORDER BY test_models.id ASC,test_models.db_null_id DESC", baseQuery),
		},
		{
			Name: "Multiple Sorts Mismatched Fields",
			Params: map[string][]string{
				"sort_columns":    {"id|null_id"},
				"sort_directions": {"asc,desc"},
			},
			ExpectErr: true,
		},
		{
			Name: "Multiple Sorts Mismatched Fields 2",
			Params: map[string][]string{
				"sort_columns":    {"id,null_id"},
				"sort_directions": {"asc|desc"},
			},
			ExpectErr: true,
		},
	}

	for _, testCase := range testCases {
		ss.T().Run(testCase.Name, func(t *testing.T) {
			s, err := scope.ForSortFromParams(context.Background(), tm, url.Values(testCase.Params))
			if testCase.ExpectErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			q := ss.DB.Session(&gorm.Session{DryRun: true}).Model(tm)
			q.Statement.SQL.Reset()
			scopeQueryFunc := q.Scopes(s).Find(&tm)

			query := scopeQueryFunc.Statement.SQL.String()
			assert.Equal(t, testCase.ExpectedQuery, query)
		})
	}
}
