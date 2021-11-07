// Copyright (c) 2021 Jing-Ying Chen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package goutil

import (
	"fmt"
	"reflect"
	"strings"
)

// Only concern string-keyed maps here

type Map map[string]interface{}
type StrMap map[string]string

func PickMap(m Map, keys ...string) Map {
	ret := Map{}
	for _, key := range keys {
		if val, ok := m[key]; ok {
			ret[key] = val
		}
	}
	return ret
}

func PickStrMap(m StrMap, keys ...string) StrMap {
	ret := StrMap{}
	for _, key := range keys {
		if val, ok := m[key]; ok {
			ret[key] = val
		}
	}
	return ret
}

func MapAllKeysExist(m Map, keys ...string) bool {
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			return false
		}
	}
	return true
}

func StrMapAllKeysExist(m StrMap, keys ...string) bool {
	for _, key := range keys {
		if _, ok := m[key]; !ok {
			return false
		}
	}
	return true
}

func StrMapAllNonEmpty(m StrMap, keys ...string) bool {
	for _, key := range keys {
		if val, ok := m[key]; !ok || strings.TrimSpace(val) == "" {
			return false
		}
	}
	return true
}

func StrMapKeys(m StrMap) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func MapKeys(m Map) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func StrMapToMap(m StrMap) Map {
	ret := Map{}
	for k, v := range m {
		ret[k] = v
	}
	return ret
}

func StrMapContains(m, n StrMap) bool {
	for k, v := range n {
		if val, ok := m[k]; !ok || val != v {
			return false
		}
	}
	return true
}

func MapContains(m, n Map) bool {
	for k, v := range n {
		if val, ok := m[k]; !ok || val != v {
			return false
		}
	}
	return true
}

// GetMapKeys return keys only for string-keyed map
func GetMapKeys(val interface{}) ([]string, error) {
	switch v := val.(type) { // let me optimize pre-maturally
	case map[string]interface{}:
	case Map:
		return MapKeys(v), nil
	case map[string]string:
	case StrMap:
		return StrMapKeys(v), nil
	}

	if !IsStringKeyMap(val) {
		return nil, fmt.Errorf("Not a string-keyed map")
	}

	keys := reflect.ValueOf(val).MapKeys()
	ret := make([]string, len(keys))
	for i, key := range keys {
		ret[i] = key.Interface().(string)
	}
	return ret, nil
}

// only concerns string keyed map sub-structure
// - slices, pointers, structs, etc are compared via deepEqual
func ContainsRecursive(v1, v2 interface{}) bool {
	if reflect.DeepEqual(v1, v2) {
		return true
	}

	t1 := reflect.TypeOf(v1)
	if t1 == nil || t1.Kind() != reflect.Map || t1.Key().Kind() != reflect.String {
		return false
	}
	t2 := reflect.TypeOf(v2)
	if t2 == nil || t2.Kind() != reflect.Map || t2.Key().Kind() != reflect.String {
		return false
	}

	vv1, vv2 := reflect.ValueOf(v1), reflect.ValueOf(v2)
	for _, key := range vv2.MapKeys() {
		val1, val2 := vv1.MapIndex(key), vv2.MapIndex(key)
		if !val1.IsValid() {
			return false
		}
		if !ContainsRecursive(val1.Interface(), val2.Interface()) {
			return false
		}
	}
	return true
}
