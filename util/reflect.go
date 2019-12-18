package util

import (
	"fmt"
	"reflect"
)

// ReflectFieldValueAsString returns the value of a field
func ReflectFieldValueAsString(v reflect.Value, field string) string {
	return fmt.Sprint(reflect.Indirect(v).FieldByName(field))
}
