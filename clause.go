package sqb

/*
	Clauses are the basic building block of queries, they determine how to filter, join and format results. Most
	general use cases should be covered by the query/table builder. Direct usage of the clause API's should only
	be used when a clause is difficult to build in a generic way.
*/

import (
	"fmt"
	"strings"
)

// Clauses Must implement the Clause interface
type Clause interface {
	Build(params *ParamList) string
}

// A filter clause has a single param and a template which defines the column it filters and how it filters it.
type FilterClause struct {
	template   string
	paramValue interface{}
}

// For simple predicates comparing primitive types, a Table will enforce a particular column exists before clause creation.
func NewPrimitiveFilterClause(columnName string, operator string, paramTemplate string, p interface{}) *FilterClause {
	if paramTemplate == "" {
		paramTemplate = "%s"
	}

	return &FilterClause{
		template:   strings.Join([]string{columnName, operator, paramTemplate}, " "),
		paramValue: p,
	}
}

// This defines how any particular clause is built
func (f *FilterClause) Build(params *ParamList) string {
	if f.paramValue == nil {
		return f.template
	}

	return fmt.Sprintf(f.template, params.RecordValueAndReturnParam(f.paramValue))
}

// A CompoundClause is necessary to effectively combine clauses
type CompoundClause struct {
	// list of predicates to be joined by the operator
	predicates []Clause

	// the operator the clauses should be joined with. Is not changed after instantiation
	operator string
}

func NewCompoundClause(operator string) *CompoundClause {
	return &CompoundClause{operator: operator}
}

// Build all sub-clauses and combine them into a single Clause
func (c *CompoundClause) Build(params *ParamList) string {
	template := "%s"

	if len(c.predicates) > 1 {
		template = "(%s)"
	}

	var builtPredicates = make([]string, 0, len(c.predicates))

	for _, predicate := range c.predicates {
		if c, ok := predicate.(Clause); ok {
			builtPredicates = append(builtPredicates, c.Build(params))
		}
	}

	return fmt.Sprintf(template, strings.Join(builtPredicates, fmt.Sprintf(` %s `, c.operator)))
}

// Add a single clause to the CompoundClause.
func (c *CompoundClause) AddClause(nc Clause) *CompoundClause {
	c.predicates = append(c.predicates, nc)
	return c
}

func (c *CompoundClause) NumClauses() int {
	return len(c.predicates)
}

type SortDirection int

const (
	Unset SortDirection = iota
	Ascending
	Descending
)

func getOrderByValue(sd SortDirection) string {
	switch sd {
	case Ascending:
		return "ASC"
	case Descending:
		return "DESC"
	default:
		return ""
	}
}

type OrderByClause struct {
	template      string
	columnName    string
	sortDirection SortDirection
}

func NewOrderByClause(columnName string, sortDirection SortDirection) *OrderByClause {
	return &OrderByClause{
		template:      "%s %s",
		columnName:    columnName,
		sortDirection: Ascending,
	}
}

func (o *OrderByClause) Build(params *ParamList) string {
	if o.columnName == "" {
		return ""
	}

	return fmt.Sprintf(o.template, params.RecordValueAndReturnParam(o.columnName), getOrderByValue(o.sortDirection))
}

type LimitClause struct {
	rowCount int64
	offset   int64
}

func NewLimitClause(rowCount int64, offset int64) *LimitClause {
	return &LimitClause{
		rowCount: rowCount,
		offset:   offset,
	}
}

func (o *LimitClause) Build() string {
	offsetClause := ""
	if o.offset > 0 {
		offsetClause = fmt.Sprint(" OFFSET ", o.offset)
	}

	return fmt.Sprint(" LIMIT ", o.rowCount, offsetClause)
}
