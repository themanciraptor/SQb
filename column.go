package sqb

import "reflect"

/*
	A column represents a column in the query. It contains data about the type a column receiver
	should have, and a pointer to the receiver value.
*/

type Column struct {
	kind     reflect.Kind
	receiver interface{}
}

func (s *Column) SetReceiver(v interface{}) {
	s.receiver = v
}

// must be initialized with a columnKind or the column would be unable to perform typeChecking
func NewColumn(columnKind reflect.Kind) *Column {
	return &Column{
		kind:     columnKind,
		receiver: nil,
	}
}
