// Copyright 2025 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import (
	"fmt"
	"regexp"
)

// regex for validating that the UUID matches RFC 4122.
// This package generates version 4 UUIDs but
// accepts any UUID version.
// http://www.ietf.org/rfc/rfc4122.txt
var (
	block1 = "[0-9a-f]{8}"
	block2 = "[0-9a-f]{4}"
	block3 = "[0-9a-f]{4}"
	block4 = "[0-9a-f]{4}"
	block5 = "[0-9a-f]{12}"

	uuidSnippet = block1 + "-" + block2 + "-" + block3 + "-" + block4 + "-" + block5
	validUUID   = regexp.MustCompile("^" + uuidSnippet + "$")
)

// IsZeroUUID checks the obtained value is zero a uuid.
var IsZeroUUID = &uuidChecker{
	example: "00000000-0000-0000-0000-000000000000",
	name:    "IsZeroUUID",
}

// IsUUID checks the obtained value is a uuid.
var IsUUID = &uuidChecker{
	name: "IsUUID",
}

// IsNonZeroUUID checks the obtained value is a non-zero uuid.
var IsNonZeroUUID = And(Not(IsZeroUUID), IsUUID)

type uuidChecker struct {
	example string
	name    string
}

func (checker *uuidChecker) Info() *CheckerInfo {
	info := CheckerInfo{
		Name:   checker.name,
		Params: []string{"obtained"},
	}
	return &info
}

func (checker *uuidChecker) Check(params []any, names []string) (result bool, error string) {
	uuid := ""
	switch v := params[0].(type) {
	case string:
		uuid = v
	case fmt.Stringer:
		uuid = v.String()
	default:
		return false, "obtained value type must be a string or fmt.Stringer"
	}

	if !validUUID.MatchString(uuid) {
		return false, "obtained value does not look like a uuid"
	}

	example := checker.example
	if example != "" && !validUUID.MatchString(example) {
		return false, "example value does not look like a uuid"
	}

	if example != "" && uuid != example {
		return false, ""
	}

	return true, ""
}
