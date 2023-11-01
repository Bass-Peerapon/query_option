## Query Operators
The file provides documentation and examples for the following query operators:
- $eq: Matches values equal to a specified value.
- $q: Full text search.
- $gt: Matches values greater than a specified value.
- $gte: Matches values greater than or equal to a specified value.
- $lt: Matches values less than a specified value.
- $lte: Matches values less than or equal to a specified value.
- $in: Matches any of the values specified in an array.
- $and: Matches all the values specified in an array.
- $or: Matches at least one of the values specified in an array.

## Structures
### QueryOption
- Field: Filter (type: map[string]any)

## Methods
### ConvertToPostgresFilter (for QueryOption)
- Returns: string, slice of interface{}

## Installation
```shell
go get -u github.com/Bass-Peerapon/query_option
```
## Usage

1. Create a new `QueryOption` instance:

```go
opt := query_option.QueryOption{
    Filter: map[string]any{
        "key": map[string]any{
            "$eq" : "test",
        },
    }
}
```

2. Use the `ConvertToPostgresFilter` method to convert the query option to a Postgres filter:

```go
filter, args := opt.ConvertToPostgresFilter()
```

3. Use the `ToSlice` and `ToSliceE` helper functions as needed for data type conversions.

