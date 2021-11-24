#Usage and Example Resource

For all examples in this file, refer to the following resource

```go
type foo {
    bar: string
    baz: int
    qux: date
    zap: int
}
```

##Filtering

###Examples

Let’s jump right into the examples which should be enough context for most filtering needs. Refer to the in-depth Parameters section below for more detail.

- `filter_columns=bar&filter_types=EQ&filter_values=test`
  - returns all resources where `bar = test`
- `filter_columns=bar|baz&filter_types=EQ|NE&filter_values=test|1&filter_logic=OR`
  - returns all resources where `bar = test OR baz != 1`
- `filter_columns=bar|baz&filter_types=EQ|NU&filter_values=test|&filter_logic=AND`
  - returns all resources where `bar = test AND baz IS NULL`
- `filter_columns=bar|baz|zap&filter_types=EQ|NE|LT&filter_values=test|1|2&filter_logic=OR|AND`
  - returns all resources where `bar = test OR baz != 1 AND zap < 2`

#Parameters

 - `filter_columns`
   - Specifies the fields that should be filtered on
     - Any fields in the resource being returned can be filtered on. For example, if the endpoint is returning a foo resource: bar, baz, qux and zap are the only valid filter_columns.
     - Multiple `filter_columns` should be specified separated by `|` characters. ex: `filter_columns=bar|baz`.

 - `filter_values`
   - Specifies the values to be used to filter the columns specified in `filter_columns`.
     - Multiple filter_values should be specified separated by a `|` character. ex: `filter_values=test|1`.
     - `filter_values` and `filter_columns` must have the same number of entries.

 - `filter_types`
   - Specifies how `filter_values` are applied to `filter_columns`.
   - Options:
     - `EQ`: the value of column in `filter_columns` must equal the value specified in `filter_values`
     - `EQT`: the value of column in `filter_columns`  must be similar to the value specified in `filter_values` but with tolerance added by the `GTE/LT` filters
     - `NE`: the value of column in `filter_columns` must NOT equal the value specified in `filter_values`
     - `EQT`: the value of the column in `filter_columns` must be equal to the integer (with any decimal point variance) specified in `filter_values`
     - `EQTD`: the value of the column in `filter_columns` must be equal to the day (any time within the 24 hour period) specified in `filter_values`
     - `DR`: the value of the column in `filter_columns` must be equal to a day within the range specified in `filter_values`
     - `DF`: The value of the column in `filter_columns` must be `DISTINCT FROM` the value specified in `filter_values`
       - The difference between `DF` and `NE` is that `DF` will also return null values, while `NE` will not
     - `LT`: the value of column in `filter_columns` must be less than the value specified in `filter_values`
     - `GT`: the value of column in `filter_columns` must greater than the value specified in `filter_values`
     - `NU`: the value of column in `filter_columns` must be `NULL`
     - `NN`: the value of column in `filter_columns` must be `NOT NULL`
     - `LK`: the value of column in `filter_columns` must be `LIKE` the value specified in `filter_values`
     - `ILK`: the value of column in `filter_columns` must be `ILIKE` the value specified in `filter_values`
     - `NLK`: the value of column in `filter_columns` must be `NOT LIKE` the value specified in `filter_values`
     - `NILK`: the value of column in `filter_columns` must be `NOT ILIKE` the value specified in `filter_values`
     - `IN`: the value of column in `filter_columns` must be `IN` the values specified in `filter_values` Multiple `filter_values` should be separated by a `,` character
     - `NIN`: the value of column in `filter_columns` must be `IN` the values specified in `filter_values`
   - Multiple `filter_types` should be separated by a `|` character. ex: `filter_types=EQ|NE`.
   - `filter_values` and `filter_columns` and `filter_types` must have the same number of entries.

 - `filter_logic`
   - Specifies how the filters should be combined.
     - `filter_logic` must have ONE LESS entry than `filter_values`, `filter_columns` and `filter_types`.
   - Options:
     - `AND`: Use and logic
     - `OR`: Use or logic
   - Multiple `filter_logic` values should be separated by a `|` character. ex: `filter_logic=AND|OR`.
   - `filter_logic` must have ONE LESS entry than the number of entries in `filter_values`, `filter_columns` and `filter_types`.

 - `filter_separator`
   - `filter_separator` overrides the default  `|`  separator when using multiple filters at once.

 - `filter_args_separator`
   - `filter_args_separator` overrides the default  `,`  argument separator when using the `IN` and `NIN` filters.

 - `filter_left_parens`
   - `filter_left_parens` indicates all clauses that should have a left parenthesis before them.  Multiple indexes are divided by the `filter_separator`.   These are used to group the logical separators in `filter_logic`.  There must be closing `filter_right_parens` as well.
   - ex. `filter_left_parens=0|1|2` would generate a query select x where `(clause[0] AND (clause[1] AND (clause[2] ...`

 - `filter_right_parens`
   - filter_right_parens indicates all clauses that have a right parenthesis after them.   There must be opening filter_left_parens as well.
   - ex. `filter_left_parens=0|1|2`, `filter_right_parens=2|2|2` would generate a query select x where `(clause[0] AND (clause[1] AND (clause[2])))`

#Custom Columns

Some resources may also have custom filter columns or sort columns.  Defined in the example below is a custom filter column for the foo model.

```go
func (f foo) GetCustomFilters() CustomColumns {
    baz_bar_filter := scope.CustomColumn{
    Name:       "baz_bar",
    ResultType: reflect.TypeOf(""),
    Statement: `(baz::TEXT || bar)`,
    }

    return CustomColumns{baz_bar_filter}
}
```

A custom column can be used as an entry in `filter_columns` if it is returned by `GetCustomFilters`.

A custom column can be used as an entry in `sort_columns` if it is returned by `GetCustomSorts`.

###Example

 - `filter_columns=baz_bar&filter_types=EQ&filter_values=12test`
   - returns all resources where `baz = 12` and `bar = test`
   - and returns all `baz = 1` and `bar = 2test`

##Filter Options

When we write an endpoint that uses the filtering options above, we will provide an additional endpoint with the path suffix `.../filter_options`.   This endpoint is used to fetch all of the available values for that field.

#Example

 - `GET /foos/filter_options?filter_column='bar'`
   - Returns all unique values of `bar` on all `foo`s that would be returned by a call to `GET /foo`

#Sorting

 - `sort_columns`
   - Specifies the fields that should be sorted on
     - Any fields in the resource being returned can be sorted on. For example, if the endpoint is returning a `foo` resource: `bar`, `baz`, `qux` and `zap` are the only valid `sort_columns`.
     - Multiple `sort_columns` should be specified separated by `|` characters. ex: `sort_columns=bar|baz`.
     - Sorting will happen in the order of the `sort_columns`, ex: `sort_columns=bar|baz` will sort by `bar`, then sort within each `bar` entry by `baz`.

 - `sort_directions`
   - Specifies the direction of each sorted field
     - `sort_directions` and `sort_columns` must have the same number of entries.
   - Options:
     - `ASC`: the values in `sort_columns` will be sorted with smallest values first.
     - `DESC`: the values in `sort_columns` will be sorted with largest values first.

#Aggregations

 - `aggregation_column`
   - Specifies the fields that should be aggregated
     - Any fields in the resource being returned can be aggregated on.  For example, if the endpoint is returning a `foo` resource: `bar`, `baz`, `qux` and `zap` are the only valid `aggregation_column` values. However, some aggregation functions are not well defined for non-numeric types, as specified by the postgres documentation for aggregate functions.

 - `aggregation_grouper_column`
   - Specifies the fields that should be grouped by for this aggregation
     - Any fields in the resource being returned can be the grouping column. For example, if the endpoint is returning a `foo` resource: `bar`, `baz`, `qux` and `zap` are the only valid a`ggregation_grouper_column` values.

 - `aggregation_type`
   - Specifies the aggregation to perform.
   - Options:
     - `SUM`: Sums the values in `aggregation_column`.
     - `AVG`: Averages the values in `aggregation_column`.
     - `MIN`: Gets the minimum of the values in `aggregation_column`.
     - `MAX`: Gets the maximum of the values in `aggregation_column`.
     - `COUNT`: Returns the number of rows returned.

###Example

 - `GET /foos/aggregate?aggregation_column='bar'&aggregation_type='SUM'`
   - Returns the numeric sum of `bar` on all `foo`s that would be returned by a call to `GET /foo`

 - `GET /foos/grouped_aggregate?aggregation_column='bar'&aggregation_type='COUNT'&aggregation_grouper_column='bax'`
   - Returns the count of `bar` on all `foo`s that would be returned by a call to `GET /foo` grouped into buckets by the values in `bax`.    Will return a list of tuples with the of the format `[{“grouper”:{{value_in_bax_1}}, “result”:1} … ]`

#Pagination

This package uses the default pagination utility provided by Go Buffalo.
https://github.com/gobuffalo/pop/blob/master/paginator.go

 - `page`
    - Specifies the current integer page of results.
      - The default `page` is 1.

 - `per_page`
   - Specifies the amount of records to return on each page of results.
     - The default `per_page` is 20.
     