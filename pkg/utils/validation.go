package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

// ValidateAndCollectErrors validates the input struct and returns any custom and validation errors found.
func ValidateAndCollectErrors(w http.ResponseWriter, statusCode int, errors map[string]string, error interface{}, obj interface{}) (map[string]string) {
	// errors := make(map[string]string)
	fieldMap := GetJSONFieldMap(obj)
	for _, err := range error.(validator.ValidationErrors) {
		// errors[err.Field()] = err.Field() + " is " + err.Tag() + " " + err.Param()
		jsonField, exists := fieldMap[err.Field()]
		if !exists {
			jsonField = err.Field() // Fallback to struct field name if JSON tag not available
		}
		errors[jsonField] = FormatErrorMessage(jsonField, err)
	}

	return errors
}

// WriteJSONEValidation dynamically gets JSON tag fields and writes a validation error response
func ValidationResponse(w http.ResponseWriter, statusCode int, error interface{}, obj interface{}) {
	errors := make(map[string]string)

	// errors := make(map[string]string)
	fieldMap := GetJSONFieldMap(obj)
	for _, err := range error.(validator.ValidationErrors) {
		// errors[err.Field()] = err.Field() + " is " + err.Tag() + " " + err.Param()
		jsonField, exists := fieldMap[err.Field()]
		if !exists {
			jsonField = err.Field() // Fallback to struct field name if JSON tag not available
		}
		errors[jsonField] = FormatErrorMessage(jsonField, err)
	}
	response := ErrorResponse{
		Success: false,
		Message: "Validation error",
		Errors:  errors,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

func FormatErrorMessage(field string, err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", field)
	case "max":
		return fmt.Sprintf("%s cannot be longer than %s characters", field, err.Param())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters long", field, err.Param())
	case "email":
		return fmt.Sprintf("%s must be a valid email address", field)
	case "boolean":
		return fmt.Sprintf("%s must be a boolean value", field)
	case "numeric":
		return fmt.Sprintf("%s must be a numeric value", field)
	case "gte":
		return fmt.Sprintf("%s must be greater than or equal to %s", field, err.Param())
	case "lte":
		return fmt.Sprintf("%s must be less than or equal to %s", field, err.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", field, err.Param())
	case "lt":
		return fmt.Sprintf("%s must be less than %s", field, err.Param())
	case "len":
		return fmt.Sprintf("%s must be exactly %s characters long", field, err.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", field, err.Param())
	case "alphanum":
		return fmt.Sprintf("%s must contain only alphanumeric characters", field)
	case "url":
		return fmt.Sprintf("%s must be a valid URL", field)
	case "uuid":
		return fmt.Sprintf("%s must be a valid UUID", field)
	case "hexadecimal":
		return fmt.Sprintf("%s must be a valid hexadecimal value", field)
	case "base64":
		return fmt.Sprintf("%s must be a valid base64 encoded string", field)
	case "contains":
		return fmt.Sprintf("%s must contain the substring '%s'", field, err.Param())
	case "excludes":
		return fmt.Sprintf("%s must not contain the substring '%s'", field, err.Param())
	case "excludesall":
		return fmt.Sprintf("%s must not contain any of the characters '%s'", field, err.Param())
	case "excludesrune":
		return fmt.Sprintf("%s must not contain the rune '%s'", field, err.Param())
	case "startswith":
		return fmt.Sprintf("%s must start with '%s'", field, err.Param())
	case "endswith":
		return fmt.Sprintf("%s must end with '%s'", field, err.Param())
	case "datetime":
		return fmt.Sprintf("%s must be a valid datetime matching the format '%s'", field, err.Param())
	default:
		return fmt.Sprintf("%s is invalid", field)
	}
}

// GetJSONFieldMap uses reflection to map struct fields to their JSON tags
func GetJSONFieldMap(obj interface{}) map[string]string {
	fieldMap := make(map[string]string)
	val := reflect.ValueOf(obj)

	// Ensure we're working with a pointer to a struct
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return fieldMap
	}

	// Recursive helper function to extract JSON tags
	var extractFields func(reflect.Value, string)
	extractFields = func(v reflect.Value, parent string) {
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i)
			fieldValue := v.Field(i)
			jsonTag := field.Tag.Get("json")

			// Skip unexported fields and those without a JSON tag
			if !fieldValue.CanInterface() || jsonTag == "-" {
				continue
			}

			// Construct the JSON field name, including parent struct names if nested
			jsonFieldName := field.Name
			if jsonTag != "" {
				jsonFieldName = strings.Split(jsonTag, ",")[0]
			}
			if parent != "" {
				jsonFieldName = parent + "." + jsonFieldName
			}

			// Add to the field map
			fieldMap[field.Name] = jsonFieldName

			// Handle nested structs recursively
			if fieldValue.Kind() == reflect.Struct {
				extractFields(fieldValue, jsonFieldName)
			}
		}
	}

	// Start the extraction process
	extractFields(val, "")
	return fieldMap
}
