package scope_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"
	"gorm.io/gorm"

	"github.com/alphaflow/scope/gorm/scope"
	"github.com/alphaflow/scope/util"
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
	ss.Equal(2, util.GetFieldByName(aggregation, "Result0").Interface())

	jsn, err := json.Marshal(aggregation)
	ss.NoError(err)
	ss.Equal(`{"count_id":2}`, string(jsn))
}

func (ss *ScopesSuite) TestGetAggregationsFromParams_Count_multiple() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column": {"id|num"},
		"aggregation_type":   {"count|count"},
	}

	aggregation, err := scope.GetAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.NoError(err)
	ss.Equal(2, util.GetFieldByName(aggregation, "Result0").Interface())
	ss.Equal(2, util.GetFieldByName(aggregation, "Result1").Interface())

	jsn, err := json.Marshal(aggregation)
	ss.NoError(err)
	ss.Equal(`{"count_id":2,"count_num":2}`, string(jsn))
}

func (ss *ScopesSuite) TestGetAggregationsFromParams_Count_duplicate() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column": {"id|id"},
		"aggregation_type":   {"count|count"},
	}

	_, err = scope.GetAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.Error(err)
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
	ss.Equal(1, util.GetFieldByName(aggregation, "Result0").Interface())
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
	ss.Equal(1, util.GetFieldByName(aggregation, "Result0").Interface())
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
	ss.Equal(float64(246), util.GetFieldByName(aggregation, "Result0").Interface())
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
	ss.Equal(float64(124), util.GetFieldByName(aggregation, "Result0").Interface())
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
	ss.Equal(float64(124), util.GetFieldByName(aggregation, "Result0").Interface())
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
	ss.Equal(float64(123), util.GetFieldByName(aggregation, "Result0").Interface())
}

func (ss *ScopesSuite) TestGetAggregations_Count() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"id"}, nil, scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeCount]})
	ss.NoError(err)
	ss.Equal(2, util.GetFieldByName(aggregation, "Result0").Interface())
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

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"id"}, sc, scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeCount]})
	ss.NoError(err)
	ss.Equal(1, util.GetFieldByName(aggregation, "Result0").Interface())
}

func (ss *ScopesSuite) TestGetAggregations_Sum() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"num"}, nil, scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeSum]})
	ss.NoError(err)
	ss.Equal(float64(246), util.GetFieldByName(aggregation, "Result0").Interface())
}

func (ss *ScopesSuite) TestGetAggregations_Avg() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 125}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"num"}, nil, scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeAvg]})
	ss.NoError(err)
	ss.Equal(float64(124), util.GetFieldByName(aggregation, "Result0").Interface())
}

func (ss *ScopesSuite) TestGetAggregations_Max() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 124}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"num"}, nil, scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeMax]})
	ss.NoError(err)
	ss.Equal(float64(124), util.GetFieldByName(aggregation, "Result0").Interface())
}

func (ss *ScopesSuite) TestGetAggregations_Min() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 124}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"num"}, nil, scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeMin]})
	ss.NoError(err)
	ss.Equal(float64(123), util.GetFieldByName(aggregation, "Result0").Interface())
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
		Result0 int     "db:\"result0\" json:\"count_id\""
	}{Grouper: 0, Result0: 2}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregationsFromParams_Count_multiple() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column":         {"id|num"},
		"aggregation_grouper_column": {"num"},
		"aggregation_type":           {"count|count"},
	}

	aggregation, err := scope.GetGroupedAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.NoError(err)

	ss.Equal([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result0 int     "db:\"result0\" json:\"count_id\""
		Result1 int     "db:\"result1\" json:\"count_num\""
	}{Grouper: 0, Result0: 2, Result1: 2}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregationsFromParams_Count_duplicate() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column":         {"id|id"},
		"aggregation_grouper_column": {"num"},
		"aggregation_type":           {"count|count"},
	}

	_, err = scope.GetGroupedAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.Error(err)
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
		Result0 int     "db:\"result0\" json:\"count_id\""
	}{Grouper: 0, Result0: 1}, struct {
		Grouper float64 "db:\"grouper\""
		Result0 int     "db:\"result0\" json:\"count_id\""
	}{Grouper: 1, Result0: 1}}, aggregation)
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
		Result0 float64 "db:\"result0\" json:\"sum_num\""
	}{Grouper: 123, Result0: 246}, struct {
		Grouper float64 "db:\"grouper\""
		Result0 float64 "db:\"result0\" json:\"sum_num\""
	}{Grouper: 124, Result0: 124}}, aggregation)
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
		Result0 float64 "db:\"result0\" json:\"avg_num\""
	}{Grouper: 123, Result0: 123}, struct {
		Grouper float64 "db:\"grouper\""
		Result0 float64 "db:\"result0\" json:\"avg_num\""
	}{Grouper: 125, Result0: 125}}, aggregation)
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
		Result0 float64    "db:\"result0\" json:\"max_num\""
	}{Grouper: nuid, Result0: 124}}, aggregation)
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
		Result0 float64    "db:\"result0\" json:\"min_num\""
	}{Grouper: nuid, Result0: 123}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregations_Count() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4())}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregations := scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeCount]}
	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"id"}, "num", nil, aggregations)
	ss.NoError(err)
	ss.Equal([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result0 int     "db:\"result0\" json:\"count_id\""
	}{Grouper: 0, Result0: 2}}, aggregation)
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

	aggregations := scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeCount]}
	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"id"}, "num", sc, aggregations)
	ss.NoError(err)
	ss.Equal([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result0 int     "db:\"result0\" json:\"count_id\""
	}{Grouper: 0, Result0: 1}}, aggregation)
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

	aggregations := scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeSum]}
	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"num"}, "num", nil, aggregations)
	ss.NoError(err)
	ss.ElementsMatch([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result0 float64 "db:\"result0\" json:\"sum_num\""
	}{Grouper: 123, Result0: 246}, struct {
		Grouper float64 "db:\"grouper\""
		Result0 float64 "db:\"result0\" json:\"sum_num\""
	}{Grouper: 124, Result0: 124}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregations_Avg() {
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Number: 125}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregations := scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeAvg]}
	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"num"}, "num", nil, aggregations)
	ss.NoError(err)

	ss.ElementsMatch([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result0 float64 "db:\"result0\" json:\"avg_num\""
	}{Grouper: 123, Result0: 123}, struct {
		Grouper float64 "db:\"grouper\""
		Result0 float64 "db:\"result0\" json:\"avg_num\""
	}{Grouper: 125, Result0: 125}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregations_Max() {
	nuid := util.UuidMust()
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Nuid: nuid, Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Nuid: nuid, Number: 124}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregations := scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeMax]}
	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"num"}, "null_id", nil, aggregations)
	ss.NoError(err)
	ss.Equal([]interface{}{struct {
		Grouper nulls.UUID "db:\"grouper\""
		Result0 float64    "db:\"result0\" json:\"max_num\""
	}{Grouper: nuid, Result0: 124}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregations_Min() {
	nuid := util.UuidMust()
	testObject := &TestObject{ID: uuid.Must(uuid.NewV4()), Nuid: nuid, Number: 123}
	err := ss.DB.Create(testObject).Error
	ss.NoError(err)

	testObject2 := &TestObject{ID: uuid.Must(uuid.NewV4()), Nuid: nuid, Number: 124}
	err = ss.DB.Create(testObject2).Error
	ss.NoError(err)

	aggregations := scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeMin]}
	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"num"}, "null_id", nil, aggregations)
	ss.NoError(err)
	ss.Equal([]interface{}{struct {
		Grouper nulls.UUID "db:\"grouper\""
		Result0 float64    "db:\"result0\" json:\"min_num\""
	}{Grouper: nuid, Result0: 123}}, aggregation)
}
