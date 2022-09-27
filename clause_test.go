package sqb_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqb "github.com/themanciraptor/SQb"
)

func Test_CanBuildPrimitiveFilterClause(t *testing.T) {
	p := sqb.NewParamList(sqb.Psql())
	expectedParams := []interface{}{
		42,
	}

	clause := sqb.NewPrimitiveFilterClause("cool", ">", "", 42)
	builtClause := clause.Build(p)

	assert.Equal(t, "cool > $1", builtClause)
	assert.Equal(t, expectedParams, p.GetParamList())
}

func Test_CanBuildPrimitiveFilterClauseWithCustomParamTemplate(t *testing.T) {
	p := sqb.NewParamList(sqb.Psql())
	expectedParams := []interface{}{
		time.Time{},
	}

	clause := sqb.NewPrimitiveFilterClause("week", ">", "EXTRACT(WEEK FROM %s)", time.Time{})
	builtClause := clause.Build(p)

	assert.Equal(t, "week > EXTRACT(WEEK FROM $1)", builtClause)
	assert.Equal(t, expectedParams, p.GetParamList())
}

func Test_CanBuildCompoundClause(t *testing.T) {
	type testCase struct {
		description string
		*sqb.CompoundClause
		expectedParams []interface{}
		expectedClause string
	}

	testCases := []testCase{
		{
			description:    "Can build empty clause",
			CompoundClause: sqb.NewCompoundClause("AND"),
			expectedParams: []interface{}{},
			expectedClause: "",
		},
		{
			description: "Can build with single clause",
			CompoundClause: sqb.NewCompoundClause("AND").
				AddClause(sqb.NewPrimitiveFilterClause("cool", ">", "", 42)),
			expectedParams: []interface{}{42},
			expectedClause: "cool > $1",
		},
		{
			description: "Can build with multiple primitive clauses",
			CompoundClause: sqb.NewCompoundClause("AND").
				AddClause(sqb.NewPrimitiveFilterClause("cool", ">", "", 42)).
				AddClause(sqb.NewPrimitiveFilterClause("we", "LIKE", "UCASE(%s)", "bleh")),
			expectedParams: []interface{}{42, "bleh"},
			expectedClause: "(cool > $1 AND we LIKE UCASE($2))",
		},
		{
			description: "Can build with multiple Compound clauses",
			CompoundClause: sqb.NewCompoundClause("AND").
				AddClause(sqb.NewPrimitiveFilterClause("cool", ">", "", 42)).
				AddClause(sqb.NewPrimitiveFilterClause("we", "LIKE", "UCASE(%s)", "bleh")).
				AddClause(
					sqb.NewCompoundClause("OR").
						AddClause(sqb.NewPrimitiveFilterClause("time", "<", "", time.Time{})).
						AddClause(sqb.NewPrimitiveFilterClause("name", "=", "", "Dovahkiin")),
				),
			expectedParams: []interface{}{42, "bleh", time.Time{}, "Dovahkiin"},
			expectedClause: "(cool > $1 AND we LIKE UCASE($2) AND (time < $3 OR name = $4))",
		},
	}

	for _, c := range testCases {
		t.Run(c.description, func(t *testing.T) {
			params := sqb.NewParamList(sqb.Psql())
			actualClause := c.CompoundClause.Build(params)

			assert.Equal(t, c.expectedClause, actualClause)
			assert.Equal(t, c.expectedParams, c.expectedParams)
		})
	}
}
