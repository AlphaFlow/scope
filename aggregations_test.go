package scope_test

import (
	"context"
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
	ss.Equal(2, aggregation)
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
	ss.Equal(1, aggregation)
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
	ss.Equal(float64(246), aggregation)
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
	ss.Equal(float64(124), aggregation)
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
	ss.Equal(float64(124), aggregation)
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
	ss.Equal(float64(123), aggregation)
}

func (ss *ScopesSuite) TestGetAggregations_Count() {
	testObject := &TestObject{}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, "id", nil, scope.StandardAggregations[scope.StandardAggregationsTypeCount])
	ss.NoError(err)
	ss.Equal(2, aggregation)
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

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, "id", sc, scope.StandardAggregations[scope.StandardAggregationsTypeCount])
	ss.NoError(err)
	ss.Equal(1, aggregation)
}

func (ss *ScopesSuite) TestGetAggregations_Sum() {
	testObject := &TestObject{Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 123}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, "num", nil, scope.StandardAggregations[scope.StandardAggregationsTypeSum])
	ss.NoError(err)
	ss.Equal(float64(246), aggregation)
}

func (ss *ScopesSuite) TestGetAggregations_Avg() {
	testObject := &TestObject{Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 125}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, "num", nil, scope.StandardAggregations[scope.StandardAggregationsTypeAvg])
	ss.NoError(err)
	ss.Equal(float64(124), aggregation)
}

func (ss *ScopesSuite) TestGetAggregations_Max() {
	testObject := &TestObject{Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 124}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, "num", nil, scope.StandardAggregations[scope.StandardAggregationsTypeMax])
	ss.NoError(err)
	ss.Equal(float64(124), aggregation)
}

func (ss *ScopesSuite) TestGetAggregations_Min() {
	testObject := &TestObject{Number: 124}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 123}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	aggregation, err := scope.GetAggregations(context.Background(), ss.DB, &[]TestObject{}, "num", nil, scope.StandardAggregations[scope.StandardAggregationsTypeMin])
	ss.NoError(err)
	ss.Equal(float64(123), aggregation)
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
		Result  int     "db:\"result\""
	}{Grouper: 0, Result: 2}}, aggregation)
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
		Result  int     "db:\"result\""
	}{Grouper: 0, Result: 1}, struct {
		Grouper float64 "db:\"grouper\""
		Result  int     "db:\"result\""
	}{Grouper: 1, Result: 1}}, aggregation)
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
		Result  float64 "db:\"result\""
	}{Grouper: 123, Result: 246}, struct {
		Grouper float64 "db:\"grouper\""
		Result  float64 "db:\"result\""
	}{Grouper: 124, Result: 124}}, aggregation)
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
		Result  float64 "db:\"result\""
	}{Grouper: 123, Result: 123}, struct {
		Grouper float64 "db:\"grouper\""
		Result  float64 "db:\"result\""
	}{Grouper: 125, Result: 125}}, aggregation)
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
		Result  float64    "db:\"result\""
	}{Grouper: nuid, Result: 124}}, aggregation)
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
		Result  float64    "db:\"result\""
	}{Grouper: nuid, Result: 123}}, aggregation)
}

func (ss *ScopesSuite) TestGetGroupedAggregations_Count() {
	testObject := &TestObject{}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, "id", "num", nil, scope.StandardAggregations[scope.StandardAggregationsTypeCount])
	ss.NoError(err)
	ss.Equal([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result  int     "db:\"result\""
	}{Grouper: 0, Result: 2}}, aggregation)
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

	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, "id", "num", sc, scope.StandardAggregations[scope.StandardAggregationsTypeCount])
	ss.NoError(err)
	ss.Equal([]interface{}{struct {
		Grouper float64 "db:\"grouper\""
		Result  int     "db:\"result\""
	}{Grouper: 0, Result: 1}}, aggregation)
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
	testObject := &TestObject{Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Number: 125}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{Nuid: nuid, Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Nuid: nuid, Number: 124}
	err = ss.DB.Create(testObject2)
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
	testObject := &TestObject{Nuid: nuid, Number: 123}
	err := ss.DB.Create(testObject)
	ss.NoError(err)

	testObject2 := &TestObject{Nuid: nuid, Number: 124}
	err = ss.DB.Create(testObject2)
	ss.NoError(err)

	aggregation, err := scope.GetGroupedAggregations(context.Background(), ss.DB, &[]TestObject{}, "num", "null_id", nil, scope.StandardAggregations[scope.StandardAggregationsTypeMin])
	ss.NoError(err)
	ss.Equal([]interface{}{struct {
		Grouper nulls.UUID "db:\"grouper\""
		Result  float64    "db:\"result\""
	}{Grouper: nuid, Result: 123}}, aggregation)
}
