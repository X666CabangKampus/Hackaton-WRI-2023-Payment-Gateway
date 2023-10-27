package util

import (
	"reflect"
	"strings"
)

func getJSONFieldName(field reflect.StructField) string {
	// Get the JSON tag for the field.
	jsonTag := field.Tag.Get("json")

	// If the JSON tag is empty, return the field name.
	if jsonTag == "" {
		return field.Name
	}

	// Split the JSON tag by comma to get the field name.
	jsonFields := strings.Split(jsonTag, ",")
	jsonFieldName := jsonFields[0]

	// Return the JSON field name.
	return jsonFieldName
}

func StructToMap(structPtr interface{}) map[string]interface{} {
	structValue := reflect.ValueOf(structPtr)
	structType := structValue.Type()
	numFields := structType.NumField()

	result := make(map[string]interface{})
	for i := 0; i < numFields; i++ {
		field := structType.Field(i)
		fieldValue := structValue.Field(i)
		fieldName := getJSONFieldName(field)
		if !reflect.DeepEqual(fieldValue.Interface(), reflect.Zero(fieldValue.Type()).Interface()) {
			result[fieldName] = fieldValue.Interface()
		}
	}
	return result
}
