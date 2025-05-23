// Copyright 2013 Canonical Ltd.
// Copyright 2014 Cloudbase Solutions SRL
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
)

// IsNonEmptyFile checker

type isNonEmptyFileChecker struct {
	*CheckerInfo
}

var IsNonEmptyFile Checker = &isNonEmptyFileChecker{
	&CheckerInfo{Name: "IsNonEmptyFile", Params: []string{"obtained"}},
}

func (checker *isNonEmptyFileChecker) Check(params []any, names []string) (result bool, error string) {
	filename, isString := stringOrStringer(params[0])
	if isString {
		fileInfo, err := os.Stat(filename)
		if os.IsNotExist(err) {
			return false, fmt.Sprintf("%s does not exist", filename)
		} else if err != nil {
			return false, fmt.Sprintf("other stat error: %v", err)
		}
		if fileInfo.Size() > 0 {
			return true, ""
		} else {
			return false, fmt.Sprintf("%s is empty", filename)
		}
	}

	value := reflect.ValueOf(params[0])
	return false, fmt.Sprintf("obtained value is not a string and has no .String(), %s:%#v", value.Kind(), params[0])
}

// IsDirectory checker

type isDirectoryChecker struct {
	*CheckerInfo
}

var IsDirectory Checker = &isDirectoryChecker{
	&CheckerInfo{Name: "IsDirectory", Params: []string{"obtained"}},
}

func (checker *isDirectoryChecker) Check(params []any, names []string) (result bool, error string) {
	path, isString := stringOrStringer(params[0])
	if isString {
		fileInfo, err := os.Stat(path)
		if os.IsNotExist(err) {
			return false, fmt.Sprintf("%s does not exist", path)
		} else if err != nil {
			return false, fmt.Sprintf("other stat error: %v", err)
		}
		if fileInfo.IsDir() {
			return true, ""
		} else {
			return false, fmt.Sprintf("%s is not a directory", path)
		}
	}

	value := reflect.ValueOf(params[0])
	return false, fmt.Sprintf("obtained value is not a string and has no .String(), %s:%#v", value.Kind(), params[0])
}

// IsSymlink checker

type isSymlinkChecker struct {
	*CheckerInfo
}

var IsSymlink Checker = &isSymlinkChecker{
	&CheckerInfo{Name: "IsSymlink", Params: []string{"obtained"}},
}

func (checker *isSymlinkChecker) Check(params []any, names []string) (result bool, error string) {
	path, isString := stringOrStringer(params[0])
	if isString {
		fileInfo, err := os.Lstat(path)
		if os.IsNotExist(err) {
			return false, fmt.Sprintf("%s does not exist", path)
		} else if err != nil {
			return false, fmt.Sprintf("other stat error: %v", err)
		}
		if fileInfo.Mode()&os.ModeSymlink != 0 {
			return true, ""
		} else {
			return false, fmt.Sprintf("%s is not a symlink: %+v", path, fileInfo)
		}
	}

	value := reflect.ValueOf(params[0])
	return false, fmt.Sprintf("obtained value is not a string and has no .String(), %s:%#v", value.Kind(), params[0])
}

// DoesNotExist checker makes sure the path specified doesn't exist.

type doesNotExistChecker struct {
	*CheckerInfo
}

var DoesNotExist Checker = &doesNotExistChecker{
	&CheckerInfo{Name: "DoesNotExist", Params: []string{"obtained"}},
}

func (checker *doesNotExistChecker) Check(params []any, names []string) (result bool, error string) {
	path, isString := stringOrStringer(params[0])
	if isString {
		_, err := os.Stat(path)
		if os.IsNotExist(err) {
			return true, ""
		} else if err != nil {
			return false, fmt.Sprintf("other stat error: %v", err)
		}
		return false, fmt.Sprintf("%s exists", path)
	}

	value := reflect.ValueOf(params[0])
	return false, fmt.Sprintf("obtained value is not a string and has no .String(), %s:%#v", value.Kind(), params[0])
}

// SymlinkDoesNotExist checker makes sure the path specified doesn't exist.

type symlinkDoesNotExistChecker struct {
	*CheckerInfo
}

var SymlinkDoesNotExist Checker = &symlinkDoesNotExistChecker{
	&CheckerInfo{Name: "SymlinkDoesNotExist", Params: []string{"obtained"}},
}

func (checker *symlinkDoesNotExistChecker) Check(params []any, names []string) (result bool, error string) {
	path, isString := stringOrStringer(params[0])
	if isString {
		_, err := os.Lstat(path)
		if os.IsNotExist(err) {
			return true, ""
		} else if err != nil {
			return false, fmt.Sprintf("other stat error: %v", err)
		}
		return false, fmt.Sprintf("%s exists", path)
	}

	value := reflect.ValueOf(params[0])
	return false, fmt.Sprintf("obtained value is not a string and has no .String(), %s:%#v", value.Kind(), params[0])
}

// Same path checker -- will check that paths are the same OS indepentent

type samePathChecker struct {
	*CheckerInfo
}

// SamePath checks paths to see whether they're the same, can follow symlinks and is OS independent
var SamePath Checker = &samePathChecker{
	&CheckerInfo{Name: "SamePath", Params: []string{"obtained", "expected"}},
}

func (checker *samePathChecker) Check(params []any, names []string) (result bool, error string) {
	// Check for panics
	defer func() {
		if panicked := recover(); panicked != nil {
			result = false
			error = fmt.Sprint(panicked)
		}
	}()

	// Convert input
	obtained, isStr := stringOrStringer(params[0])
	if !isStr {
		return false, fmt.Sprintf("obtained value is not a string and has no .String(), %T:%#v", params[0], params[0])
	}
	expected, isStr := stringOrStringer(params[1])
	if !isStr {
		return false, fmt.Sprintf("obtained value is not a string and has no .String(), %T:%#v", params[1], params[1])
	}

	// Convert paths to proper format
	obtained = filepath.FromSlash(obtained)
	expected = filepath.FromSlash(expected)

	// If running on Windows, paths will be case-insensitive and thus we
	// normalize the inputs to a default of all upper-case
	if runtime.GOOS == "windows" {
		obtained = strings.ToUpper(obtained)
		expected = strings.ToUpper(expected)
	}

	// Same path do not check further
	if obtained == expected {
		return true, ""
	}

	// If it's not the same path, check if it points to the same file.
	// Thus, the cases with windows-shortened paths are accounted for
	// This will throw an error if it's not a file
	ob, err := os.Stat(obtained)
	if err != nil {
		return false, err.Error()
	}

	ex, err := os.Stat(expected)
	if err != nil {
		return false, err.Error()
	}

	res := os.SameFile(ob, ex)
	if res {
		return true, ""
	}
	return false, "Not the same file"
}
