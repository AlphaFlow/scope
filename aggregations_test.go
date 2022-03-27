package scope_test

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/gobuffalo/nulls"

	"github.com/alphaflow/scope"
	"github.com/alphaflow/scope/util"
)

func (ss *ScopesSuite) TestGetAggregationsFromParams_Count() {
	testObject := &TestObject{}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	params := map[string][]string{
		"aggregation_column": {"id|id"},
		"aggregation_type":   {"count|count"},
	}

	_, err = scope.GetAggregationsFromParams(context.Background(), ss.DB, &[]TestObject{}, url.Values(params), nil)
	ss.Error(err)
}

func (ss *ScopesSuite) TestGetAggregationsFromParams_Count_scoped() {
	testObject := &TestObject{}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{}
	err = ss.DB.Create(testObject2)
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

func (ss *ScopesSuite) TestGetAggregationsFromParams_Sum() {
	testObject := &TestObject{Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 123}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 125}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 124}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{Number: 124}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 123}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"id"}, nil, scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeCount]})
	ss.NoError(err)
	ss.Equal(2, util.GetFieldByName(aggregation, "Result0").Interface())
}

func (ss *ScopesSuite) TestGetAggregations_Count_scoped() {
	testObject := &TestObject{}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	sc := scope.NewCollection(ss.DB)
	sc.Push(scope.ForID(testObject.ID.String()))

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"id"}, sc, scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeCount]})
	ss.NoError(err)
	ss.Equal(1, util.GetFieldByName(aggregation, "Result0").Interface())
}

func (ss *ScopesSuite) TestGetAggregations_Sum() {
	testObject := &TestObject{Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 123}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"num"}, nil, scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeSum]})
	ss.NoError(err)
	ss.Equal(float64(246), util.GetFieldByName(aggregation, "Result0").Interface())
}

func (ss *ScopesSuite) TestGetAggregations_Avg() {
	testObject := &TestObject{Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 125}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"num"}, nil, scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeAvg]})
	ss.NoError(err)
	ss.Equal(float64(124), util.GetFieldByName(aggregation, "Result0").Interface())
}

func (ss *ScopesSuite) TestGetAggregations_Max() {
	testObject := &TestObject{Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 124}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"num"}, nil, scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeMax]})
	ss.NoError(err)
	ss.Equal(float64(124), util.GetFieldByName(aggregation, "Result0").Interface())
}

func (ss *ScopesSuite) TestGetAggregations_Min() {
	testObject := &TestObject{Number: 124}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 123}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"num"}, nil, scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeMin]})
	ss.NoError(err)
	ss.Equal(float64(123), util.GetFieldByName(aggregation, "Result0").Interface())
}

func (ss *ScopesSuite) TestGetGroupedAggregationsFromParams_Count() {
	testObject := &TestObject{}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 1}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 123}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	testObject3 := &TestObject{Number: 124}
	err = ss.DB.Create(testObject3)
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
	testObject := &TestObject{Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 125}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{Nuid: nuid, Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Nuid: nuid, Number: 124}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{Nuid: nuid, Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Nuid: nuid, Number: 124}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 123}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	testObject3 := &TestObject{Number: 124}
	err = ss.DB.Create(testObject3)
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
	testObject := &TestObject{Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 125}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{Nuid: nuid, Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Nuid: nuid, Number: 124}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{Nuid: nuid, Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Nuid: nuid, Number: 124}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	aggregations := scope.Aggregations{scope.StandardAggregations[scope.StandardAggregationsTypeMin]}
	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, []string{"num"}, "null_id", nil, aggregations)
	ss.NoError(err)
	ss.Equal([]interface{}{struct {
		Grouper nulls.UUID "db:\"grouper\""
		Result0 float64    "db:\"result0\" json:\"min_num\""
	}{Grouper: nuid, Result0: 123}}, aggregation)
}
