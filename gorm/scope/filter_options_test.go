package scope_test

import (
	"context"
	"reflect"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"

	"github.com/alphaflow/scope/gorm/scope"
)

type TestObject struct {
	ID        uuid.UUID  `json:"id" db:"id" gorm:"primaryKey;column:id;default:uuid_generate_v4()"`
	Nuid      nulls.UUID `json:"null_id" db:"db_null_id" gorm:"column:db_null_id"`
	Number    float64    `json:"num" db:"num" gorm:"column:num"`
	NotInDb   int        `json:"not_in_db" db:"-" gorm:"-"`
	NotInJson int        `json:"-" db:"not_in_json" gorm:"column:not_in_json"`
}

func (t TestObject) TableName() string {
	return "objects"
}

func (t TestObject) GetCustomFilters(ctx context.Context) scope.CustomColumns {
	customFilter := scope.CustomColumn{
		Name:       "custom_filter",
		Statement:  `(SELECT '1234'::TEXT)`,
		ResultType: reflect.TypeOf("1234"),
	}
	customUUIDFilter := scope.CustomColumn{
		Name:       "custom_uuid_filter",
		Statement:  `(SELECT uuid_nil())`,
		ResultType: reflect.TypeOf(uuid.UUID{}),
	}
	customNullsUUIDFilter := scope.CustomColumn{
		Name:       "custom_nulls_uuid_filter",
		Statement:  `(SELECT uuid_nil())`,
		ResultType: reflect.TypeOf(nulls.UUID{}),
	}
	nullFilter := scope.CustomColumn{
		Name:       "null_filter",
		Statement:  `NULL`,
		ResultType: reflect.TypeOf(nulls.UUID{}),
	}
	return scope.CustomColumns{customFilter, customUUIDFilter, customNullsUUIDFilter, nullFilter}
}

func (t TestObject) GetCustomSorts(ctx context.Context) scope.CustomColumns {
	customSort := scope.CustomColumn{
		Name:       "custom_sort",
		Statement:  `(SELECT '1234'::TEXT)`,
		ResultType: reflect.TypeOf("1234"),
	}
	customUUIDSort := scope.CustomColumn{
		Name:       "custom_uuid_sort",
		Statement:  `(SELECT uuid_nil())`,
		ResultType: reflect.TypeOf(uuid.UUID{}),
	}
	customNullsUUIDSort := scope.CustomColumn{
		Name:       "custom_nulls_uuid_sort",
		Statement:  `(SELECT uuid_nil())`,
		ResultType: reflect.TypeOf(nulls.UUID{}),
	}
	nullSort := scope.CustomColumn{
		Name:       "null_sort",
		Statement:  `NULL`,
		ResultType: reflect.TypeOf(nulls.UUID{}),
	}
	return scope.CustomColumns{customSort, customUUIDSort, customNullsUUIDSort, nullSort}
}

func (ss *ScopesSuite) TestGetFilterOptions() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	filterOptions, err := scope.GetFilterOptions(context.Background(), ss.DB, &[]TestObject{}, "id", nil)
	ss.NoError(err)
	ss.ElementsMatch([]interface{}{testObject.ID, testObject2.ID}, filterOptions)
}

func (ss *ScopesSuite) TestGetFilterOptions_noNulls() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	filterOptions, err := scope.GetFilterOptions(context.Background(), ss.DB, &[]TestObject{}, "null_id", nil)
	ss.NoError(err)
	ss.ElementsMatch([]interface{}{}, filterOptions)
}

func (ss *ScopesSuite) TestGetFilterOptions_withScopes() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	scopes := scope.NewCollection(ss.DB)
	scopes.Push(scope.ForUuidID(testObject.ID))
	scopes.Push(scope.ForOrder("id ASC"))
	scopes.Push(scope.ForLimit(2))

	filterOptions, err := scope.GetFilterOptions(context.Background(), ss.DB, &[]TestObject{}, "id", scopes)
	ss.NoError(err)
	ss.Equal([]interface{}{testObject.ID}, filterOptions)
}

func (ss *ScopesSuite) TestGetFilterOptions_withCustomFilters() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	scopes := scope.NewCollection(ss.DB)
	scopes.Push(scope.ForUuidID(testObject.ID))
	scopes.Push(scope.ForOrder("id ASC"))
	scopes.Push(scope.ForLimit(2))

	filterOptions, err := scope.GetFilterOptions(context.Background(), ss.DB, &[]TestObject{}, "custom_filter", scopes)
	ss.NoError(err)
	ss.Equal([]interface{}{"1234"}, filterOptions)
}

func (ss *ScopesSuite) TestGetFilterOptions_withCustomUUIDFilters() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	scopes := scope.NewCollection(ss.DB)
	scopes.Push(scope.ForUuidID(testObject.ID))
	scopes.Push(scope.ForOrder("id ASC"))
	scopes.Push(scope.ForLimit(2))

	filterOptions, err := scope.GetFilterOptions(context.Background(), ss.DB, &[]TestObject{}, "custom_uuid_filter", scopes)
	ss.NoError(err)
	ss.Equal([]interface{}{uuid.Nil}, filterOptions)
}

func (ss *ScopesSuite) TestGetFilterOptions_withCustomNullsUUIDFilters() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	scopes := scope.NewCollection(ss.DB)
	scopes.Push(scope.ForUuidID(testObject.ID))
	scopes.Push(scope.ForOrder("id ASC"))
	scopes.Push(scope.ForLimit(2))

	filterOptions, err := scope.GetFilterOptions(context.Background(), ss.DB, &[]TestObject{}, "custom_nulls_uuid_filter", scopes)
	ss.NoError(err)
	ss.Equal([]interface{}{nulls.NewUUID(uuid.Nil)}, filterOptions)
}

func (ss *ScopesSuite) TestGetFilterOptions_withCustomNullResult() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	scopes := scope.NewCollection(ss.DB)
	scopes.Push(scope.ForUuidID(testObject.ID))
	scopes.Push(scope.ForOrder("id ASC"))
	scopes.Push(scope.ForLimit(2))

	filterOptions, err := scope.GetFilterOptions(context.Background(), ss.DB, &[]TestObject{}, "null_filter", scopes)
	ss.NoError(err)
	ss.Equal([]interface{}{}, filterOptions)
}

type TestObjectPtrTablename struct {
	ID        uuid.UUID  `json:"id" db:"id" gorm:"primaryKey,column:id"`
	Nuid      nulls.UUID `json:"null_id" db:"db_null_id" gorm:"column:db_null_id"`
	NotInDb   int        `json:"not_in_db" db:"-" gorm:"-"`
	NotInJson int        `json:"-" db:"not_in_json" gorm:"column:not_in_json"`
}

type TestObjectPtrTablenames []TestObjectPtrTablename

func (t *TestObjectPtrTablename) TableName() string {
	return "objects"
}

func (t *TestObjectPtrTablenames) TableName() string {
	return "objects"
}

func (t TestObjectPtrTablename) GetCustomFilters(ctx context.Context) scope.CustomColumns {
	customFilter := scope.CustomColumn{
		Name:       "custom_filter",
		Statement:  `(SELECT '1234'::TEXT)`,
		ResultType: reflect.TypeOf("1234"),
	}
	customUUIDFilter := scope.CustomColumn{
		Name:       "custom_uuid_filter",
		Statement:  `(SELECT uuid_nil())`,
		ResultType: reflect.TypeOf(uuid.UUID{}),
	}
	customNullsUUIDFilter := scope.CustomColumn{
		Name:       "custom_nulls_uuid_filter",
		Statement:  `(SELECT uuid_nil())`,
		ResultType: reflect.TypeOf(nulls.UUID{}),
	}
	nullFilter := scope.CustomColumn{
		Name:       "null_filter",
		Statement:  `NULL`,
		ResultType: reflect.TypeOf(nulls.UUID{}),
	}
	return scope.CustomColumns{customFilter, customUUIDFilter, customNullsUUIDFilter, nullFilter}
}

func (t TestObjectPtrTablename) GetCustomSorts(ctx context.Context) scope.CustomColumns {
	customSort := scope.CustomColumn{
		Name:       "custom_sort",
		Statement:  `(SELECT '1234'::TEXT)`,
		ResultType: reflect.TypeOf("1234"),
	}
	customUUIDSort := scope.CustomColumn{
		Name:       "custom_uuid_sort",
		Statement:  `(SELECT uuid_nil())`,
		ResultType: reflect.TypeOf(uuid.UUID{}),
	}
	customNullsUUIDSort := scope.CustomColumn{
		Name:       "custom_nulls_uuid_sort",
		Statement:  `(SELECT uuid_nil())`,
		ResultType: reflect.TypeOf(nulls.UUID{}),
	}
	nullSort := scope.CustomColumn{
		Name:       "null_sort",
		Statement:  `NULL`,
		ResultType: reflect.TypeOf(nulls.UUID{}),
	}
	return scope.CustomColumns{customSort, customUUIDSort, customNullsUUIDSort, nullSort}
}

func (ss *ScopesSuite) TestGetFilterOptions_ptrTablename() {
	testObject := &TestObjectPtrTablename{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObjectPtrTablename{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	filterOptions, err := scope.GetFilterOptions(context.Background(), ss.DB, &TestObjectPtrTablenames{}, "id", nil)
	ss.NoError(err)
	ss.ElementsMatch([]interface{}{testObject.ID, testObject2.ID}, filterOptions)
}

func (ss *ScopesSuite) TestGetFilterOptions_noNulls_ptrTablename() {
	testObject := &TestObjectPtrTablename{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObjectPtrTablename{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	filterOptions, err := scope.GetFilterOptions(context.Background(), ss.DB, &TestObjectPtrTablenames{}, "null_id", nil)
	ss.NoError(err)
	ss.ElementsMatch([]interface{}{}, filterOptions)
}

func (ss *ScopesSuite) TestGetFilterOptions_withScopes_ptrTablename() {
	testObject := &TestObjectPtrTablename{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObjectPtrTablename{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	scopes := scope.NewCollection(ss.DB)
	scopes.Push(scope.ForUuidID(testObject.ID))
	scopes.Push(scope.ForOrder("id ASC"))
	scopes.Push(scope.ForLimit(2))

	filterOptions, err := scope.GetFilterOptions(context.Background(), ss.DB, &TestObjectPtrTablenames{}, "id", scopes)
	ss.NoError(err)
	ss.Equal([]interface{}{testObject.ID}, filterOptions)
}

func (ss *ScopesSuite) TestGetFilterOptions_withCustomFilters_ptrTablename() {
	testObject := &TestObjectPtrTablename{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObjectPtrTablename{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	scopes := scope.NewCollection(ss.DB)
	scopes.Push(scope.ForUuidID(testObject.ID))
	scopes.Push(scope.ForOrder("id ASC"))
	scopes.Push(scope.ForLimit(2))

	filterOptions, err := scope.GetFilterOptions(context.Background(), ss.DB, &TestObjectPtrTablenames{}, "custom_filter", scopes)
	ss.NoError(err)
	ss.Equal([]interface{}{"1234"}, filterOptions)
}

func (ss *ScopesSuite) TestGetFilterOptions_withCustomUUIDFilters_ptrTablename() {
	testObject := &TestObjectPtrTablename{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObjectPtrTablename{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	scopes := scope.NewCollection(ss.DB)
	scopes.Push(scope.ForUuidID(testObject.ID))
	scopes.Push(scope.ForOrder("id ASC"))
	scopes.Push(scope.ForLimit(2))

	filterOptions, err := scope.GetFilterOptions(context.Background(), ss.DB, &TestObjectPtrTablenames{}, "custom_uuid_filter", scopes)
	ss.NoError(err)
	ss.Equal([]interface{}{uuid.Nil}, filterOptions)
}

func (ss *ScopesSuite) TestGetFilterOptions_withCustomNullsUUIDFilters_ptrTablename() {
	testObject := &TestObjectPtrTablename{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObjectPtrTablename{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	scopes := scope.NewCollection(ss.DB)
	scopes.Push(scope.ForUuidID(testObject.ID))
	scopes.Push(scope.ForOrder("id ASC"))
	scopes.Push(scope.ForLimit(2))

	filterOptions, err := scope.GetFilterOptions(context.Background(), ss.DB, &TestObjectPtrTablenames{}, "custom_nulls_uuid_filter", scopes)
	ss.NoError(err)
	ss.Equal([]interface{}{nulls.NewUUID(uuid.Nil)}, filterOptions)
}

func (ss *ScopesSuite) TestGetFilterOptions_withCustomNullResult_ptrTablename() {
	testObject := &TestObjectPtrTablename{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObjectPtrTablename{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	scopes := scope.NewCollection(ss.DB)
	scopes.Push(scope.ForUuidID(testObject.ID))
	scopes.Push(scope.ForOrder("id ASC"))
	scopes.Push(scope.ForLimit(2))

	filterOptions, err := scope.GetFilterOptions(context.Background(), ss.DB, &TestObjectPtrTablenames{}, "null_filter", scopes)
	ss.NoError(err)
	ss.Equal([]interface{}{}, filterOptions)
}
