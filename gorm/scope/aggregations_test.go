package scope_test

import (
	"context"
	"fmt"
	"net/url"

	"github.com/alphaflow/scope/util"
	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"


	"github.com/alphaflow/scope/gorm/scope"
)

func (ss *ScopesSuite) TestGetAggregationsFromParams_Count() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column": {"id"},
		"aggregation_type":   {"count"},
	}

	aggregation, err := scope.GetAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.NoError(err)
	ss.Equal(2, aggregation)
}

func (ss *ScopesSuite) TestGetAggregationsFromParams_Count_scoped() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column": {"id"},
		"aggregation_type":   {"count"},
	}

	sc := scope.NewCollection(ss.DB)
	sc.Push(scope.ForID(testObject.ID.String()))

	aggregation, err := scope.GetAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), sc)
	ss.NoError(err)
	ss.Equal(1, aggregation)
}

// There is a known incompatibility between this package and Gorm.  See https://github.com/go-gorm/gorm/issues/5170.
func (ss *ScopesSuite) TestGetAggregationsFromParams_Count_scopedAtSymbol() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column": {"id"},
		"aggregation_type":   {"count"},
	}

	sc := scope.NewCollection(ss.DB)

	atScopeFunc := func(q *gorm.DB) *gorm.DB {
		return q.Where("'@' = ?", "@")
	}

	sc.Push(atScopeFunc)

	_, err = scope.GetAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), sc)
	ss.Error(err)
}

func (ss *ScopesSuite) TestGetAggregationsFromParams_Count_scopedGenericFiltering() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column": {"id"},
		"aggregation_type":   {"count"},
		"filter_columns":     {"id"},
		"filter_types":       {"in"},
		"filter_values":      {fmt.Sprintf("%s,%s", testObject.ID.String(), testObject.ID.String())},
	}

	sc := scope.NewCollection(ss.DB)
	sc.Push(scope.ForID(testObject.ID.String()))
	s, err := scope.ForFiltersFromParams(context.Background(), TestObject{}, url.Values(params))
	ss.NoError(err)
	sc.Push(s)

	aggregation, err := scope.GetAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), sc)
	ss.NoError(err)
	ss.Equal(1, aggregation)
}

func (ss *ScopesSuite) TestGetAggregationsFromParams_Sum() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column": {"num"},
		"aggregation_type":   {"sum"},
	}

	aggregation, err := scope.GetAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.NoError(err)
	ss.Equal(float64(246), aggregation)
}

func (ss *ScopesSuite) TestGetAggregationsFromParams_Avg() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 125}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column": {"num"},
		"aggregation_type":   {"avg"},
	}

	aggregation, err := scope.GetAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.NoError(err)
	ss.Equal(float64(124), aggregation)
}

func (ss *ScopesSuite) TestGetAggregationsFromParams_Max() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 124}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column": {"num"},
		"aggregation_type":   {"max"},
	}

	aggregation, err := scope.GetAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.NoError(err)
	ss.Equal(float64(124), aggregation)
}

func (ss *ScopesSuite) TestGetAggregationsFromParams_Min() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 124}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column": {"num"},
		"aggregation_type":   {"min"},
	}

	aggregation, err := scope.GetAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.NoError(err)
	ss.Equal(float64(123), aggregation)
}

func (ss *ScopesSuite) TestGetAggregations_Count() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}

	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, "id", nil, scope.StandardAggregations[scope.StandardAggregationsTypeCount])
	ss.NoError(err)
	ss.Equal(2, aggregation)
}

func (ss *ScopesSuite) TestGetAggregations_Count_scoped() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	sc := scope.NewCollection(ss.DB)
	sc.Push(scope.ForID(testObject.ID.String()))

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, "id", sc, scope.StandardAggregations[scope.StandardAggregationsTypeCount])
	ss.NoError(err)
	ss.Equal(1, aggregation)
}

func (ss *ScopesSuite) TestGetAggregations_Sum() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, "num", nil, scope.StandardAggregations[scope.StandardAggregationsTypeSum])
	ss.NoError(err)
	ss.Equal(float64(246), aggregation)
}

func (ss *ScopesSuite) TestGetAggregations_Avg() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 125}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, "num", nil, scope.StandardAggregations[scope.StandardAggregationsTypeAvg])
	ss.NoError(err)
	ss.Equal(float64(124), aggregation)
}

func (ss *ScopesSuite) TestGetAggregations_Max() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 124}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, "num", nil, scope.StandardAggregations[scope.StandardAggregationsTypeMax])
	ss.NoError(err)
	ss.Equal(float64(124), aggregation)
}

func (ss *ScopesSuite) TestGetAggregations_Min() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 124}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, "num", nil, scope.StandardAggregations[scope.StandardAggregationsTypeMin])
	ss.NoError(err)
	ss.Equal(float64(123), aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregationsFromParams_Count() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column":         {"id"},
		"aggregation_grouper_column": {"num"},
		"aggregation_type":           {"count"},
	}

	aggregation, err := scope.GetGroupedAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.NoError(err)

	ss.Equal([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result  int     "db:\"result\""
	}{Grouper: 0, Result: 2}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregationsFromParams_Count_grouped() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 1}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column":         {"id"},
		"aggregation_grouper_column": {"num"},
		"aggregation_type":           {"COUNT"},
	}

	aggregation, err := scope.GetGroupedAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.NoError(err)
	ss.ElementsMatch([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result  int     "db:\"result\""
	}{Grouper: 0, Result: 1}, struct {
		Grouper float64 "db:\"grouper\""
		Result  int     "db:\"result\""
	}{Grouper: 1, Result: 1}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregationsFromParams_Sum() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	testObject3 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 124}
	err = ss.DB.Create(testObject3).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column":         {"num"},
		"aggregation_grouper_column": {"num"},
		"aggregation_type":           {"sum"},
	}

	aggregation, err := scope.GetGroupedAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.NoError(err)
	ss.ElementsMatch([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result  float64 "db:\"result\""
	}{Grouper: 123, Result: 246}, struct {
		Grouper float64 "db:\"grouper\""
		Result  float64 "db:\"result\""
	}{Grouper: 124, Result: 124}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregationsFromParams_Avg() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 125}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column":         {"num"},
		"aggregation_grouper_column": {"num"},
		"aggregation_type":           {"avg"},
	}

	aggregation, err := scope.GetGroupedAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.NoError(err)
	ss.ElementsMatch([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result  float64 "db:\"result\""
	}{Grouper: 123, Result: 123}, struct {
		Grouper float64 "db:\"grouper\""
		Result  float64 "db:\"result\""
	}{Grouper: 125, Result: 125}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregationsFromParams_Max() {
	nuid := util.UuidMust()
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Nuid: nuid, Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Nuid: nuid, Number: 124}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column":         {"num"},
		"aggregation_grouper_column": {"null_id"},
		"aggregation_type":           {"max"},
	}

	aggregation, err := scope.GetGroupedAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.NoError(err)
	ss.Equal([]interface{}{struct {
		Grouper nulls.UUID "db:\"grouper\""
		Result  float64    "db:\"result\""
	}{Grouper: nuid, Result: 124}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregationsFromParams_Min() {
	nuid := util.UuidMust()
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Nuid: nuid, Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Nuid: nuid, Number: 124}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column":         {"num"},
		"aggregation_grouper_column": {"null_id"},
		"aggregation_type":           {"mIn"},
	}

	aggregation, err := scope.GetGroupedAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.NoError(err)
	ss.Equal([]interface{}{struct {
		Grouper nulls.UUID "db:\"grouper\""
		Result  float64    "db:\"result\""
	}{Grouper: nuid, Result: 123}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregations_Count() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, "id", "num", nil, scope.StandardAggregations[scope.StandardAggregationsTypeCount])
	ss.NoError(err)
	ss.Equal([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result  int     "db:\"result\""
	}{Grouper: 0, Result: 2}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregations_Count_scoped() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	sc := scope.NewCollection(ss.DB)
	sc.Push(scope.ForID(testObject.ID.String()))

	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, "id", "num", sc, scope.StandardAggregations[scope.StandardAggregationsTypeCount])
	ss.NoError(err)
	ss.Equal([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result  int     "db:\"result\""
	}{Grouper: 0, Result: 1}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregations_Sum() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	testObject3 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 124}
	err = ss.DB.Create(testObject3).Error
	ss.NoError(err)

	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, "num", "num", nil, scope.StandardAggregations[scope.StandardAggregationsTypeSum])
	ss.NoError(err)
	ss.ElementsMatch([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result  float64 "db:\"result\""
	}{Grouper: 123, Result: 246}, struct {
		Grouper float64 "db:\"grouper\""
		Result  float64 "db:\"result\""
	}{Grouper: 124, Result: 124}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregations_Avg() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 125}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, "num", "num", nil, scope.StandardAggregations[scope.StandardAggregationsTypeAvg])
	ss.NoError(err)

	ss.ElementsMatch([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result  float64 "db:\"result\""
	}{Grouper: 123, Result: 123}, struct {
		Grouper float64 "db:\"grouper\""
		Result  float64 "db:\"result\""
	}{Grouper: 125, Result: 125}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregations_Max() {
	nuid := util.UuidMust()
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Nuid: nuid, Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Nuid: nuid, Number: 124}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, "num", "null_id", nil, scope.StandardAggregations[scope.StandardAggregationsTypeMax])
	ss.NoError(err)
	ss.Equal([]interface{}{struct {
		Grouper nulls.UUID "db:\"grouper\""
		Result  float64    "db:\"result\""
	}{Grouper: nuid, Result: 124}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregations_Min() {
	nuid := util.UuidMust()
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Nuid: nuid, Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Nuid: nuid, Number: 124}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, "num", "null_id", nil, scope.StandardAggregations[scope.StandardAggregationsTypeMin])
	ss.NoError(err)
	ss.Equal([]interface{}{struct {
		Grouper nulls.UUID "db:\"grouper\""
		Result  float64    "db:\"result\""
	}{Grouper: nuid, Result: 123}}, aggregation)
}
