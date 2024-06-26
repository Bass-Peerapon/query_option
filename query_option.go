package query_option

import (
	"fmt"
	"strings"
)

// $eq: Matches values that are equal to a specified value.
// Example:
//
//	{
//	  "key": {
//	    "$eq": "value"
//	  }
//	}
//
// or the simplest form
//
//	{
//	  "key": "value"
//	}
//
// $q: Full text search (matches values where the whole text value matches the specified value)
// Example:
//
//	{
//	  "key": {
//	    "$q": "value
//	  }
//	}
//
// $gt: Matches values that are greater than a specified value.
// Example:
//
//	{
//	  "key": {
//	    "$gt": 4
//	  }
//	}
//
// $gte: Matches values that are greater than or equal to a specified value.
// Example:
//
//	{
//	  "key": {
//	    "$gte": 4
//	  }
//	}
//
// $lt: Matches values that are less than a specified value.
// Example:
//
//	{
//	  "key": {
//	    "$lt": 4
//	  }
//	}
//
// $lte: Matches values that are less than or equal to a specified value.
// Example:
//
//	{
//	  "key": {
//	    "$lte": 4
//	  }
//	}
//
// $in: Matches any of the values specified in an array.
// Example:
//
//	{
//	  "key": {
//	    "$in": [
//	      1,
//	      2,
//	      4
//	    ]
//	  }
//	}
//
// $and: Matches all the values specified in an array.
// Example:
//
//	{
//	  "$and": [
//	    {
//	      "key": {
//	        "$in": [
//	          1,
//	          2,
//	          4
//	        ]
//	      }
//	    },
//	    {
//	      "some_other_key": 10
//	    }
//	  ]
//	}
//
// $or: Matches at least one of the values specified in an array.
// Example:
//
//	{
//	  "$or": [
//	    {
//	      "key": {
//	        "$in": [
//	          1,
//	          2,
//	          4
//	        ]
//	      }
//	    },
//	    {
//	      "key2": 10
//	    }
//	  ]
//	}
type M map[string]any

type QueryOption struct {
	Filter  M             `json:"filter_conditions,omitempty"`
	Sort    []*SortOption `json:"sort,omitempty"`
	PerPage int           `json:"per_page,omitempty"` // pagination option: limit number of results
	Page    int           `json:"page,omitempty"`     // pagination option: offset to return items from
}

type SortOption struct {
	Field     string `json:"field"`     // field name to sort by,from json tags(in camel case), for example created_at
	Direction int    `json:"direction"` // [-1, 1] Direction can be ascending (1) or descending (-1)
}

func (queryOption QueryOption) ConvertToPostgresFilter() (string, []interface{}) {
	var conditions []string
	var args []interface{}

	for key, value := range queryOption.Filter {
		conds, valArgs := handleLogicalOperators(key, value)
		conditions = append(conditions, conds...)
		args = append(args, valArgs...)
	}

	filterClause := ""
	if len(conditions) > 0 {
		filterClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	return filterClause, args
}

func (queryOption QueryOption) ConvertToPostgresSort() string {
	var conditions []string
	for _, sort := range queryOption.Sort {
		var direction string
		if sort.Direction > 0 {
			direction = "ASC"
		} else {
			direction = "DESC"
		}
		conditions = append(conditions, fmt.Sprintf("%s %s", sort.Field, direction))
	}
	var sql string
	if len(conditions) > 0 {
		sql += "ORDER " + strings.Join(conditions, ", ")
	}
	return sql
}

func (quertOption QueryOption) ConvertToPostgresPagination() string {
	var paginatorSql string
	limit := quertOption.PerPage
	skipItem := (quertOption.Page - 1) * quertOption.PerPage
	paginatorSql = fmt.Sprintf(`
			LIMIT %d
			OFFSET %d
			`,
		limit,
		skipItem,
	)

	return paginatorSql
}

func handleLogicalOperators(key string, value interface{}) ([]string, []interface{}) {
	var conditions []string
	var args []interface{}
	if inValue, ok := value.(map[string]interface{}); ok {
		for k, v := range inValue {
			switch k {
			case "$eq":
				conditions = append(conditions, fmt.Sprintf("%s = ?", key))
				args = append(args, v)
			case "$q":
				conditions = append(conditions, fmt.Sprintf("%s ILIKE ?", key))
				args = append(args, "%"+v.(string)+"%")
			case "$gt":
				conditions = append(conditions, fmt.Sprintf("%s > ?", key))
				args = append(args, v)
			case "$gte":
				conditions = append(conditions, fmt.Sprintf("%s >= ?", key))
				args = append(args, v)
			case "$lt":
				conditions = append(conditions, fmt.Sprintf("%s < ?", key))
				args = append(args, v)
			case "$lte":
				conditions = append(conditions, fmt.Sprintf("%s <= ?", key))
				args = append(args, v)
			case "$in":
				inValues := toSlice(v)
				placeholders := make([]string, len(inValues))
				for i := range inValues {
					args = append(args, inValues[i])
					placeholders[i] = "?"
				}
				conditions = append(conditions, fmt.Sprintf("%s IN (%s)", key, strings.Join(placeholders, ", ")))

			}
		}
	} else if key == "$and" || key == "$or" {
		if values, ok := value.([]map[string]interface{}); ok {
			var subConditions []string
			var subArgs []interface{}
			for _, m := range values {
				for k, v := range m {
					conds, valArgs := handleLogicalOperators(k, v)
					subConditions = append(subConditions, conds...)
					subArgs = append(subArgs, valArgs...)
				}
			}
			if len(subConditions) > 0 {
				conditions = append(conditions, "("+strings.Join(subConditions, fmt.Sprintf(" %s ", strings.Trim(strings.ToUpper(key), "$")))+")")
				args = append(args, subArgs...)
			}
		}
	}

	return conditions, args
}

// toSlice casts an interface to a []interface{} type.
func toSlice(i interface{}) []interface{} {
	v, _ := toSliceE(i)
	return v
}

// toSliceE casts an interface to a []interface{} type.
func toSliceE(i interface{}) ([]interface{}, error) {
	var s []interface{}

	switch v := i.(type) {
	case []interface{}:
		return append(s, v...), nil
	case []map[string]interface{}:
		for _, u := range v {
			s = append(s, u)
		}
		return s, nil
	case []int:
		for _, v2 := range v {
			s = append(s, v2)
		}
		return s, nil

	case []string:
		for _, v2 := range v {
			s = append(s, v2)
		}
		return s, nil

	case []float32:
		for _, v2 := range v {
			s = append(s, v2)
		}
		return s, nil

	case []float64:
		for _, v2 := range v {
			s = append(s, v2)
		}
		return s, nil

	default:
		return s, fmt.Errorf("unable to cast %#v of type %T to []interface{}", i, i)
	}
}

// test merge
