package util

import "reflect"

func GetFieldIntValue(v interface{}, field string) int {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return int(f.Int())
}

func GetFieldStringValue(v interface{}, field string) string {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.String()
}

func GetFieldBoolValue(v interface{}, field string) bool {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f.Bool()
}

func GetFieldValue(v interface{}, field string) reflect.Value {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(field)
	return f
}
