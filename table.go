package sqb

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/lib/pq"
)

/*
Table's hold information about a specific database table. They are essentially helper structs
developer can use to easily build queries against our database.
*/
type Table struct {
	// The literal name of the table within the database
	tableName string

	// Map table columns to a value receivers.
	fields map[string]*Column

	// The filter clause applied to the table
	filter *CompoundClause

	// Ordering Applied to Table
	orderBy *CompoundClause

	// Limit
	limit *LimitClause
}

// Set receiver for a particular table column. The column must exist on the table.
func (t *Table) SetColumnReceiver(columnName string, scanTo interface{}) *Table {
	if reflect.TypeOf(scanTo).Kind() != reflect.Ptr {
		panic("SetField: scanTo must be reference pointer")
	}

	if pair, ok := t.fields[columnName]; ok {
		indirectVal := reflect.Indirect(reflect.ValueOf(scanTo))

		if indirectVal.Kind() != pair.kind && indirectVal.Type() != nullType(pair.kind) {
			panic(fmt.Sprintf("SetField: attempted to scan to invalid type for field: cannot assign %v to: %v", reflect.TypeOf(scanTo).Kind(), pair.kind))
		}

		if reflect.Indirect(reflect.ValueOf(scanTo)).Kind() == reflect.Slice {
			t.fields[columnName].SetReceiver(pq.Array(scanTo))
		}

		t.fields[columnName].SetReceiver(scanTo)
		return t
	}

	panic(fmt.Sprintf("SetField: column not included in table %v for %s", t.fields, columnName))
}

func (t *Table) LoadReceiversFromAccumulator(a Accumulator) *Table {
	columnErrors := map[string]interface{}{}

	receivers := a.GetColumnReceiverMap()

	for columnName, receiver := range receivers {
		if reflect.TypeOf(receiver).Kind() != reflect.Ptr {
			columnErrors[columnName] = "receiver must be reference pointer"
		}

		if pair, ok := t.fields[columnName]; ok {
			indirectVal := reflect.Indirect(reflect.ValueOf(receiver))

			if indirectVal.Kind() != pair.kind && indirectVal.Type() != nullType(pair.kind) {
				columnErrors[columnName] = fmt.Sprintf("Invalid type for field: cannot assign %v to: %v", reflect.TypeOf(receiver).Kind(), pair.kind)
				continue
			}

			if reflect.Indirect(reflect.ValueOf(receiver)).Kind() == reflect.Slice {
				t.fields[columnName].SetReceiver(pq.Array(receiver))
				continue
			}

			t.fields[columnName].SetReceiver(receiver)
		} else {
			columnErrors[columnName] = "Column not included in table"
		}
	}

	if len(columnErrors) > 0 {
		e := strings.Builder{}

		for c, v := range columnErrors {
			e.WriteString(fmt.Sprintf("%s: %s\n", c, v))
		}

		panic(e.String())
	}

	return t
}

func NewTable(tableName string, dialect Dialect, model interface{}) *Table {
	if reflect.TypeOf(model).Kind() != reflect.Ptr {
		panic("QueryBuilder: Table model must be pointer to struct type")
	}

	modelValue := reflect.Indirect(reflect.ValueOf(model))

	if modelValue.Kind() != reflect.Struct {
		panic("QueryBuilder: Table model must be pointer to struct type")
	}

	table := &Table{
		tableName: tableName,
		fields:    map[string]*Column{},
		filter:    NewCompoundClause("AND"),
		orderBy:   NewCompoundClause(","),
	}

	// Provide default columns based on the table model
	for i := 0; i < modelValue.NumField(); i++ {
		kind := reflect.TypeOf(model).Elem().Field(i).Type.Kind()

		c := reflect.TypeOf(model).Elem().Field(i).Tag.Get(dialect.StructTag())
		table.fields[c] = NewColumn(kind)
	}

	return table
}

// TODO: (SSC-3681): Query building should be owned by a separate query builder
// TODO: (SSC-3682): Joins will be supported through query builder, not tables themselves.
//
//	May require a separate table-like object to define joins and their
//	types. Ideally, can support within the table struct when query building
//	has been extracted.
//
// TODO: (SSC-3683): We may want to adjust the select statement depending on the tables
//
//	context. This support does not rely on the query builder per se. But
//	having the query builder already will make implementation easier.
func (t *Table) Build(a Accumulator, dialect Dialect) *Query {
	selectedFields := make([]string, 0, len(t.fields))
	scanList := make([]interface{}, 0, len(t.fields))
	paramList := NewParamList(dialect)
	filters := ""

	if t.filter.NumClauses() > 0 {
		filters = fmt.Sprint(` WHERE `, t.filter.Build(paramList))
	}

	keys := make([]string, 0, len(t.fields))
	for key := range t.fields {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, columnName := range keys {
		sp := t.fields[columnName]

		if sp.receiver != nil {
			selectedFields = append(selectedFields, columnName)
			scanList = append(scanList, sp.receiver)
		}
	}

	limitClause := ""
	if t.limit != nil {
		limitClause = t.limit.Build()
	}

	orderClause := ""
	if t.orderBy.NumClauses() > 0 {
		orderClause = t.orderBy.Build(paramList)
	}

	return &Query{
		query:    fmt.Sprint(`SELECT `, strings.Join(selectedFields, ", "), ` FROM `, t.tableName, filters, orderClause, limitClause),
		scanList: scanList,
		params:   paramList.GetParamList(),

		accumulator: a,
	}
}
