# SQb (Not production ready)

A generic query builder intended to support query building in a generic way and make more of the query code we write reusable. For now, this is only intended to perform accessor queries.

The main advantage will be composability. Joins are intended to be programmatic, allowing us to create functions that not only return simple queries for tables, but also return more complex joins as a single object on which more filters or joins can be performed.

This is not intended to be built like an ORM. Result objects are defined separately from domain models.

# How it works

We use the metadata provided by our models to allow us to perform type checking on our queries during development and perform a lot of the gruntwork automatically. To define a typical query, the following steps are necessary:

1. Create a table.
2. Define a result accumulator. See [Example Accumulator](common_test.go)
3. Build the query.
4. Use the query's run function to run the query.
5. Get results from the accumulator.

# Terminology

### Column

A mapping between a column name, and a 2-tuple of its kind and receiver. If a receiver value is set, it is automatically added to the select statement against that table.

### CompoundClause

A container for multiple clauses. Builds its own clauses iteratively. Clauses can be simple clauses or more CompoundClauses. These are intended to be abstracted away from developers except in cases where the provided filters do not cover the logic necessary. Before building a CompoundClause, always check to see if more generic filters will support your use case.

### Dialect

Currently used as a catch-all for major differences between sql implementations. This holds information like how to define a variable within a sql statement.

### FilterClause

Defines filter statements within an sql query. Filters are defined in the following format:
`column operator input`, e.g. `name = "Carl"`

### ParamList

Handles the mapping between params and their corresponding SQL variable, for sql prepared
statements.

### Query

Contains information necessary to query an sql table such as the query string, the params
for user input, and the receivers each query row will be place in.

### Table

currently acts as a combined table definition and query builder. A table does the following:

- builds and maintains columns definitions (map\[column\]{type, receiver}).
- includes a CompoundClause for filtering
- Handles building a query (as a first slice)
