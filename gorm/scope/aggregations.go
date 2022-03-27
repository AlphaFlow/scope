package scope

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
	"gorm.io/gorm"

	"github.com/alphaflow/scope/util"
)

type Aggregation struct {
	Name       string
	Statement  string
	ResultType reflect.Type
}

type Aggregations []Aggregation
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
func GetAggregationsFromParams(ctx context.Context, tx *gorm.DB, modelsPtr interface{}, params buffalo.ParamValues, scopes *Collection) (interface{}, error) {
	filterSeparator := getFilterSeparator(params)

	columns := make([]string, 0)
	if !util.IsBlank(params.Get("aggregation_column")) {
		columns = strings.Split(params.Get("aggregation_column"), filterSeparator)
	}
	types := make([]string, 0)
	if !util.IsBlank(params.Get("aggregation_type")) {
		types = strings.Split(params.Get("aggregation_type"), filterSeparator)
	}

	if len(columns) != len(types) {
		// We must have the same number of all aggregation params.
		return nil, errors.New("missing or mismatched aggregation parameters")
	}

	aggregations := Aggregations{}
	for _, aggregationType := range types {
		aggregation, ok := StandardAggregations[StandardAggregationsType(strings.ToUpper(aggregationType))]
		if !ok {
			return nil, errors.New("unknown aggregation type")
		}

		aggregations = append(aggregations, aggregation)
	}

	aggregationResult, err := GetAggregations(ctx, tx, modelsPtr, columns, scopes, aggregations)
	if err != nil {
		return nil, err
	}

	return aggregationResult, nil
}

// GetGroupedAggregationsFromParams groups and aggregates a modelsPtr based on params, restricting by the scope collection scopes.
func GetGroupedAggregationsFromParams(ctx context.Context, tx *gorm.DB, modelsPtr interface{}, params buffalo.ParamValues, scopes *Collection) ([]interface{}, error) {
	filterSeparator := getFilterSeparator(params)

	columns := make([]string, 0)
	if !util.IsBlank(params.Get("aggregation_column")) {
		columns = strings.Split(params.Get("aggregation_column"), filterSeparator)
	}
	types := make([]string, 0)
	if !util.IsBlank(params.Get("aggregation_type")) {
		types = strings.Split(params.Get("aggregation_type"), filterSeparator)
	}

	if len(columns) != len(types) {
		// We must have the same number of all aggregation params, and only 1 grouper.
		return nil, errors.New("missing or mismatched aggregation parameters")
	}

	aggregations := Aggregations{}
	for _, aggregationType := range types {
		aggregation, ok := StandardAggregations[StandardAggregationsType(strings.ToUpper(aggregationType))]
		if !ok {
			return nil, errors.New("unknown aggregation type")
		}

		aggregations = append(aggregations, aggregation)
	}

	aggregationGrouperColumn := params.Get("aggregation_grouper_column")

	aggregationResult, err := GetGroupedAggregations(ctx, tx, modelsPtr, columns, aggregationGrouperColumn, scopes, aggregations)
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
func GetAggregations(ctx context.Context, tx *gorm.DB, modelsPtr interface{}, columnNames []string, scopes *Collection, aggregations Aggregations) (interface{}, error) {
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

	customColumns := CustomColumns{}
	for _, columnName := range columnNames {
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

		customColumns = append(customColumns, *column)
	}

	tableName := TableName(modelPtr)
	return getCustomAggregations(tx, tableName, customColumns, scopes, aggregations)
}

// GetGroupedAggregations returns the aggregated value for column `columnName` of modelsPtr, grouped by `grouperName` of
// modelsPtr, restricting by the scope collection `scopes`.
//
// `columnName` is either a CustomColumn returned by the CustomFilterable interface, or a field specified by the json
// tag.  This is the same as the acceptable values for `filter_columns` in ForFiltersFromParams.
func GetGroupedAggregations(ctx context.Context, tx *gorm.DB, modelsPtr interface{}, columnNames []string, grouperName string, scopes *Collection, aggregations Aggregations) ([]interface{}, error) {
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

	customColumns := CustomColumns{}
	var grouper *CustomColumn
	for _, columnName := range columnNames {
		var column *CustomColumn
		for i, filterColumn := range filterColumns {
			if columnName == filterColumn.Name {
				column = &filterColumns[i]
			}

			if grouper == nil && grouperName == filterColumn.Name {
				grouper = &filterColumns[i]
			}

			if grouper != nil && column != nil {
				break
			}
		}

		if column == nil {
			return nil, errors.Errorf("invalid filter field: %v", columnName)
		}

		customColumns = append(customColumns, *column)
	}
	if grouper == nil {
		return nil, errors.Errorf("invalid filter field: %v", grouperName)
	}

	tableName := TableName(modelPtr)
	return getCustomGroupedAggregations(tx, tableName, customColumns, *grouper, scopes, aggregations)
}

// getCustomAggregations returns the aggregated value for column for the provided `customColumn` from the table `tableName`,
// after scoping said table by `scopes`.
//
// In order to do this, we must build a custom struct with the correct ResultType for customColumn.  We then build a
// scoped GROUP BY query to retrieve all values for customColumn into that struct.
func getCustomAggregations(tx *gorm.DB, tableName string, customColumns CustomColumns, scopes *Collection, aggregations Aggregations) (interface{}, error) {
	clauses := ""

	structFields := make([]reflect.StructField, len(aggregations))
	queryStubs := make([]string, len(aggregations))
	aggregationScopes := NewCollection(tx)
	jsonKeySet := make(map[string]bool)
	for i, aggregation := range aggregations {
		jsonKey := fmt.Sprintf("%v_%v", strings.ToLower(aggregation.Name), customColumns[i].Name)
		if _, ok := jsonKeySet[jsonKey]; ok {
			return nil, errors.New("duplicate aggregation parameter")
		}

		jsonKeySet[jsonKey] = true
		structFieldName := fmt.Sprintf("Result%v", i)
		structFieldTag := fmt.Sprintf(`db:"result%v" json:"%v"`, i, jsonKey)

		// We need to build a struct of type "ResultType", so that we can correctly marshall the output types from the DB.
		templateStructField, ok := reflect.ValueOf(aggregationsQueryResult{}).Type().FieldByName("Result")
		if !ok {
			return nil, errors.New("unable to build aggregation query result")
		}

		templateStructField.Name = structFieldName
		templateStructField.Type = aggregation.ResultType
		templateStructField.Tag = reflect.StructTag(structFieldTag)

		// If the aggregation doesn't have a defined output type, assume the output is the same type as the field.
		if templateStructField.Type == nil {
			templateStructField.Type = customColumns[i].ResultType
		}

		trailingComma := ","
		if i == len(aggregations)-1 {
			trailingComma = ""
		}
		structFields[i] = templateStructField
		queryStubs[i] = fmt.Sprintf("%v(%v) AS result%v%v", aggregation.Statement, customColumns[i].Statement, i, trailingComma)

		// We never return null as a filter option.
		aggregationScopes.Push(ForNotNull(customColumns[i].Statement))
	}

	//templateStructFieldDBTag := templateStructField.Tag.Get("db")
	typedStructWithDBTag := reflect.New(reflect.StructOf(structFields))

	if scopes != nil && len(scopes.scopes) > 0 {
		aggregationScopes.Push(scopes.scopes...)
	}

	q := tx.Session(&gorm.Session{DryRun: true}).Model(__stub__{})
	q.Statement.SQL.Reset()
	scopeQueryFunc := aggregationScopes.Flatten()(q).Find(q.Statement.Model)
	scopeQuerySQL := scopeQueryFunc.Statement.SQL.String()
	scopeQuerySQL = strings.Replace(scopeQuerySQL, "SELECT *", "SELECT ", 1)
	scopeQueryArgs := scopeQueryFunc.Statement.Vars
	stubRegex := regexp.MustCompile(stubRegex)
	clauses = stubRegex.ReplaceAllString(scopeQuerySQL, "")

	// Strip all order by columns out of the query, since they don't matter.
	orderRegex := regexp.MustCompile(`ORDER\s+BY\s+\w+(\s+ASC|\s+DESC)?([\s,]*\w+(\s+ASC|\s+DESC)?)*`)
	clauses = orderRegex.ReplaceAllString(clauses, "")

	generatedStatement := fmt.Sprintf("SELECT %v FROM %v %v", strings.Join(queryStubs, " "), tableName, clauses)
	err := tx.Raw(generatedStatement, scopeQueryArgs...).First(typedStructWithDBTag.Interface()).Error
	if err != nil {
		return nil, err
	}

	return typedStructWithDBTag.Interface(), nil
}

// getCustomGroupedAggregations returns the aggregated value for column for the provided `customColumn` from the table `tableName`,
// after scoping said table by `scopes`, grouping by `groupColumn`.
//
// In order to do this, we must build a custom struct with the correct ResultType for customColumn.  We then build a
// scoped GROUP BY query to retrieve all values for customColumn into that struct.
func getCustomGroupedAggregations(tx *gorm.DB, tableName string, customColumns CustomColumns, groupColumn CustomColumn, scopes *Collection, aggregations Aggregations) ([]interface{}, error) {
	clauses := ""

	structFields := make([]reflect.StructField, len(aggregations)+1)

	templateGrouperStructField, ok := reflect.ValueOf(groupedAggregationsQueryResult{}).Type().FieldByName("Grouper")
	if !ok {
		return nil, errors.New("unable to build grouped aggregation query result")
	}

	templateGrouperStructField.Type = groupColumn.ResultType
	structFields[0] = templateGrouperStructField

	queryStubs := make([]string, len(aggregations))
	aggregationScopes := NewCollection(tx)
	jsonKeySet := make(map[string]bool)
	for i, aggregation := range aggregations {
		jsonKey := fmt.Sprintf("%v_%v", strings.ToLower(aggregation.Name), customColumns[i].Name)
		if _, ok := jsonKeySet[jsonKey]; ok {
			return nil, errors.New("duplicate aggregation parameter")
		}

		jsonKeySet[jsonKey] = true
		structFieldName := fmt.Sprintf("Result%v", i)
		structFieldTag := fmt.Sprintf(`db:"result%v" json:"%v"`, i, jsonKey)

		// We need to build a struct of type "ResultType", so that we can correctly marshall the output types from the DB.
		templateStructField, ok := reflect.ValueOf(aggregationsQueryResult{}).Type().FieldByName("Result")
		if !ok {
			return nil, errors.New("unable to build aggregation query result")
		}

		templateStructField.Name = structFieldName
		templateStructField.Type = aggregation.ResultType
		templateStructField.Tag = reflect.StructTag(structFieldTag)

		// If the aggregation doesn't have a defined output type, assume the output is the same type as the field.
		if templateStructField.Type == nil {
			templateStructField.Type = customColumns[i].ResultType
		}

		trailingComma := ","
		if i == len(aggregations)-1 {
			trailingComma = ""
		}

		structFields[i+1] = templateStructField
		queryStubs[i] = fmt.Sprintf("%v(%v) AS result%v%v", aggregation.Statement, customColumns[i].Statement, i, trailingComma)

		// We never return null as a filter option.
		aggregationScopes.Push(ForNotNull(customColumns[i].Statement))
	}

	typedStructArrayPtrWithDBTag := reflect.New(reflect.SliceOf(reflect.StructOf(structFields)))

	if scopes != nil && len(scopes.scopes) > 0 {
		aggregationScopes.Push(scopes.scopes...)
	}

	q := tx.Session(&gorm.Session{DryRun: true}).Model(__stub__{})
	q.Statement.SQL.Reset()
	scopeQueryFunc := aggregationScopes.Flatten()(q).Find(q.Statement.Model)
	scopeQuerySQL := scopeQueryFunc.Statement.SQL.String()
	scopeQuerySQL = strings.Replace(scopeQuerySQL, "SELECT *", "SELECT ", 1)
	scopeQueryArgs := scopeQueryFunc.Statement.Vars
	stubRegex := regexp.MustCompile(stubRegex)
	clauses = stubRegex.ReplaceAllString(scopeQuerySQL, "")

	// Strip all order by columns out of the query, since they don't matter.
	orderRegex := regexp.MustCompile(`ORDER\s+BY\s+\w+(\s+ASC|\s+DESC)?([\s,]*\w+(\s+ASC|\s+DESC)?)*`)
	clauses = orderRegex.ReplaceAllString(clauses, "")

	generatedStatement := fmt.Sprintf("SELECT %v AS grouper, %v FROM %v %v GROUP BY %v", groupColumn.Statement, strings.Join(queryStubs, " "), tableName, clauses, groupColumn.Statement)
	err := tx.Raw(generatedStatement, scopeQueryArgs...).Find(typedStructArrayPtrWithDBTag.Interface()).Error
	if err != nil {
		return nil, err
	}

	return util.InterfaceSlice(reflect.Indirect(typedStructArrayPtrWithDBTag).Interface()), nil
}
