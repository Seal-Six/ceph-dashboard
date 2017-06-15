package main

import (
	"fmt"
	"reflect"
)

func HasField(v interface{}, name string) bool {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	if rv.Kind() != reflect.Struct {
		return false
	}
	return rv.FieldByName(name).IsValid()
}

func SprintField(v interface{}, name string) string {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	if rv.Kind() == reflect.Struct {
		if !rv.FieldByName(name).IsValid() {
			return ""
		}

		return fmt.Sprint(rv.FieldByName(name))
	}

	if rv.Kind() == reflect.Map {
		rt := reflect.TypeOf(v)
		if rt.Key().Kind() != reflect.String {
			return ""
		}

		v := rv.MapIndex(reflect.ValueOf(name))
		if v.Kind() == reflect.Invalid {
			return ""
		}

		return fmt.Sprint(v)
	}

	return ""
}
