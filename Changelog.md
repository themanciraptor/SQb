## 0.0.1
Add the following features:
- CompoundClause enables us to build complex filters in a generic way
- filterClause provides simple clauses of the form `column operator input`, e.g. `name = "Carl"`
- Dialect allows support for other dialects to be defined separately and injected.
- ParamList handles the mapping between params and their corresponding SQL variable.
- Query provides a way for us to generically run queries against the database.
- Column is a mapping between a column name and its value receiver. This allows the query buildingto define scan order automatically.
- Table currently acts as a combined table definition and query builder. A table does thefollowing:
    - builds and maintains columns.
    - includes a CompoundClause for filtering
    - Handles building a query (as a first slice)
        