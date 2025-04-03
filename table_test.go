package sqb_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqb "github.com/themanciraptor/SQb"
)

func Test_NewTable_Panics(t *testing.T) {
	type testCase struct {
		description string
		tableModel  interface{}
		panicMsg    string
	}

	testCases := []testCase{
		{
			description: "when a non-struct pointer passed into NewTable constructor",
			tableModel:  new(string),
			panicMsg:    "QueryBuilder: Table model must be pointer to struct type",
		},
		{
			description: "when non-pointer object passed into NewTable constructor",
			tableModel:  exampleModel{},
			panicMsg:    "QueryBuilder: Table model must be pointer to struct type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			assert.PanicsWithValue(t, tc.panicMsg, func() { sqb.NewTable[exampleResult]("example", sqb.Psql(), tc.tableModel) })
		})
	}
}

func Test_SetColumnReceiver_Panics(t *testing.T) {
	refModel := &exampleModel{}
	tt := sqb.NewTable[exampleResult]("exampleTable", sqb.Psql(), refModel)

	type testPanicMsg func()

	type testCase struct {
		description string
		columnName  string
		scanTo      interface{}
		testPanicMsg
	}

	testCases := []testCase{
		{
			description: "when an address to a struct field not included in the reference model passed in",
			columnName:  "created_tim",
			scanTo:      new(time.Time),
			testPanicMsg: func() {
				r := recover()
				assert.Contains(t, r, "SetField: column not included in table")
				panic(r)
			},
		},
		{
			description: "when non-pointer reference passed in for scanTo",
			columnName:  "created_time",
			scanTo:      time.Time{},
			testPanicMsg: func() {
				r := recover()
				assert.Contains(t, r, "SetField: scanTo must be reference pointer")
				panic(r)
			},
		},
		{
			description: "when receiver and referenceModel field are not the same type",
			columnName:  "cool",
			scanTo:      new(time.Time),
			testPanicMsg: func() {
				r := recover()
				assert.Contains(t, r, "SetField: attempted to scan to invalid type for field: cannot assign")
				panic(r)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			assert.Panics(t, func() {
				defer tc.testPanicMsg()
				tt.SetColumnReceiver(tc.columnName, tc.scanTo)
			})
		})
	}
}

func Test_CanBuildSimpleSelectQuery(t *testing.T) {
	r := exampleResult{}
	acc := exampleResultAccumulator{}

	// alternate order's are possible because we are iterating over a map
	expectedQuery := "SELECT cool, created_time FROM exampleTable"
	expectedScanList := []interface{}{&r.Name, &r.Created}

	e := exampleModel{}
	actualQuery := sqb.NewTable[exampleResult]("exampleTable", sqb.Psql(), &e).
		SetColumnReceiver("cool", &r.Name).
		SetColumnReceiver("created_time", &r.Created).
		Build(&acc, sqb.Psql())

	assert.Equal(t, actualQuery.GetQuery(), expectedQuery)
	if actualQuery.GetQuery() == expectedQuery {
		assert.Equal(t, expectedScanList, actualQuery.GetScanList())
	}
}

func Test_TableWithQueryFiltersIsBuiltProperly(t *testing.T) {
	r := exampleResult{}
	acc := exampleResultAccumulator{}

	// alternate order's are possible because we are iterating over a map
	expectedQuery := "SELECT cool, loves FROM exampleTable WHERE cool IS NULL"

	e := exampleModel{}
	actualQuery := sqb.NewTable[exampleResult]("exampleTable", sqb.Psql(), &e).
		SetColumnReceiver("cool", &r.Name).
		SetColumnReceiver("loves", &e.Loves).
		ColumnNull("cool").
		Build(&acc, sqb.Psql())

	assert.Equal(t, expectedQuery, actualQuery.GetQuery())
}

func Test_CanAddReceiversFromReceiverMap(t *testing.T) {
	acc := NewResultAccumulator()

	e := exampleModel{}
	assert.NotPanics(t, func() {
		sqb.NewTable[exampleResult]("exampleTable", sqb.Psql(), &e).LoadReceiversFromAccumulator(acc)
	})
}

func Test_ReceiverMapPanics(t *testing.T) {
	e := exampleModel{}

	type testPanicMsg func()

	type testCase struct {
		description string
		columnName  string
		receiver    interface{}
		testPanicMsg
	}

	testCases := []testCase{
		{
			description: "when a column name that doesn't exist in the table is passed in",
			columnName:  "non_existent_column",
			receiver:    new(time.Time),
			testPanicMsg: func() {
				r := recover()
				assert.Contains(t, r, "non_existent_column: Column not included in table")
				panic(r)
			},
		},
		{
			description: "when non-pointer reference passed in for receiver",
			columnName:  "created_time",
			receiver:    time.Time{},
			testPanicMsg: func() {
				r := recover()
				assert.Contains(t, r, "created_time: receiver must be reference pointer")
				panic(r)
			},
		},
		{
			description: "when receiver and referenceModel field are not the same type",
			columnName:  "cool",
			receiver:    new(time.Time),
			testPanicMsg: func() {
				r := recover()
				assert.Contains(t, r, "cool: Invalid type for field: cannot assign")
				panic(r)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			assert.Panics(t, func() {
				defer tc.testPanicMsg()
				acc := &exampleResultAccumulator{
					receiverMap: map[string]interface{}{
						tc.columnName: tc.receiver,
					},
				}

				sqb.NewTable[exampleResult]("exampleTable", sqb.Psql(), &e).LoadReceiversFromAccumulator(acc)
			})
		})
	}
}
