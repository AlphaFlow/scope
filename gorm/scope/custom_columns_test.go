package scope_test

import (
	"context"
	"reflect"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"

	"github.com/alphaflow/api-core/gorm/scope"
)

type TestObjectWithOverride struct {
	ID       uuid.UUID  `json:"id" db:"id" gorm:"primaryKey;column:id"`
	ObjectId nulls.UUID `json:"object_id" db:"object_id" gorm:"column:object_id"`
}

func (t TestObjectWithOverride) GetCustomFilters(ctx context.Context) scope.CustomColumns {
	overrideFilter := scope.CustomColumn{
		Name:       "object_id",
		Statement:  `NULL`,
		ResultType: reflect.TypeOf(nulls.UUID{}),
	}
	return scope.CustomColumns{overrideFilter}
}

func (t TestObjectWithOverride) GetCustomSorts(ctx context.Context) scope.CustomColumns {
	overrideSort := scope.CustomColumn{
		Name:       "object_id",
		Statement:  `NULL`,
		ResultType: reflect.TypeOf(nulls.UUID{}),
	}
	return scope.CustomColumns{overrideSort}
}

type TestSubObject struct {
	ID       uuid.UUID  `json:"id" db:"id" gorm:"primaryKey;column:id"`
	ObjectId nulls.UUID `json:"object_id" db:"object_id" gorm:"column:object_id"`
}

func (ss *ScopesSuite) TestGetAllFilterColumns() {
	testObject := &TestObject{}

	filters, err := scope.GetAllFilterColumns(context.Background(), testObject)
	ss.NoError(err)

	validColumnNames := make([]string, len(filters))
	for i := range filters {
		validColumnNames[i] = filters[i].Name
	}

	expectedFilters := []string{"id", "null_id", "num", "custom_filter", "custom_uuid_filter", "custom_nulls_uuid_filter", "null_filter"}
	ss.ElementsMatch(expectedFilters, validColumnNames)
}

func (ss *ScopesSuite) TestGetAllFilterColumns_override() {
	testObject := &TestObjectWithOverride{}

	filters, err := scope.GetAllFilterColumns(context.Background(), testObject)
	ss.NoError(err)

	validColumnNames := make([]string, len(filters))
	for i := range filters {
		validColumnNames[i] = filters[i].Name
	}

	expectedFilters := []string{"id", "object_id"}
	ss.ElementsMatch(expectedFilters, validColumnNames)

	filtersMap := make(map[string]scope.CustomColumn, len(filters))
	for i, col := range filters {
		filtersMap[col.Name] = filters[i]
	}
	ss.Equal(`NULL`, filtersMap["object_id"].Statement)
}

func (ss *ScopesSuite) TestGetAllFilterColumnNames() {
	testObject := &TestObject{}

	filters, err := scope.GetAllFilterColumnNames(context.Background(), testObject)
	ss.NoError(err)

	expectedFilters := []string{"id", "null_id", "num", "custom_filter", "custom_uuid_filter", "custom_nulls_uuid_filter", "null_filter"}
	ss.ElementsMatch(expectedFilters, filters)
}

func (ss *ScopesSuite) TestGetAllSortColumns() {
	testObject := &TestObject{}

	sorts, err := scope.GetAllSortColumns(context.Background(), testObject)
	ss.NoError(err)

	validColumnNames := make([]string, len(sorts))
	for i := range sorts {
		validColumnNames[i] = sorts[i].Name
	}

	expectedSorts := []string{"id", "null_id", "num", "custom_sort", "custom_uuid_sort", "custom_nulls_uuid_sort", "null_sort"}
	ss.ElementsMatch(expectedSorts, validColumnNames)
}

func (ss *ScopesSuite) TestGetAllSortColumns_override() {
	testObject := &TestObjectWithOverride{}

	sorts, err := scope.GetAllSortColumns(context.Background(), testObject)
	ss.NoError(err)

	validColumnNames := make([]string, len(sorts))
	for i := range sorts {
		validColumnNames[i] = sorts[i].Name
	}

	expectedSorts := []string{"id", "object_id"}
	ss.ElementsMatch(expectedSorts, validColumnNames)

	sortsMap := make(map[string]scope.CustomColumn, len(sorts))
	for i, col := range sorts {
		sortsMap[col.Name] = sorts[i]
	}
	ss.Equal(`NULL`, sortsMap["object_id"].Statement)
}

func (ss *ScopesSuite) TestGetAllSortColumnNames() {
	testObject := &TestObject{}

	sorts, err := scope.GetAllSortColumnNames(context.Background(), testObject)
	ss.NoError(err)

	expectedSorts := []string{"id", "null_id", "num", "custom_sort", "custom_uuid_sort", "custom_nulls_uuid_sort", "null_sort"}
	ss.ElementsMatch(expectedSorts, sorts)
}

func (ss *ScopesSuite) TestGenerateSubobjectCustomFilterColumns() {
	testSubObject := &TestSubObject{}

	customFilterColumns, err := scope.GenerateCustomColumnsForSubobject(testSubObject, "subobject", "id = subobject.object_id")
	ss.NoError(err)

	expectedCustomFilterColumns := scope.CustomColumns{
		scope.CustomColumn{
			Name:       "subobject.id",
			ResultType: reflect.TypeOf(uuid.Nil),
			Statement:  "(select id from test_sub_objects where id = subobject.object_id)",
		},
		scope.CustomColumn{
			Name:       "subobject.object_id",
			ResultType: reflect.TypeOf(nulls.UUID{}),
			Statement:  "(select object_id from test_sub_objects where id = subobject.object_id)",
		},
	}

	ss.Equal(expectedCustomFilterColumns, customFilterColumns)
}
