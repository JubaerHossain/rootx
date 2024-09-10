package jsonQuery

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/JubaerHossain/rootx/pkg/splitters"
)

// JQ (JSON Query) struct
type JQ struct {
	Data any
}

// NewFileQuery - Create a new &JQ from a JSON file.
func NewFileQuery(jsonFile string) (*JQ, error) {
	raw, err := os.ReadFile(jsonFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	var data any
	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON from file: %w", err)
	}
	return &JQ{Data: data}, nil
}

// NewStringQuery - Create a new &JQ from a raw JSON string.
func NewStringQuery(jsonString string) (*JQ, error) {
	var data any
	if err := json.Unmarshal([]byte(jsonString), &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON string: %w", err)
	}
	return &JQ{Data: data}, nil
}

// NewQuery - Create a &JQ from any parsed by json.Unmarshal
func NewQuery(jsonObject any) *JQ {
	return &JQ{Data: jsonObject}
}

// Query - queries against the JSON with the expression passed in. The exp is separated by dots (".")
func (jq *JQ) Query(exp string) (any, error) {
	if exp == "." {
		return jq.Data, nil
	}

	paths, err := splitters.SplitArgs(exp, ".", false)
	if err != nil {
		return nil, fmt.Errorf("failed to split query expression: %w", err)
	}

	var context = jq.Data
	for _, path := range paths {
		if isArrayPath(path) {
			index, err := parseArrayIndex(path)
			if err != nil {
				return nil, err
			}

			arr, ok := context.([]any)
			if !ok {
				return nil, fmt.Errorf("%s is not an array", path)
			}
			if index >= len(arr) {
				return nil, fmt.Errorf("index %d out of range in %s", index, path)
			}
			context = arr[index]
		} else {
			obj, ok := context.(map[string]any)
			if !ok {
				return nil, fmt.Errorf("%s is not a JSON object", path)
			}
			val, exists := obj[path]
			if !exists {
				return nil, fmt.Errorf("%s does not exist in the JSON object", path)
			}
			context = val
		}
	}
	return context, nil
}

// QueryToMap - Queries and converts the result to a map[string]any
func (jq *JQ) QueryToMap(exp string) (map[string]any, error) {
	result, err := jq.Query(exp)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	ret, ok := result.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected a JSON object but got a different type for: %s", exp)
	}
	return ret, nil
}

// QueryToArray - Queries and converts the result to a slice of any ([]any)
func (jq *JQ) QueryToArray(exp string) ([]any, error) {
	result, err := jq.Query(exp)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}

	ret, ok := result.([]any)
	if !ok {
		return nil, fmt.Errorf("expected an array but got a different type for: %s", exp)
	}
	return ret, nil
}

// QueryToString - Queries and converts the result to a string
func (jq *JQ) QueryToString(exp string) (string, error) {
	result, err := jq.Query(exp)
	if err != nil {
		return "", fmt.Errorf("failed to query: %w", err)
	}

	ret, ok := result.(string)
	if !ok {
		return "", fmt.Errorf("expected a string but got a different type for: %s", exp)
	}
	return ret, nil
}

// QueryToInt64 - Queries and converts the result to an int64
func (jq *JQ) QueryToInt64(exp string) (int64, error) {
	result, err := jq.Query(exp)
	if err != nil {
		return 0, fmt.Errorf("failed to query: %w", err)
	}

	switch v := result.(type) {
	case float64:
		return int64(v), nil
	case string:
		ret, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to convert string to int64: %w", err)
		}
		return ret, nil
	default:
		return 0, fmt.Errorf("expected an int64-compatible value but got: %v", result)
	}
}

// QueryToFloat64 - Queries and converts the result to a float64
func (jq *JQ) QueryToFloat64(exp string) (float64, error) {
	result, err := jq.Query(exp)
	if err != nil {
		return 0, fmt.Errorf("failed to query: %w", err)
	}

	ret, ok := result.(float64)
	if !ok {
		return 0, fmt.Errorf("expected a float64 but got a different type for: %s", exp)
	}
	return ret, nil
}

// QueryToBool - Queries and converts the result to a boolean
func (jq *JQ) QueryToBool(exp string) (bool, error) {
	result, err := jq.Query(exp)
	if err != nil {
		return false, fmt.Errorf("failed to query: %w", err)
	}

	ret, ok := result.(bool)
	if !ok {
		return false, fmt.Errorf("expected a boolean but got a different type for: %s", exp)
	}
	return ret, nil
}

// Helper functions

func isArrayPath(path string) bool {
	return strings.HasPrefix(path, "[") && strings.HasSuffix(path, "]")
}

func parseArrayIndex(path string) (int, error) {
	index, err := strconv.Atoi(path[1 : len(path)-1])
	if err != nil {
		return 0, fmt.Errorf("invalid array index %s: %w", path, err)
	}
	return index, nil
}
