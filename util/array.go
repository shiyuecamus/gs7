// Copyright 2024 shiyuecamus. All Rights Reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package util

import (
	"reflect"
)

// Int interface type
type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Uint interface type
type Uint interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Xint interface type. alias of Integer
type Xint interface {
	Int | Uint
}

// Integer interface type. all int or uint types
type Integer interface {
	Int | Uint
}

// Float interface type
type Float interface {
	~float32 | ~float64
}

// IntOrFloat interface type. all int and float types
type IntOrFloat interface {
	Int | Float
}

// XintOrFloat interface type. all int, uint and float types
type XintOrFloat interface {
	Int | Uint | Float
}

// SortedType interface type. same of constraints.Ordered
//
// it can be ordered, that supports the operators < <= >= >.
//
// contains: (x)int, float, ~string types
type SortedType interface {
	Int | Uint | Float | ~string
}

// Compared type. alias of constraints.SortedType
// type Compared = SortedType
type Compared interface {
	Int | Uint | Float | ~string
}

// SimpleType interface type. alias of ScalarType
//
// contains: (x)int, float, ~string, ~bool types
type SimpleType interface {
	Int | Uint | Float | ~string | ~bool
}

// ScalarType interface type.
//
// TIP: has bool type, it cannot be ordered
//
// contains: (x)int, float, ~string, ~bool types
type ScalarType interface {
	Int | Uint | Float | ~string | ~bool
}

//
//
// Matcher type
//
//

// Matcher interface
type Matcher[T any] interface {
	Match(s T) bool
}

// MatchFunc definition. implements Matcher interface
type MatchFunc[T any] func(v T) bool

// Match satisfies the Matcher interface
func (fn MatchFunc[T]) Match(v T) bool {
	return fn(v)
}

// StringMatcher interface
type StringMatcher interface {
	Match(s string) bool
}

// StringMatchFunc definition
type StringMatchFunc func(s string) bool

// Match satisfies the StringMatcher interface
func (fn StringMatchFunc) Match(s string) bool {
	return fn(s)
}

// StringHandler interface
type StringHandler interface {
	Handle(s string) string
}

// StringHandleFunc definition
type StringHandleFunc func(s string) string

// Handle satisfies the StringHandler interface
func (fn StringHandleFunc) Handle(s string) string {
	return fn(s)
}

// Reverse any T slice.
//
// eg: []string{"site", "user", "info", "0"} -> []string{"0", "info", "user", "site"}
func Reverse[T any](ls []T) {
	ln := len(ls)
	for i := 0; i < ln/2; i++ {
		li := ln - i - 1
		ls[i], ls[li] = ls[li], ls[i]
	}
}

// Remove give element from slice []T.
//
// eg: []string{"site", "user", "info", "0"} -> []string{"site", "user", "info"}
func Remove[T Compared](ls []T, val T) []T {
	return Filter(ls, func(el T) bool {
		return el != val
	})
}

// Filter given slice, default will filter zero value.
//
// Usage:
//
//	// output: [a, b]
//	ss := arrutil.Filter([]string{"a", "", "b", ""})
func Filter[T any](ls []T, filter ...MatchFunc[T]) []T {
	var fn MatchFunc[T]
	if len(filter) > 0 && filter[0] != nil {
		fn = filter[0]
	} else {
		fn = func(el T) bool {
			return !reflect.ValueOf(el).IsZero()
		}
	}

	newLs := make([]T, 0, len(ls))
	for _, el := range ls {
		if fn(el) {
			newLs = append(newLs, el)
		}
	}
	return newLs
}

// MapFn map handle function type.
type MapFn[T any, V any] func(input T) (target V, find bool)

// Map a list to new list
//
// eg: mapping [object0{},object1{},...] to flatten list [object0.someKey, object1.someKey, ...]
func Map[T any, V any](list []T, mapFn MapFn[T, V]) []V {
	flatArr := make([]V, 0, len(list))

	for _, obj := range list {
		if target, ok := mapFn(obj); ok {
			flatArr = append(flatArr, target)
		}
	}
	return flatArr
}

// Column alias of Map func
func Column[T any, V any](list []T, mapFn func(obj T) (val V, find bool)) []V {
	return Map(list, mapFn)
}

// Unique value in the given slice data.
func Unique[T ~string | XintOrFloat](list []T) []T {
	if len(list) < 2 {
		return list
	}

	valMap := make(map[T]struct{}, len(list))
	uniArr := make([]T, 0, len(list))

	for _, t := range list {
		if _, ok := valMap[t]; !ok {
			valMap[t] = struct{}{}
			uniArr = append(uniArr, t)
		}
	}
	return uniArr
}

// IndexOf value in given slice.
func IndexOf[T ~string | XintOrFloat](val T, list []T) int {
	for i, v := range list {
		if v == val {
			return i
		}
	}
	return -1
}
