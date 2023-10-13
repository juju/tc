// Copyright 2020 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

//go:build go1.13 && !go1.12 && !go1.14
// +build go1.13,!go1.12,!go1.14

package check

const (
	flagRO flagField = 1<<5 | 1<<6
)
