// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

// Must calls f with no arguments and asserts that the error return is nil
// and returns the remaining single return value.
// Must is an alias of [Must0_1].
func Must[T any, E error, F func() (T, E)](c LikeC, f F) T {
	c.Helper()
	return Must0_1(c, f)
}

// Must0 calls f with no arguments and asserts that the error return is nil
// and returns the remaining single return value.
// Must0 is an alias of [Must0_1].
func Must0[T any, E error, F func() (T, E)](c LikeC, f F) T {
	c.Helper()
	return Must0_1(c, f)
}

// Must1 calls f with one argument and asserts that the error return is nil
// and returns the remaining single return value.
// Must1 is an alias of [Must1_1].
func Must1[T any, E error, A any, F func(A) (T, E)](c LikeC, f F, a A) T {
	c.Helper()
	return Must1_1(c, f, a)
}

// Must2 calls f with two arguments and asserts that the error return is nil
// and returns the remaining single return value.
// Must2 is an alias of [Must2_1].
func Must2[T any, E error, A any, B any, F func(A, B) (T, E)](c LikeC, f F, a A, b B) T {
	c.Helper()
	return Must2_1(c, f, a, b)
}

// Must0_0 calls f with no arguments and asserts that the error return is nil
// and returns no values.
func Must0_0[E error, F func() E](c LikeC, f F) {
	c.Helper()
	err := f()
	c.Assert(err, ErrorIsNil)
}

// Must1_0 calls f with one argument and asserts that the error return is nil
// and returns no values.
func Must1_0[E error, A any, F func(A) E](c LikeC, f F, a A) {
	c.Helper()
	err := f(a)
	c.Assert(err, ErrorIsNil)
}

// Must2_0 calls f with two arguments and asserts that the error return is nil
// and returns no values.
func Must2_0[E error, A any, B any, F func(A, B) E](c LikeC, f F, a A, b B) {
	c.Helper()
	err := f(a, b)
	c.Assert(err, ErrorIsNil)
}

// Must0_1 calls f with no arguments and asserts that the error return is nil
// and returns the remaining single return value.
func Must0_1[T any, E error, F func() (T, E)](c LikeC, f F) T {
	c.Helper()
	r, err := f()
	c.Assert(err, ErrorIsNil)
	return r
}

// Must1_1 calls f with one argument and asserts that the error return is nil
// and returns the remaining single return value.
func Must1_1[T any, E error, A any, F func(A) (T, E)](c LikeC, f F, a A) T {
	c.Helper()
	r, err := f(a)
	c.Assert(err, ErrorIsNil)
	return r
}

// Must2_1 calls f with two arguments and asserts that the error return is nil
// and returns the remaining single return value.
func Must2_1[T any, E error, A any, B any, F func(A, B) (T, E)](c LikeC, f F, a A, b B) T {
	c.Helper()
	r, err := f(a, b)
	c.Assert(err, ErrorIsNil)
	return r
}

// Must0_2 calls f with no arguments and asserts that the error return is nil
// and returns the remaining two return values.
func Must0_2[T any, T2 any, E error, F func() (T, T2, E)](c LikeC, f F) (T, T2) {
	c.Helper()
	r1, r2, err := f()
	c.Assert(err, ErrorIsNil)
	return r1, r2
}

// Must1_2 calls f with one argument and asserts that the error return is nil
// and returns the remaining two return values.
func Must1_2[T any, T2 any, E error, A any, F func(A) (T, T2, E)](c LikeC, f F, a A) (T, T2) {
	c.Helper()
	r1, r2, err := f(a)
	c.Assert(err, ErrorIsNil)
	return r1, r2
}

// Must2_2 calls f with two arguments and asserts that the error return is nil
// and returns the remaining two return values.
func Must2_2[T any, T2 any, E error, A any, B any, F func(A, B) (T, T2, E)](c LikeC, f F, a A, b B) (T, T2) {
	c.Helper()
	r1, r2, err := f(a, b)
	c.Assert(err, ErrorIsNil)
	return r1, r2
}
