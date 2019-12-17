package util

import (
	"fmt"
	"reflect"
)

// GetFieldValueAsString returns the value of a field
func GetFieldValueAsString(v interface{}, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return fmt.Sprint(f)
}

// ReflectFieldValueAsString returns the value of a field
func ReflectFieldValueAsString(v reflect.Value, field string) string {
	return fmt.Sprint(reflect.Indirect(v).FieldByName(field))
}
