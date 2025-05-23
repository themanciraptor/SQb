package sqb

import (
	"fmt"
	"reflect"
)

// TODO (SSC-3684): There is likely many cases where the default builder filters don't
//
//	adequately cover our needs. In these cases, we will want to be able
//	to build custom filters.
//
// Provide public helper so custom clauses can easily check if they are valid
func (t *Table[T]) AssertFilterClauseValid(columnName string, param interface{}) {
	var column *Column

	if c, ok := t.fields[columnName]; !ok {
		keys := make([]string, 0, len(t.fields))
		for key := range t.fields {
			keys = append(keys, key)
		}

		panic(fmt.Sprintf("No column named %s found for table %s, available columns: %v", columnName, t.tableName, keys))
	} else {
		column = c
	}

	if reflect.TypeOf(param).Kind() != column.kind && reflect.TypeOf(param) != nullType(column.kind) {
		panic(fmt.Sprintf("Incorrect type for column. Need %s, got %s",
			column.kind,
			reflect.TypeOf(param),
		))
	}
}

func (t *Table[T]) AssertColumnExists(columnName string) {
	if _, ok := t.fields[columnName]; !ok {
		panic(fmt.Sprintf("No column named %s found for table %s", columnName, t.tableName))
	}
}

func (t *Table[T]) ColumnEquals(columnName string, v interface{}) *Table[T] {
	t.AssertFilterClauseValid(columnName, v)

	t.filter.AddClause(NewPrimitiveFilterClause(columnName, "=", "%s", v))

	return t
}

func (t *Table[T]) ColumnNull(columnName string) *Table[T] {
	t.AssertColumnExists(columnName)

	t.filter.AddClause(NewPrimitiveFilterClause(columnName, "IS", "NULL", nil))

	return t
}

func (t *Table[T]) BuildFilter(params *ParamList) string {
	return t.filter.Build(params)
}

func (t *Table[T]) AddOrderByClause(columnName string, sortDirection SortDirection) *Table[T] {
	t.orderBy.AddClause(NewOrderByClause(columnName, sortDirection))

	return t
}

func (t *Table[T]) Limit(rowCount int64, offset int64) *Table[T] {
	t.limit = NewLimitClause(rowCount, offset)

	return t
}
