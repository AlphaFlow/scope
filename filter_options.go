package scope

import (
	"context"
	"fmt"
	"reflect"
	"regexp"

	"github.com/alphaflow/scope/util"
	"github.com/gobuffalo/pop/v5"
	"github.com/pkg/errors"
)

// filterOptionsQueryResult is a struct with an interface column.  The type of interface is swapped out using
// reflection in getCustomFilterOptions in order to be able to scan DB values into any type as needed.
type filterOptionsQueryResult struct {
	Result interface{} `db:"result"`
}

// GetFilterOptions returns all of unique values for column 'columnName' of modelsPtr, restricting by the scope collection
// scopes.
//
// 'columnName' is either a CustomColumn returned by the CustomFilterable interface, or a field specified by the json
// tag.  This is the same as the acceptable values for 'filter_columns' in ForFiltersFromParams.
func GetFilterOptions(ctx context.Context, tx *pop.Connection, modelsPtr interface{}, columnName string, scopes *Collection) ([]interface{}, error) {
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
	return getCustomFilterOptions(tx, tableName, *column, scopes)
}

// getCustomFilterOptions returns all of the potential values for the provided 'customColumn' from the table 'tableName',
// after scoping said table by 'scopes'.
//
// In order to do this, we must build a custom struct with the correct ResultType for customColumn.  We then build a
// scoped GROUP BY query to retrieve all values for customColumn into that struct.
func getCustomFilterOptions(tx *pop.Connection, tableName string, customColumn CustomColumn, scopes *Collection) ([]interface{}, error) {
	type __stub__ struct{}
	clauses := ""

	// We need to build a struct of type "ResultType", so that we can correctly marshall the output types from the DB.
	templateStructField := reflect.ValueOf(filterOptionsQueryResult{}).Type().Field(0)
	templateStructField.Type = customColumn.ResultType
	templateStructFieldDBTag := templateStructField.Tag.Get("db")
	typedStructWithDBTag := reflect.StructOf([]reflect.StructField{templateStructField})
	typedStructArrayPtrWithDBTag := reflect.New(reflect.SliceOf(typedStructWithDBTag))
	filterOptionScopes := NewCollection(tx)

	// We never return null as a filter option.
	filterOptionScopes.Push(ForNotNull(customColumn.Statement))

	if scopes != nil && len(scopes.scopes) > 0 {
		filterOptionScopes.Push(scopes.scopes...)
	}

	// Covert this query to a GROUP BY query in order to only get distinct results.
	scopeQueryFunc := filterOptionScopes.Flatten()(tx.Q()).GroupBy(templateStructFieldDBTag)
	scopeQuerySQL, scopeQueryArgs := scopeQueryFunc.ToSQL(&pop.Model{Value: __stub__{}})
	stubRegex := regexp.MustCompile(`^SELECT\s+FROM stubs AS stubs\s+`)
	clauses = stubRegex.ReplaceAllString(scopeQuerySQL, "")

	// Strip all order by columns out of the query, since they don't matter and will break our GROUP BY.
	orderRegex := regexp.MustCompile(`ORDER\s+BY\s+\w+(\s+ASC|\s+DESC)?([\s,]*\w+(\s+ASC|\s+DESC)?)*`)
	clauses = orderRegex.ReplaceAllString(clauses, "")

	generatedStatement := fmt.Sprintf("select %v as result from %v %v", customColumn.Statement, tableName, clauses)
	err := tx.RawQuery(generatedStatement, scopeQueryArgs...).All(typedStructArrayPtrWithDBTag.Interface())
	if err != nil {
		return nil, err
	}

	resultLen := reflect.Indirect(typedStructArrayPtrWithDBTag).Len()
	output := reflect.MakeSlice(reflect.SliceOf(customColumn.ResultType), resultLen, resultLen)
	for i := 0; i < resultLen; i++ {
		value := reflect.Indirect(typedStructArrayPtrWithDBTag).Index(i).FieldByName(templateStructField.Name).Convert(customColumn.ResultType)
		output.Index(i).Set(value)
	}

	return util.InterfaceSlice(output.Interface()), nil
}
