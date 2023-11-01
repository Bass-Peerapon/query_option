package query_option_test

import (
	"reflect"
	"testing"

	"github.com/Bass-Peerapon/query_option"
)

func TestConvertToPostgresFilter(t *testing.T) {
	tests := []struct {
		name     string
		input    query_option.QueryOption
		expected string
		args     []interface{}
	}{
		{
			name:     "Test $eq condition",
			input:    query_option.QueryOption{Filter: map[string]interface{}{"key": map[string]interface{}{"$eq": "value"}}},
			expected: "WHERE key = ?",
			args:     []interface{}{"value"},
		},
		{
			name:     "Test $q condition",
			input:    query_option.QueryOption{Filter: map[string]interface{}{"key": map[string]interface{}{"$q": "value"}}},
			expected: "WHERE key ILIKE ?",
			args:     []interface{}{"%value%"},
		},
		{
			name:     "Test $gt condition",
			input:    query_option.QueryOption{Filter: map[string]interface{}{"key": map[string]interface{}{"$gt": 5}}},
			expected: "WHERE key > ?",
			args:     []interface{}{5},
		},
		{
			name:     "Test $gte condition",
			input:    query_option.QueryOption{Filter: map[string]interface{}{"key": map[string]interface{}{"$gte": 5}}},
			expected: "WHERE key >= ?",
			args:     []interface{}{5},
		},
		{
			name:     "Test $lt condition",
			input:    query_option.QueryOption{Filter: map[string]interface{}{"key": map[string]interface{}{"$lt": 5}}},
			expected: "WHERE key < ?",
			args:     []interface{}{5},
		},
		{
			name:     "Test $lte condition",
			input:    query_option.QueryOption{Filter: map[string]interface{}{"key": map[string]interface{}{"$lte": 5}}},
			expected: "WHERE key <= ?",
			args:     []interface{}{5},
		},
		{
			name:     "Test $in condition",
			input:    query_option.QueryOption{Filter: map[string]interface{}{"key": map[string]interface{}{"$in": []int{1, 2, 3}}}},
			expected: "WHERE key IN (?, ?, ?)",
			args:     []interface{}{1, 2, 3},
		},
		{
			name:     "Test $and condition with multiple conditions",
			input:    query_option.QueryOption{Filter: map[string]interface{}{"$and": []map[string]interface{}{{"key1": map[string]interface{}{"$eq": "value1"}}, {"key2": map[string]interface{}{"$gt": 5}}}}},
			expected: "WHERE (key1 = ? AND key2 > ?)",
			args:     []interface{}{"value1", 5},
		},
		{
			name:     "Test $or condition with multiple conditions",
			input:    query_option.QueryOption{Filter: map[string]interface{}{"$or": []map[string]interface{}{{"key1": map[string]interface{}{"$eq": "value1"}}, {"key2": map[string]interface{}{"$lt": 5}}}}},
			expected: "WHERE (key1 = ? OR key2 < ?)",
			args:     []interface{}{"value1", 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, args := tt.input.ConvertToPostgresFilter()
			if result != tt.expected || !reflect.DeepEqual(args, tt.args) {
				t.Errorf("expected %s %v, got %s %v", tt.expected, tt.args, result, args)
			}
		})
	}
}
