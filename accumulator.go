package sqb

// In order to properly assign the result, in a generic way without requiring the
// developer to perform type checks. Developers should create a separate accumulator
// object to accumulate results. See common_test.go for an example accumulator
// implementation
type Accumulator[T any] interface {
	Acc()
	GetColumnReceiverMap() map[string]interface{}
	GetResults() []T
}

func NewAccumulator[T any](receiver func(r *T) map[string]interface{}) Accumulator[T] {
	r := new(T)
	return &genericAccumulator[T]{
		ColumnReceiverMap: receiver(r),
		receiver:          r,
	}
}

type genericAccumulator[T any] struct {
	ColumnReceiverMap map[string]interface{}
	// Accumulated Results
	results []T

	// Receiver struct
	receiver *T
}

func (r *genericAccumulator[T]) Acc() {
	r.results = append(r.results, *r.receiver)
}

func (r *genericAccumulator[T]) GetColumnReceiverMap() map[string]interface{} {
	return r.ColumnReceiverMap
}

func (r *genericAccumulator[T]) GetResults() []T {
	return r.results
}
