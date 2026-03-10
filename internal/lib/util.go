package lib

import (
	"reflect"
)

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	v := reflect.ValueOf(i)
	if v.Kind() == reflect.Ptr {
		return v.IsNil()
	}
	return false
}

func Pointer[T any](v T) *T {
	return &v
}

func UniqueStringSlice(list []string) []string {
	seen := make(map[string]struct{})
	var result []string

	for _, v := range list {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		result = append(result, v)
	}
	return result
}
