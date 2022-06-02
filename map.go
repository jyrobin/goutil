// Copyright (c) 2021 Jing-Ying Chen. Subject to the MIT License.

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

func OmitMap(m Map, skips ...string) Map {
	b := len(skips) == 0
	ret := Map{}
	for k, v := range m {
		if b || !ContainsString(skips, k) {
			ret[k] = v
		}

	}
	return ret
}

func OmitStrMap(m StrMap, skips ...string) StrMap {
	b := len(skips) == 0
	ret := StrMap{}
	for k, v := range m {
		if b || !ContainsString(skips, k) {
			ret[k] = v
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

func MapContainsBasic(m, n Map) bool {
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

func DefaultParamListString(args ...string) string {
	return ParamListString(',', '=', args...)
}

// no checking, use with care
func ParamListString(sep, keySep rune, args ...string) string {
	var sb strings.Builder
	for i, n := 0, len(args); i < n; i += 2 {
		key, val := args[i], args[i+1]
		key = CutLeft(key, sep, keySep)
		val = CutLeft(val, sep, keySep)
		if i > 0 {
			sb.WriteRune(sep)
		}
		sb.WriteString(key)
		sb.WriteRune(keySep)
		sb.WriteString(val)
	}
	return sb.String()
}

func Params(args ...string) map[string]string {
	ret := map[string]string{}
	for i, n := 0, len(args); i < n; i += 2 {
		ret[args[i]] = args[i+1]
	}
	return ret
}

func ParseParams(params string, sep string, keySeps ...string) map[string]string {
	ret := map[string]string{}
	if params == "" {
		return ret
	}

	keySep := "="
	if len(keySeps) > 0 {
		keySep = keySeps[0]
	}

	words := strings.Split(strings.TrimSpace(params), sep)
	for _, word := range words {
		pair := strings.SplitN(word, keySep, 2)
		key := strings.TrimSpace(pair[0])
		var val string
		if len(pair) == 2 {
			val = strings.TrimSpace(pair[1])
		}
		ret[key] = val
	}
	return ret
}

// args[0] the separator, args[1] the key-value separator
func ParamsString(params map[string]string, args ...string) string {
	np := len(params)
	if np == 0 {
		return ""
	}

	sep := ";"
	keySep := "="
	if n := len(args); n > 1 {
		sep, keySep = args[0], args[1]
	} else if n == 1 {
		sep = args[0]
	}

	var sb strings.Builder
	for k, v := range params {
		if sb.Len() > 0 {
			sb.WriteString(sep)
		}
		sb.WriteString(k)
		sb.WriteString(keySep)
		sb.WriteString(v)
	}
	return sb.String()
}

func MergeStrMap(dest, src map[string]string, srcs ...map[string]string) map[string]string {
	if dest == nil {
		dest = map[string]string{}
	}
	for k, v := range src {
		dest[k] = v
	}
	for i := range srcs {
		for k, v := range srcs[i] {
			dest[k] = v
		}
	}
	return dest
}
