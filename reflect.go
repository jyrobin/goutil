// Copyright (c) 2021 Jing-Ying Chen. Subject to the MIT License.

package goutil

import (
	"fmt"
	"reflect"
)

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
