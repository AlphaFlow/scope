package scope

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/alphaflow/api-core/destructify"
)

// CustomColumn represents a SQL statement that can be used like a column in a SQL query. CustomColumns are used in order
// to create filter and sort options the are more complex than simply sorting on a single field.  For example:
//
// Given an object with a db column 'test_int', in order to sort on the TEXT converted value of 'test_int' you would
// implement 'CustomSortable' and provide the CustomColumn:
//
//  {
//	  Name:       "test_int_text_value",
//	  Statement:  "test_int::TEXT",
//	  ResultType: reflect.TypeOf(""),
//   }
//
// Statement can be any valid sql Statement that returns a single value.
// ResultType is the type of the value returned by statement, and is used to scan the value from the DB.
type CustomColumn struct {
	Name       string
	Statement  string
	ResultType reflect.Type
}

type CustomColumns []CustomColumn

type CustomFilterable interface {
	GetCustomFilters(ctx context.Context) CustomColumns
}

type CustomSortable interface {
	GetCustomSorts(ctx context.Context) CustomColumns
}

// GetAllFilterColumns is a utility in order to automatically get a list of all columns that can be filtered on for
// the referenced model.
func GetAllFilterColumns(ctx context.Context, modelPtr interface{}) ([]CustomColumn, error) {
	v := reflect.ValueOf(modelPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, errors.New("pointer to struct expected")
	}

	model := v.Elem().Interface()

	validColumns, err := getAllQueryableColumns(modelPtr)
	if err != nil {
		return nil, err
	}

	validColumnsMap := make(map[string]CustomColumn, len(validColumns))
	for i, col := range validColumns {
		validColumnsMap[col.Name] = validColumns[i]
	}

	// Get all of the custom filters.
	if customFilterable, ok := model.(CustomFilterable); ok {
		customColumns := customFilterable.GetCustomFilters(ctx)

		for i, col := range customColumns {
			validColumnsMap[col.Name] = customColumns[i]
		}
	}

	allColumns := make([]CustomColumn, len(validColumnsMap))
	i := 0
	for j := range validColumnsMap {
		allColumns[i] = validColumnsMap[j]
		i++
	}

	return allColumns, nil
}

// GetAllFilterColumnNames is a utility in order to automatically get a list of all tags that can be filtered on for
// the referenced model.
func GetAllFilterColumnNames(ctx context.Context, modelPtr interface{}) ([]string, error) {
	validColumns, err := GetAllFilterColumns(ctx, modelPtr)
	if err != nil {
		return nil, err
	}

	validColumnNames := make([]string, len(validColumns))
	for i := range validColumns {
		validColumnNames[i] = validColumns[i].Name
	}

	return validColumnNames, nil
}

// GetAllSortColumns is a utility in order to automatically get a list of all columns that can be sorted on for
// the referenced model.
func GetAllSortColumns(ctx context.Context, modelPtr interface{}) ([]CustomColumn, error) {
	v := reflect.ValueOf(modelPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, errors.New("pointer to struct expected")
	}

	model := v.Elem().Interface()

	validColumns, err := getAllQueryableColumns(modelPtr)
	if err != nil {
		return nil, err
	}

	validColumnsMap := make(map[string]CustomColumn, len(validColumns))
	for i, col := range validColumns {
		validColumnsMap[col.Name] = validColumns[i]
	}

	// Get all of the custom sorts.
	if customSortable, ok := model.(CustomSortable); ok {
		customColumns := customSortable.GetCustomSorts(ctx)

		for i, col := range customColumns {
			validColumnsMap[col.Name] = customColumns[i]
		}
	}

	allColumns := make([]CustomColumn, len(validColumnsMap))
	i := 0
	for j := range validColumnsMap {
		allColumns[i] = validColumnsMap[j]
		i++
	}

	return allColumns, nil
}

// GetAllSortColumnNames is a utility in order to automatically get a list of all tags that can be sorted on for
// the referenced model.
func GetAllSortColumnNames(ctx context.Context, modelPtr interface{}) ([]string, error) {
	validColumns, err := GetAllSortColumns(ctx, modelPtr)
	if err != nil {
		return nil, err
	}

	validColumnNames := make([]string, len(validColumns))
	for i := range validColumns {
		validColumnNames[i] = validColumns[i].Name
	}

	return validColumnNames, nil
}

// GenerateCustomColumnsForSubobject is a utility in order to automatically create custom filter columns for a model
// related to the model you are filtering.
//
// For example, assume you have a table Houses, and each house has an address_id field pointing to an Addresses table.
// In order to make houses sortable by the columns of the Address table, you will need to implement 'CustomSortable' and
// return the following:
//
// GenerateCustomColumnsForSubobject(&Address{}, "address", "addresses.id = houses.address_id")
//
// Which will return a list of CustomColumns derived from the fields of the `&Address{}` model.
//  CustomColumns{
//    {
//	    Name:       "address.id",
//	    Statement:  "(select id from addresses where addresses.id = houses.address_id)",
//	    ResultType: reflect.TypeOf(""),
//    },
//    {
//	    Name:       "address.city",
//	    Statement:  "(select city from addresses where addresses.id = houses.address_id)",
//	    ResultType: reflect.TypeOf(""),
//    },
//    ...
//  }
func GenerateCustomColumnsForSubobject(subobjectPtr interface{}, subobjectJsonTag, joinClause string, optionalTablename ...string) (CustomColumns, error) {
	v := reflect.ValueOf(subobjectPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, errors.New("pointer to struct expected")
	}

	subobject := v.Interface()

	customColumns := CustomColumns{}

	fields := destructify.ValuesForStructTag(subobject, "json")
	for _, col := range fields {
		field, ok := destructify.FieldWithJsonTagValue(subobject, col)
		if !ok {
			continue
		}

		// Look up the database column name for this field.
		dbColumn, ok := destructify.LookupForStructFieldTag(subobject, field, "db")
		if !ok || dbColumn == "-" {
			continue
		}

		var tablename string
		if len(optionalTablename) > 0 {
			tablename = optionalTablename[0]
		} else {
			tablename = TableName(subobject)
		}

		customColumn := CustomColumn{
			Name:       fmt.Sprintf("%v.%v", subobjectJsonTag, col),
			ResultType: destructify.GetFieldByName(subobjectPtr, field).Type(),

			// Select [field_db_tag] from [subobject tablename] where [join clause]
			Statement: fmt.Sprintf("(select %v from %v where %v)", dbColumn, tablename, joinClause),
		}

		customColumns = append(customColumns, customColumn)
	}

	return customColumns, nil
}

// getAllQueryableColumns is a utility in order to automatically get a list of custom columns for all tags that can be
// filtered on this model by default, without including the interfaces CustomFilterable and CustomSortable.  In other
// words, this does not include any custom columns that may have been added, it only returns the columns on this model
// with both json and db tags.
func getAllQueryableColumns(modelPtr interface{}) ([]CustomColumn, error) {
	v := reflect.ValueOf(modelPtr)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return nil, errors.New("pointer to struct expected")
	}

	// Fetch the associated table for this model.
	tableName := TableName(modelPtr)

	model := v.Elem().Interface()

	validColumns := make([]CustomColumn, 0)

	// Get all of the json tags that are exposed to the user.
	structJsonTags := destructify.ValuesForStructTag(model, "json")
	for i, structJsonTag := range structJsonTags {
		structFieldName, ok := destructify.FieldWithJsonTagValue(model, structJsonTag)
		if !ok {
			continue
		}

		dbColumn, ok := destructify.LookupForStructFieldTag(model, structFieldName, "db")
		if !ok || dbColumn == "-" {
			continue
		}

		customColumn := CustomColumn{
			Name:       structJsonTags[i],
			ResultType: destructify.GetFieldByName(modelPtr, structFieldName).Type(),
			Statement:  fmt.Sprintf("%v.%v", tableName, dbColumn),
		}

		validColumns = append(validColumns, customColumn)
	}

	return validColumns, nil
}
