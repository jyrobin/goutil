// Copyright (c) 2021 Jing-Ying Chen. Subject to the MIT License.

package goutil

import "testing"

func TestContainsRecursive(t *testing.T) {
	var tests = []struct {
		v1 interface{}
		v2 interface{}
		ok bool
	}{
		{1, 1, true},
		{1, 2, false},
		{[]int{1, 2, 3}, []int{1, 2, 3}, true},
		{[]int{1, 2, 3}, []int{1, 2}, false},
		{
			map[string]int{"a": 1, "b": 2},
			map[string]int{"b": 2, "a": 1},
			true,
		},
		{
			map[string]int{"a": 1, "b": 2},
			map[string]int{"b": 1, "a": 1},
			false,
		},
		{
			map[string]int{"a": 1, "b": 2},
			map[string]int{"a": 1, "c": 3},
			false,
		},
	}

	for _, tt := range tests {
		if ContainsRecursive(tt.v1, tt.v2) != tt.ok {
			if tt.ok {
				t.Fatalf("%+v should contain %+v", tt.v1, tt.v2)
			} else {
				t.Fatalf("%+v should not contain %+v", tt.v1, tt.v2)
			}
		}
	}
}

func TestJsonContains(t *testing.T) {
	var tests = []struct {
		v1 string
		v2 string
		ok bool
	}{
		{"1", "1", true},
		{"[1, 2, 3]", "[1, 2, 3]", true},
		{"[1, 2, 3]", "[1, 2]", false},
		{`{"a":1,"b":2}`, `{"b":2,"a":1}`, true},
		{`{"a":1,"b":2}`, `{"b":2}`, true},
		{`{"a":1,"b":2}`, `{"b":1,"a":1}`, false},
		{`{"a":1,"b":2}`, `{"c":2,"a":1}`, false},
	}

	for _, tt := range tests {
		if JsonContains([]byte(tt.v1), []byte(tt.v2)) != tt.ok {
			if tt.ok {
				t.Fatalf("%s should contain %s", tt.v1, tt.v2)
			} else {
				t.Fatalf("%s should not contain %s", tt.v1, tt.v2)
			}
		}
	}
}
