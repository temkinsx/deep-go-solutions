package main

import (
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go test -v generic_test.go

type Person struct {
	Name    string `properties:"name"`
	Address string `properties:"address,omitempty"`
	Age     int    `properties:"age"`
	Married bool   `properties:"married"`
}

func Serialize[T any](v T) string {
	rv := reflect.ValueOf(v)
	if !rv.IsValid() {
		return ""
	}

	if rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return ""
		}
		rv = rv.Elem()
	}

	if rv.Kind() != reflect.Struct {
		return ""
	}

	t := rv.Type()
	first := true
	var b strings.Builder

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if field.PkgPath != "" {
			continue
		}

		tag := field.Tag.Get("properties")
		fieldName := ""
		omitEmpty := false

		if tag != "" {
			parts := strings.Split(tag, ",")
			fieldName = strings.TrimSpace(parts[0])
			if fieldName == "-" {
				continue
			}

			for _, p := range parts[1:] {
				p = strings.TrimSpace(p)
				if p == "omitempty" {
					omitEmpty = true
				}
			}
		}

		if fieldName == "" {
			fieldName = strings.ToLower(field.Name)
		}

		fv := rv.Field(i)

		if omitEmpty && fv.IsZero() {
			continue
		}

		if !first {
			b.WriteByte('\n')
		}
		first = false

		fmt.Fprintf(&b, "%s=%v", fieldName, fv.Interface())

	}
	return b.String()
}

func TestSerialization(t *testing.T) {
	tests := map[string]struct {
		person Person
		result string
	}{
		"test case with empty fields": {
			result: "name=\nage=0\nmarried=false",
		},
		"test case with fields": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
			},
			result: "name=John Doe\nage=30\nmarried=true",
		},
		"test case with omitempty field": {
			person: Person{
				Name:    "John Doe",
				Age:     30,
				Married: true,
				Address: "Paris",
			},
			result: "name=John Doe\naddress=Paris\nage=30\nmarried=true",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			result := Serialize(test.person)
			assert.Equal(t, test.result, result)
		})
	}
}
