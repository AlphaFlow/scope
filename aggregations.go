package scope

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"

	"github.com/alphaflow/scope/util"
)

type Aggregation struct {
	Name       string
	Statement  string
	ResultType reflect.Type
}

type StandardAggregationsType string

const (
	StandardAggregationsTypeCount StandardAggregationsType = "COUNT"
	StandardAggregationsTypeSum   StandardAggregationsType = "SUM"
	StandardAggregationsTypeAvg   StandardAggregationsType = "AVG"
	StandardAggregationsTypeMax   StandardAggregationsType = "MAX"
	StandardAggregationsTypeMin   StandardAggregationsType = "MIN"
)

var StandardAggregations = map[StandardAggregationsType]Aggregation{
	StandardAggregationsTypeCount: {
		Name:       string(StandardAggregationsTypeCount),
		Statement:  string(StandardAggregationsTypeCount),
		ResultType: reflect.TypeOf(0),
	},
	StandardAggregationsTypeSum: {
		Name:       string(StandardAggregationsTypeSum),
		Statement:  string(StandardAggregationsTypeSum),
		ResultType: nil,
	},
	StandardAggregationsTypeAvg: {
		Name:       string(StandardAggregationsTypeAvg),
		Statement:  string(StandardAggregationsTypeAvg),
		ResultType: nil,
	},
	StandardAggregationsTypeMax: {
		Name:       string(StandardAggregationsTypeMax),
		Statement:  string(StandardAggregationsTypeMax),
		ResultType: nil,
	},
	StandardAggregationsTypeMin: {
		Name:       string(StandardAggregationsTypeMin),
		Statement:  string(StandardAggregationsTypeMin),
		ResultType: nil,
	},
}

// aggregationsQueryResult is a struct with an interface column.  The type of interface is swapped out using
// reflection in getCustomAggregations in order to be able to scan DB values into any type as needed.
type aggregationsQueryResult struct {
	Result interface{} `db:"result"`
}

// groupedAggregationsQueryResult is a struct with an interface column.  The type of interface is swapped out using
// reflection in getCustomGroupedAggregations in order to be able to scan DB values into any type as needed.
type groupedAggregationsQueryResult struct {
	Grouper interface{} `db:"grouper"`
	Result  interface{} `db:"result"`
}

// GetAggregationsFromParams aggregates a modelsPtr based on params, restricting by the scope collection scopes.
func GetAggregationsFromParams(ctx context.Context, tx *pop.Connection, modelsPtr interface{}, params buffalo.ParamValues, scopes *Collection) (interface{}, error) {
	aggregation, ok := StandardAggregations[StandardAggregationsType(strings.ToUpper(params.Get("aggregation_type")))]
	if !ok {
		return nil, errors.New("unknown aggregation type")
	}

	aggregationResult, err := GetAggregations(ctx, tx, modelsPtr, params.Get("aggregation_column"), scopes, aggregation)
	if err != nil {
		return nil, err
	}

	return aggregationResult, nil
}

// GetGroupedAggregationsFromParams groups and aggregates a modelsPtr based on params, restricting by the scope collection scopes.
func GetGroupedAggregationsFromParams(ctx context.Context, tx *pop.Connection, modelsPtr interface{}, params buffalo.ParamValues, scopes *Collection) ([]interface{}, error) {
	aggregation, ok := StandardAggregations[StandardAggregationsType(strings.ToUpper(params.Get("aggregation_type")))]
	if !ok {
		return nil, errors.New("unknown aggregation type")
	}

	aggregationResult, err := GetGroupedAggregations(ctx, tx, modelsPtr, params.Get("aggregation_column"), params.Get("aggregation_grouper_column"), scopes, aggregation)
	if err != nil {
		return nil, err
	}

	return aggregationResult, nil
}

// GetAggregations returns the aggregated value for column `columnName` of modelsPtr, restricting by the scope collection
// `scopes`.
//
// `columnName` is either a CustomColumn returned by the CustomFilterable interface, or a field specified by the json
// tag.  This is the same as the acceptable values for `filter_columns` in ForFiltersFromParams.
func GetAggregations(ctx context.Context, tx *pop.Connection, modelsPtr interface{}, columnName string, scopes *Collection, aggregation Aggregation) (interface{}, error) {
	v := reflect.ValueOf(modelsPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return nil, errors.New("pointer to slice expected")
	}

	models := v.Elem()
	modelPtr := reflect.New(models.Type().Elem()).Interface()

	filterColumns, err := GetAllFilterColumns(ctx, modelPtr)
	if err != nil {
		return nil, err
	}

	var column *CustomColumn
	for i, filterColumn := range filterColumns {
		if columnName == filterColumn.Name {
			column = &filterColumns[i]
			break
		}
	}

	if column == nil {
		return nil, errors.Errorf("invalid filter field: %v", columnName)
	}

	tableName := (&pop.Model{Value: modelPtr}).TableName()
	return getCustomAggregations(tx, tableName, *column, scopes, aggregation)
}

// GetGroupedAggregations returns the aggregated value for column `columnName` of modelsPtr, grouped by `grouperName` of
// modelsPtr, restricting by the scope collection `scopes`.
//
// `columnName` is either a CustomColumn returned by the CustomFilterable interface, or a field specified by the json
// tag.  This is the same as the acceptable values for `filter_columns` in ForFiltersFromParams.
func GetGroupedAggregations(ctx context.Context, tx *pop.Connection, modelsPtr interface{}, columnName, grouperName string, scopes *Collection, aggregation Aggregation) ([]interface{}, error) {
	v := reflect.ValueOf(modelsPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		return nil, errors.New("pointer to slice expected")
	}

	models := v.Elem()
	modelPtr := reflect.New(models.Type().Elem()).Interface()

	filterColumns, err := GetAllFilterColumns(ctx, modelPtr)
	if err != nil {
		return nil, err
	}

	var column *CustomColumn
	var grouper *CustomColumn
	for i, filterColumn := range filterColumns {
		if columnName == filterColumn.Name {
			column = &filterColumns[i]
		}
		if grouperName == filterColumn.Name {
			grouper = &filterColumns[i]
		}
		if column != nil && grouper != nil {
			break
		}
	}

	if column == nil {
		return nil, errors.Errorf("invalid filter field: %v", columnName)
	}
	if grouper == nil {
		return nil, errors.Errorf("invalid filter field: %v", grouperName)
	}

	tableName := (&pop.Model{Value: modelPtr}).TableName()
	return getCustomGroupedAggregations(tx, tableName, *column, *grouper, scopes, aggregation)
}

// getCustomAggregations returns the aggregated value for column for the provided `customColumn` from the table `tableName`,
// after scoping said table by `scopes`.
//
// In order to do this, we must build a custom struct with the correct ResultType for customColumn.  We then build a
// scoped GROUP BY query to retrieve all values for customColumn into that struct.
func getCustomAggregations(tx *pop.Connection, tableName string, customColumn CustomColumn, scopes *Collection, aggregation Aggregation) (interface{}, error) {
	type __stub__ struct{}
	clauses := ""

	// We need to build a struct of type "ResultType", so that we can correctly marshall the output types from the DB.
	templateStructField, ok := reflect.ValueOf(aggregationsQueryResult{}).Type().FieldByName("Result")
	if !ok {
		return nil, errors.New("unable to build aggregation query result")
	}

	templateStructField.Type = aggregation.ResultType

	// If the aggregation doesn't have a defined output type, assume the output is the same type as the field.
	if templateStructField.Type == nil {
		templateStructField.Type = customColumn.ResultType
	}

	//templateStructFieldDBTag := templateStructField.Tag.Get("db")
	typedStructWithDBTag := reflect.New(reflect.StructOf([]reflect.StructField{templateStructField}))
	aggregationScopes := NewCollection(tx)

	// We never return null as a filter option.
	aggregationScopes.Push(ForNotNull(customColumn.Statement))

	if scopes != nil && len(scopes.scopes) > 0 {
		aggregationScopes.Push(scopes.scopes...)
	}

	scopeQueryFunc := aggregationScopes.Flatten()(tx.Q())
	scopeQuerySQL, scopeQueryArgs := scopeQueryFunc.ToSQL(&pop.Model{Value: __stub__{}})
	stubRegex := regexp.MustCompile(`^SELECT\s+FROM stubs AS stubs\s+`)
	clauses = stubRegex.ReplaceAllString(scopeQuerySQL, "")

	// Strip all order by columns out of the query, since they don't matter.
	orderRegex := regexp.MustCompile(`ORDER\s+BY\s+\w+(\s+ASC|\s+DESC)?([\s,]*\w+(\s+ASC|\s+DESC)?)*`)
	clauses = orderRegex.ReplaceAllString(clauses, "")

	generatedStatement := fmt.Sprintf("SELECT %v(%v) AS result FROM %v %v", aggregation.Statement, customColumn.Statement, tableName, clauses)
	err := tx.RawQuery(generatedStatement, scopeQueryArgs...).First(typedStructWithDBTag.Interface())
	if err != nil {
		return nil, err
	}

	return reflect.Indirect(typedStructWithDBTag).FieldByName("Result").Interface(), nil
}

// getCustomGroupedAggregations returns the aggregated value for column for the provided `customColumn` from the table `tableName`,
// after scoping said table by `scopes`, grouping by `groupColumn`.
//
// In order to do this, we must build a custom struct with the correct ResultType for customColumn.  We then build a
// scoped GROUP BY query to retrieve all values for customColumn into that struct.
func getCustomGroupedAggregations(tx *pop.Connection, tableName string, customColumn, groupColumn CustomColumn, scopes *Collection, aggregation Aggregation) ([]interface{}, error) {
	type __stub__ struct{}
	clauses := ""

	// We need to build a struct of type "ResultType", so that we can correctly marshall the output types from the DB.
	templateResultStructField, ok := reflect.ValueOf(groupedAggregationsQueryResult{}).Type().FieldByName("Result")
	if !ok {
		return nil, errors.New("unable to build grouped aggregation query result")
	}
	templateGrouperStructField, ok := reflect.ValueOf(groupedAggregationsQueryResult{}).Type().FieldByName("Grouper")
	if !ok {
		return nil, errors.New("unable to build grouped aggregation query result")
	}

	templateResultStructField.Type = aggregation.ResultType
	templateGrouperStructField.Type = groupColumn.ResultType

	// If the aggregation doesn't have a defined output type, assume the output is the same type as the field.
	if templateResultStructField.Type == nil {
		templateResultStructField.Type = customColumn.ResultType
	}

	typedStructArrayPtrWithDBTag := reflect.New(reflect.SliceOf(reflect.StructOf([]reflect.StructField{templateGrouperStructField, templateResultStructField})))
	aggregationScopes := NewCollection(tx)

	// We never return null as a filter option.
	aggregationScopes.Push(ForNotNull(customColumn.Statement))

	if scopes != nil && len(scopes.scopes) > 0 {
		aggregationScopes.Push(scopes.scopes...)
	}

	scopeQueryFunc := aggregationScopes.Flatten()(tx.Q())
	scopeQuerySQL, scopeQueryArgs := scopeQueryFunc.ToSQL(&pop.Model{Value: __stub__{}})
	stubRegex := regexp.MustCompile(`^SELECT\s+FROM stubs AS stubs\s+`)
	clauses = stubRegex.ReplaceAllString(scopeQuerySQL, "")

	// Strip all order by columns out of the query, since they don't matter.
	orderRegex := regexp.MustCompile(`ORDER\s+BY\s+\w+(\s+ASC|\s+DESC)?([\s,]*\w+(\s+ASC|\s+DESC)?)*`)
	clauses = orderRegex.ReplaceAllString(clauses, "")

	generatedStatement := fmt.Sprintf("SELECT %v AS grouper, %v(%v) AS result FROM %v %v GROUP BY %v", groupColumn.Statement, aggregation.Statement, customColumn.Statement, tableName, clauses, groupColumn.Statement)
	err := tx.RawQuery(generatedStatement, scopeQueryArgs...).All(typedStructArrayPtrWithDBTag.Interface())
	if err != nil {
		return nil, err
	}

	return util.InterfaceSlice(reflect.Indirect(typedStructArrayPtrWithDBTag).Interface()), nil
}
