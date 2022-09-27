package sqb

import "fmt"

// Any variance in dialects should be accounted for here.
type Dialect interface {
	StructTag() string
	FormatParam(n int) string
}

type psql struct{}

func (p psql) StructTag() string {
	return "psql"
}

func (p psql) FormatParam(n int) string {
	return fmt.Sprintf("$%d", n)
}

func Psql() Dialect {
	return psql{}
}
