package sqb

import (
	"context"

	"github.com/vendasta/gosdks/verrors"
)

// In order to properly assign the result, in a generic way without requiring the
// developer to perform type checks. Developers should create a separate accumulator
// object to accumulate results. See common_test.go for an example accumulator
// implementation
type Accumulator interface {
	Acc()
	GetColumnReceiverMap() map[string]interface{}
}

type Query struct {
	scanList []interface{}
	query    string
	params   []interface{}

	accumulator Accumulator
}

func (q *Query) GetScanList() []interface{} {
	return q.scanList
}

func (q *Query) GetQuery() string {
	return q.query
}

func (q *Query) GetParams() []interface{} {
	return q.params
}

// tempRows is an interface that allows us to mock the rows object for testing purposes
type tempRows interface {
	Next() bool
	Scan(dest ...interface{}) error
}

type tempDriver interface {
	RunQuery(ctx context.Context, query string, params []interface{}) (res tempRows, closer func(ctx context.Context), err error)
}

func (q *Query) Run(ctx context.Context, psql tempDriver) error {
	res, closer, err := psql.RunQuery(ctx, q.query, q.params)
	if err != nil {
		return verrors.WrapError(err, "Failed to run query")
	}
	defer closer(ctx)

	for res.Next() {
		err := res.Scan(q.scanList...)
		if err != nil {
			return verrors.WrapError(err, "Failed to convert row")
		}

		q.accumulator.Acc()
	}
	return nil
}
