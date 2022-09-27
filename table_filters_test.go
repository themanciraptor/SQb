package sqb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	sqb "github.com/themanciraptor/SQb"
)

func Test_PrimitiveFilters_BuildCorrectly(t *testing.T) {
	type TableFilterBuilder = func(tt *sqb.Table)

	type testCase struct {
		description    string
		expectedClause string
		TableFilterBuilder
	}

	testCases := []testCase{
		{
			description:    "equals filter",
			expectedClause: "cool = $1",
			TableFilterBuilder: func(tt *sqb.Table) {
				tt.ColumnEquals("cool", "don't care")
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			tt := sqb.NewTable("exampleTable", sqb.Psql(), &exampleModel{})
			params := sqb.NewParamList(sqb.Psql())

			tc.TableFilterBuilder(tt)

			assert.Equal(t, tc.expectedClause, tt.BuildFilter(params))
		})
	}
}

func Test_IsNull_BuildsCorrectly(t *testing.T) {
	params := sqb.NewParamList(sqb.Psql())
	tt := sqb.NewTable("exampleTable", sqb.Psql(), &exampleModel{}).ColumnNull("cool")

	expected := "cool IS NULL"
	actual := tt.BuildFilter(params)

	expectedParams := []interface{}{}

	assert.Equal(t, expected, actual)
	assert.Equal(t, expectedParams, params.GetParamList())
}

func Test_LimitClause_BuildsCorrectly(t *testing.T) {
	params := sqb.NewParamList(sqb.Psql())
	tt := sqb.NewTable("exampleTable", sqb.Psql(), &exampleModel{}).Limit(25, 5)

	expected := "LIMIT 25 OFFSET 5"
	actual := tt.Build(&exampleResultAccumulator{}, sqb.Psql())

	expectedParams := []interface{}{}

	assert.Contains(t, actual.GetQuery(), expected)
	assert.Equal(t, expectedParams, params.GetParamList())
}
