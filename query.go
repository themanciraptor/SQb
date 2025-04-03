package sqb

import (
	"context"
	"errors"
)

type Query[T any] struct {
	scanList []interface{}
	query    string
	params   []interface{}

	accumulator Accumulator[T]
}

func (q *Query[T]) GetScanList() []interface{} {
	return q.scanList
}

func (q *Query[T]) GetQuery() string {
	return q.query
}

func (q *Query[T]) GetParams() []interface{} {
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

func (q *Query[T]) Run(ctx context.Context, psql tempDriver) error {
	res, closer, err := psql.RunQuery(ctx, q.query, q.params)
	if err != nil {
		return errors.Join(err, errors.New("failed to run query"))
	}
	defer closer(ctx)

	for res.Next() {
		err := res.Scan(q.scanList...)
		if err != nil {
			return errors.Join(err, errors.New("failed to scan row"))
		}

		q.accumulator.Acc()
	}
	return nil
}
