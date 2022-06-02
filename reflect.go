// Copyright (c) 2021 Jing-Ying Chen. Subject to the MIT License.

package goutil

import (
	"fmt"
	"reflect"
)

type Nilable interface {
	IsNil() bool
}

func IsNull(v Nilable) bool {
	return v == nil || v.IsNil()
}

// panic on purpose (for interface param known to be struct pointer)
func IsNilPtr(v interface{}) bool {
	return v == nil || reflect.ValueOf(v).IsNil()
}

// for app initialization or less frequent uses
func IsNil(v interface{}) bool {
	if v == nil {
		return true
	}

	if nb, ok := v.(Nilable); ok {
		return nb.IsNil()
	}

	switch reflect.TypeOf(v).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(v).IsNil()
	}
	return false
}

func HasNil(vs ...interface{}) bool {
	for i := range vs {
		if IsNil(vs[i]) {
			return true
		}
	}
	return false
}

func IsStringKeyMap(v interface{}) bool {
	t := reflect.TypeOf(v)
	return t != nil && t.Kind() == reflect.Map && t.Key().Kind() == reflect.String
}

func StructToMap(val interface{}, keys ...string) (map[string]interface{}, error) {
	v := reflect.Indirect(reflect.ValueOf(val))

	if v.Kind() != reflect.Struct { // including !v.IsValid() where v.Kind() == reflect.Invalid
		return nil, fmt.Errorf("Not a structure: %v", val)
	}

	ret := map[string]interface{}{}
	if len(keys) == 0 {
		for i, n := 0, v.NumField(); i < n; i++ {
			name := v.Type().Field(i).Name
			f := v.Field(i)
			if f.IsValid() {
				ret[name] = f.Interface()
			}
		}
	} else {
		for _, key := range keys {
			f := v.FieldByName(key)
			if f.IsValid() {
				ret[key] = f.Interface()
			}
		}
	}
	return ret, nil
}

// assume v valid, will deferenced once
func isBasic(v reflect.Value, noPtr bool) bool {
	switch reflect.TypeOf(v).Kind() {
	case reflect.Ptr: // TODO: remove this
		if !noPtr {
			val := v.Elem().Interface()
			return isBasic(reflect.ValueOf(val), true)
		}
	case reflect.Bool, reflect.String, reflect.Int, reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
		return true
	}
	return false
}

// will deference twice... TODO: remove this
func IsBasic(val interface{}) bool {
	v := reflect.Indirect(reflect.ValueOf(val))
	return v.IsValid() && isBasic(v, false)
}

func StructToBasicMap(val interface{}) map[string]interface{} {
	v := reflect.Indirect(reflect.ValueOf(val))

	if v.Kind() != reflect.Struct {
		return nil
	}

	ret := map[string]interface{}{}
	for i, n := 0, v.NumField(); i < n; i++ {
		key := v.Type().Field(i).Name
		f := v.Field(i)
		if f.IsValid() {
			switch reflect.TypeOf(v).Kind() {
			case reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
			case reflect.Ptr:
				ret[key] = f.Elem().Interface()
			default:
				ret[key] = f.Interface()
			}
		}
	}
	return ret
}

func StructFieldNames(val interface{}) []string {
	v := reflect.Indirect(reflect.ValueOf(val))

	if v.Kind() != reflect.Struct {
		return nil
	}

	ret := []string{}
	for i, n := 0, v.NumField(); i < n; i++ {
		ret = append(ret, v.Type().Field(i).Name)
	}
	return ret
}
