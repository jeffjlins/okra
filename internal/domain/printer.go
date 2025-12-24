package domain

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func print(f reflect.Value) string {
	//var the_value = reflect.ValueOf(f.Interface())
	switch f.Kind() {
	case reflect.Ptr:
		if f.Interface() == nil {
			return "*<nil>"
		}
		return "*" + print(f.Elem())
	case reflect.Struct:
		return printStruct(f)
	case reflect.String:
		return "\"" + f.String() + "\""
	case reflect.Bool:
		return strconv.FormatBool(f.Bool())
	case reflect.Int:
		return strconv.FormatInt(f.Int(), 10)
	case reflect.Float32:
		return strconv.FormatFloat(f.Float(), 'g', -1, 32)
	case reflect.Float64:
		return strconv.FormatFloat(f.Float(), 'g', -1, 64)
	case reflect.Slice:
		return printSequence(f)
	case reflect.Array:
		return printSequence(f)
	}
	return fmt.Sprintf("%s", f)
}

func printStruct(r reflect.Value) string {
	var fields []string
	for i := 0; i < r.NumField(); i++ {
		fields = append(fields, r.Type().Field(i).Name+"="+print(r.Field((i))))
	}
	return r.Type().Name() + "{ " + strings.Join(fields, ", ") + " }"
}

func printSequence(r reflect.Value) string {
	if r.Len() == 0 {
		return "[]"
	}
	var elements []string
	for i := 0; i < r.Len(); i++ {
		elements = append(elements, print(r.Index(i)))
	}
	return "[" + strings.Join(elements, ", ") + "]"
}
