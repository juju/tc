// Copied with small adaptations from the reflect package in the
// Go source tree.

// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE-golang file.

package tc

import (
	"fmt"
	"math/big"
	"reflect"
	"strings"
	"time"
	"unsafe"
)

var timeType = reflect.TypeOf(time.Time{})

// During deepValueEqual, must keep track of checks that are
// in progress.  The comparison algorithm assumes that all
// checks in progress are true when it reencounters them.
// Visited comparisons are stored in a map indexed by visit.
type visit struct {
	a1  uintptr
	a2  uintptr
	typ reflect.Type
}

type mismatchError struct {
	v1, v2 reflect.Value
	path   string
	how    string
}

func (err *mismatchError) Error() string {
	path := err.path
	if path == "" {
		path = "top level"
	}
	return fmt.Sprintf("mismatch at %s: %s; obtained %#v; expected %#v", path, err.how, printable(err.v1), printable(err.v2))
}

func printable(v reflect.Value) any {
	vi := interfaceOf(v)
	switch vi := vi.(type) {
	case time.Time:
		return vi.UTC().Format(time.RFC3339Nano)
	default:
		return vi
	}
}

// Tests for deep equality using reflected types. The map argument tracks
// comparisons that have already been seen, which allows short circuiting on
// recursive types.
func deepValueEqual(
	path string, v1, v2 reflect.Value, visited map[visit]bool, depth int,
	equal func(a, b any) bool,
	length func(path string, a, b int) bool,
	customCheckFunc CustomCheckFunc,
) (ok bool, err error) {
	errorf := func(f string, a ...any) error {
		return &mismatchError{
			v1:   v1,
			v2:   v2,
			path: strings.Replace(path, topLevel, "", 1),
			how:  fmt.Sprintf(f, a...),
		}
	}
	if !v1.IsValid() {
		if !v2.IsValid() {
			return true, nil
		}
		return false, errorf("validity mismatch")
	} else if !v2.IsValid() {
		v2 = reflect.Zero(v1.Type())
	}

	if v1.CanAddr() && v2.CanAddr() {
		addr1 := v1.UnsafeAddr()
		addr2 := v2.UnsafeAddr()
		if addr1 > addr2 {
			// Canonicalize order to reduce number of entries in visited.
			addr1, addr2 = addr2, addr1
		}
		// Avoid infinite recursion.
		typ := v1.Type()
		v := visit{addr1, addr2, typ}
		if visited[v] {
			return true, nil
		}
		visited[v] = true
	}

	if customCheckFunc != nil {
		useDefault, equal, err := customCheckFunc(path, interfaceOf(v1), interfaceOf(v2))
		if !useDefault {
			return equal, err
		}
	}

	if v1.Type() != v2.Type() {
		return false, errorf("type mismatch %s vs %s", v1.Type(), v2.Type())
	}

	if v1.CanInterface() && v2.CanInterface() {
		switch v1.Type() {
		case reflect.TypeFor[*big.Int]():
			if bigInt1, ok := v1.Interface().(*big.Int); ok {
				if bigInt2, ok := v2.Interface().(*big.Int); ok {
					if equal(bigInt1.Cmp(bigInt2), 0) {
						return true, nil
					} else {
						return false, errorf("unequal big int")
					}
				}
			}
		case reflect.TypeFor[*big.Float]():
			if bigFloat1, ok := v1.Interface().(*big.Float); ok {
				if bigFloat2, ok := v2.Interface().(*big.Float); ok {
					if equal(bigFloat1.Cmp(bigFloat2), 0) {
						return true, nil
					} else {
						return false, errorf("unequal big float")
					}
				}
			}
		case reflect.TypeFor[*big.Rat]():
			if bigRat1, ok := v1.Interface().(*big.Rat); ok {
				if bigRat2, ok := v2.Interface().(*big.Rat); ok {
					if equal(bigRat1.Cmp(bigRat2), 0) {
						return true, nil
					} else {
						return false, errorf("unequal big rational")
					}
				}
			}
		}
	}

	switch v1.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v1.Len(); i++ {
			rhsValue := reflect.Zero(v1.Type())
			if i < v2.Len() {
				rhsValue = v2.Index(i)
			}
			ok, err := deepValueEqual(
				fmt.Sprintf("%s[%d]", path, i),
				v1.Index(i), rhsValue,
				visited, depth+1,
				equal, length, customCheckFunc)
			if !ok {
				return false, err
			}
		}
		if !length(path, v1.Len(), v2.Len()) {
			return false, errorf("slice/array length mismatch, %d vs %d",
				v1.Len(), v2.Len())
		}
		return true, nil
	case reflect.Interface:
		return deepValueEqual(path, v1.Elem(), v2.Elem(),
			visited, depth+1, equal, length, customCheckFunc)
	case reflect.Ptr:
		return deepValueEqual("(*"+path+")", v1.Elem(), v2.Elem(),
			visited, depth+1, equal, length, customCheckFunc)
	case reflect.Struct:
		if v1.Type() == timeType {
			// Special case for time - we ignore the time zone.
			t1 := interfaceOf(v1).(time.Time)
			t2 := interfaceOf(v2).(time.Time)
			if equal(t1.Equal(t2), true) {
				return true, nil
			}
			return false, errorf("unequal")
		}
		for i, n := 0, v1.NumField(); i < n; i++ {
			path := path + "." + v1.Type().Field(i).Name
			ok, err := deepValueEqual(path, v1.Field(i), v2.Field(i),
				visited, depth+1,
				equal, length, customCheckFunc)
			if !ok {
				return false, err
			}
		}
		return true, nil
	case reflect.Map:
		for _, k := range v1.MapKeys() {
			var p string
			if k.CanInterface() {
				p = path + "[" + fmt.Sprintf("%#v", k.Interface()) + "]"
			} else {
				p = path + "[someKey]"
			}
			ok, err := deepValueEqual(p, v1.MapIndex(k), v2.MapIndex(k),
				visited, depth+1,
				equal, length, customCheckFunc)
			if !ok {
				return false, err
			}
		}
		if !length(path, v1.Len(), v2.Len()) {
			return false, errorf("map length mismatch, %d vs %d",
				v1.Len(), v2.Len())
		}
		return true, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if !equal(v1.Int(), v2.Int()) {
			return false, errorf("unequal")
		}
		return true, nil
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if !equal(v1.Uint(), v2.Uint()) {
			return false, errorf("unequal")
		}
		return true, nil
	case reflect.Float32, reflect.Float64:
		if !equal(v1.Float(), v2.Float()) {
			return false, errorf("unequal")
		}
		return true, nil
	case reflect.Complex64, reflect.Complex128:
		if !equal(v1.Complex(), v2.Complex()) {
			return false, errorf("unequal")
		}
		return true, nil
	case reflect.Bool:
		if !equal(v1.Bool(), v2.Bool()) {
			return false, errorf("unequal")
		}
		return true, nil
	case reflect.String:
		if !equal(v1.String(), v2.String()) {
			return false, errorf("unequal")
		}
		return true, nil
	case reflect.Chan, reflect.UnsafePointer:
		if !equal(v1.Pointer(), v2.Pointer()) {
			return false, errorf("unequal")
		}
		return true, nil
	case reflect.Func:
		if v1.IsNil() && equal(true, v2.IsNil()) {
			return true, nil
		}
		if equal(v1.Pointer(), v2.Pointer()) {
			return true, nil
		}
		// Can't do better than this:
		return false, errorf("non-nil functions")
	default:
		panic("unexpected type " + v1.Type().String())
	}
}

// DeepEqual tests for deep equality. It uses normal == equality where
// possible but will scan elements of arrays, slices, maps, and fields
// of structs. In maps, keys are compared with == but elements use deep
// equality. DeepEqual correctly handles recursive types. Functions are
// equal only if they are both nil.
//
// DeepEqual differs from reflect.DeepEqual in two ways:
// - an empty slice is considered equal to a nil slice.
// - two time.Time values that represent the same instant
// but with different time zones are considered equal.
//
// If the two values compare unequal, the resulting error holds the
// first difference encountered.
func DeepEqual(a1, a2 any) (bool, error) {
	errorf := func(f string, a ...any) error {
		return &mismatchError{
			v1:   reflect.ValueOf(a1),
			v2:   reflect.ValueOf(a2),
			path: "",
			how:  fmt.Sprintf(f, a...),
		}
	}
	if a1 == nil || a2 == nil {
		if a1 == a2 {
			return true, nil
		}
		return false, errorf("nil vs non-nil mismatch")
	}
	v1 := reflect.ValueOf(a1)
	v2 := reflect.ValueOf(a2)
	if v1.Type() != v2.Type() {
		return false, errorf("type mismatch %s vs %s", v1.Type(), v2.Type())
	}
	return deepValueEqual(topLevel, v1, v2, make(map[visit]bool), 0,
		func(a, b any) bool { return a == b },
		func(path string, a, b int) bool { return a == b },
		nil)
}

// CustomCheckFunc should return true for useDefault if DeepEqualWithCustomCheck should behave like DeepEqual.
// Otherwise the result of the CustomCheckFunc is used.
type CustomCheckFunc func(path string, a1 any, a2 any) (useDefault bool, equal bool, err error)

// interfaceOf returns v.Interface() even if v.CanInterface() == false.
// This enables us to call fmt.Printf on a value even if it's derived
// from inside an unexported field.
// See https://code.google.com/p/go/issues/detail?id=8965
// for a possible future alternative to this hack.
func interfaceOf(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}
	return bypassCanInterface(v).Interface()
}

type flagField uintptr

var flagValOffset = func() uintptr {
	field, ok := reflect.TypeOf(reflect.Value{}).FieldByName("flag")
	if !ok {
		panic("reflect.Value has no flag field")
	}
	return field.Offset
}()

func toFlagField(v *reflect.Value) *flagField {
	return (*flagField)(unsafe.Pointer(uintptr(unsafe.Pointer(v)) + flagValOffset))
}

// bypassCanInterface returns a version of v that
// bypasses the CanInterface check.
func bypassCanInterface(v reflect.Value) reflect.Value {
	if !v.IsValid() || v.CanInterface() {
		return v
	}
	*toFlagField(&v) &^= flagRO
	return v
}

// Sanity checks against future reflect package changes
// to the type or semantics of the Value.flag field.
func init() {
	field, ok := reflect.TypeOf(reflect.Value{}).FieldByName("flag")
	if !ok {
		panic("reflect.Value has no flag field")
	}
	if field.Type.Kind() != reflect.TypeOf(flagField(0)).Kind() {
		panic("reflect.Value flag field has changed kind")
	}
}
