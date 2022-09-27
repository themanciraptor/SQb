package sqb_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	sqb "github.com/themanciraptor/SQb"
)

type exampleModel struct {
	Name       string    `psql:"cool"`
	Created    time.Time `psql:"created_time"`
	NumFoods   int32     `psql:"number_of_food"`
	NumStars   int64     `psql:"number_of_star"`
	MoonRadius float64   `psql:"radius_of_moon"`
	IsTrue     bool      `psql:"is_true_true"`
	Loves      []string  `psql:"loves"`
}

type exampleResult struct {
	Name       string
	Created    time.Time
	NumFoods   int32
	NumStars   int64
	MoonRadius float64
	IsTrue     bool
	Loves      []string
}

/*
Used by the queryBuilder to assign values to unexported fields.
*/
type exampleResultAccumulator struct {
	// Map columns to their result receivers, this is only necessary because
	// the receiver may be a different type from it's final value in the result
	receiverMap map[string]interface{}

	// Accumulated Results
	results []exampleResult

	// Receiver struct
	receiver *exampleResult
}

func (r *exampleResultAccumulator) Acc() {
	r.results = append(r.results, *r.receiver)
}

func (r *exampleResultAccumulator) GetColumnReceiverMap() map[string]interface{} {
	return r.receiverMap
}

func (r *exampleResultAccumulator) GetResults() []exampleResult {
	return r.results
}

func NewResultAccumulator() *exampleResultAccumulator {
	r := new(exampleResult)
	return &exampleResultAccumulator{
		results:  []exampleResult{},
		receiver: r,
		receiverMap: map[string]interface{}{
			"cool":           &r.Name,
			"created_time":   &r.Created,
			"number_of_food": sqb.NewNullInt32(&r.NumFoods),
			"number_of_star": sqb.NewNullInt64(&r.NumStars),
			"radius_of_moon": sqb.NewNullFloat64(&r.MoonRadius),
			"is_true_true":   sqb.NewNullBool(&r.IsTrue),
			"loves":          &r.Loves,
		},
	}
}

/*
The following tests are for the example accumulator. They may also
serve as an example of how to test an accumulator.
*/
func Test_ResultAccumulator_AccumulatesResults(t *testing.T) {
	a := NewResultAccumulator()

	r := a.GetColumnReceiverMap()

	expectedTime := time.Date(2011, 11, 11, 00, 0, 0, 0, time.UTC)

	*r["cool"].(*string) = "doom"
	*r["created_time"].(*time.Time) = expectedTime
	*r["number_of_food"].(*sqb.NullInt32).Int32 = 32
	*r["number_of_star"].(*sqb.NullInt64).Int64 = 64
	*r["radius_of_moon"].(*sqb.NullFloat64).Float64 = 64.64
	*r["is_true_true"].(*sqb.NullBool).Bool = true
	*r["loves"].(*[]string) = append(*r["loves"].(*[]string), "doom")

	a.Acc()

	expected := []exampleResult{
		{
			"doom",
			expectedTime,
			32,
			64,
			64.64,
			true,
			[]string{"doom"},
		},
	}

	actual := a.GetResults()

	assert.Equal(t, expected, actual)
}
