package goutil

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// REF - https://stackoverflow.com/questions/56616196/how-to-convert-camel-case-string-to-snake-case/56616250

// ASCII only

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
var matchLodash = regexp.MustCompile("[_]+")
var matchDash = regexp.MustCompile("[-]+")
var matchSpace = regexp.MustCompile("\\s+")

func slugify(str, f string) string {
	slug := matchFirstCap.ReplaceAllString(str, f)
	slug = matchAllCap.ReplaceAllString(slug, f)
	return strings.ToLower(slug)
}

func Slugify(str, sep string) string {
	return slugify(str, fmt.Sprintf("${1}%s${2}", sep))
}

func CamelToSnake(str string) string {
	return slugify(str, "${1}_${2}")
}
func CamelToSlug(str string) string {
	return slugify(str, "${1}-${2}")
}

func SnakeToCamel(str string, capital bool) string {
	return camelize(str, capital, matchLodash)
}
func SlugToCamel(str string, capital bool) string {
	return camelize(str, capital, matchDash)
}

func camelize(str string, capital bool, re *regexp.Regexp) string {
	s := re.ReplaceAllString(str, " ")
	s = strings.TrimSpace(s)
	s = strings.Title(s)
	s = matchSpace.ReplaceAllString(s, "")
	if len(s) > 0 {
		r := []rune(s)
		if capital && !unicode.IsUpper(r[0]) {
			r[0] = unicode.ToUpper(r[0])
		} else if !capital && unicode.IsUpper(r[0]) {
			r[0] = unicode.ToLower(r[0])
		}
		s = string(r)
	}
	return s
}
