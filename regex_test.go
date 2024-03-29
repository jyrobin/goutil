// Copyright (c) 2021 Jing-Ying Chen. Subject to the MIT License.

package goutil

import "testing"

func TestCamelToSnakeAndSlug(t *testing.T) {
	var tests = []struct {
		s1 string
		s2 string
		s3 string
	}{
		{"abcDefGhi", "abc_def_ghi", "abc-def-ghi"},
		{"AbcDefGhi", "abc_def_ghi", "abc-def-ghi"},
		{"ABCDefGhi", "abc_def_ghi", "abc-def-ghi"},
		{"AbcDEFGhi", "abc_def_ghi", "abc-def-ghi"},
	}

	for _, tt := range tests {
		s := CamelToSnake(tt.s1)
		if s != tt.s2 {
			t.Fatalf("%s should become %s, got %s", tt.s1, tt.s2, s)
		}
		s = CamelToSlug(tt.s1)
		if s != tt.s3 {
			t.Fatalf("%s should become %s, got %s", tt.s1, tt.s3, s)
		}
	}
}

func TestCamelize(t *testing.T) {
	var tests = []struct {
		s1 string
		s2 string
		s3 string
		s4 string
	}{
		{"abc_def_ghi", "abc-def-ghi", "AbcDefGhi", "abcDefGhi"},
		{"abc__def___ghi", "abc--def---ghi", "AbcDefGhi", "abcDefGhi"},
		{"_abc__def___ghi__", "-abc--def---ghi--", "AbcDefGhi", "abcDefGhi"},
	}

	for _, tt := range tests {
		s := SnakeToCamel(tt.s1, true)
		if s != tt.s3 {
			t.Fatalf("%s should become %s, got %s", tt.s1, tt.s3, s)
		}
		s = SnakeToCamel(tt.s1, false)
		if s != tt.s4 {
			t.Fatalf("%s should become %s, got %s", tt.s1, tt.s4, s)
		}
		s = SlugToCamel(tt.s2, true)
		if s != tt.s3 {
			t.Fatalf("%s should become %s, got %s", tt.s2, tt.s3, s)
		}
		s = SlugToCamel(tt.s2, false)
		if s != tt.s4 {
			t.Fatalf("%s should become %s, got %s", tt.s2, tt.s4, s)
		}
	}
}
