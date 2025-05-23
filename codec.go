// Copyright 2012-2014 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package tc

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"
)

type codecEqualChecker struct {
	name      string
	marshal   func(any) ([]byte, error)
	unmarshal func([]byte, any) error
}

// JSONEquals defines a checker that checks whether a byte slice, when
// unmarshaled as JSON, is equal to the given value.
// Rather than unmarshaling into something of the expected
// body type, we reform the expected body in JSON and
// back to any, so we can check the whole content.
// Otherwise we lose information when unmarshaling.
var JSONEquals = &codecEqualChecker{
	name:      "JSONEquals",
	marshal:   json.Marshal,
	unmarshal: json.Unmarshal,
}

// YAMLEquals defines a checker that checks whether a byte slice, when
// unmarshaled as YAML, is equal to the given value.
// Rather than unmarshaling into something of the expected
// body type, we reform the expected body in YAML and
// back to any, so we can check the whole content.
// Otherwise we lose information when unmarshaling.
var YAMLEquals = &codecEqualChecker{
	name:      "YAMLEquals",
	marshal:   yaml.Marshal,
	unmarshal: yaml.Unmarshal,
}

func (checker *codecEqualChecker) Info() *CheckerInfo {
	return &CheckerInfo{
		Name:   checker.name,
		Params: []string{"obtained", "expected"},
	}
}

func (checker *codecEqualChecker) Check(params []any, names []string) (result bool, error string) {
	gotContent, ok := params[0].(string)
	if !ok {
		return false, fmt.Sprintf("expected string, got %T", params[0])
	}
	expectContent := params[1]
	expectContentBytes, err := checker.marshal(expectContent)
	if err != nil {
		return false, fmt.Sprintf("cannot marshal expected contents: %v", err)
	}
	var expectContentVal any
	if err := checker.unmarshal(expectContentBytes, &expectContentVal); err != nil {
		return false, fmt.Sprintf("cannot unmarshal expected contents: %v", err)
	}

	var gotContentVal any
	if err := checker.unmarshal([]byte(gotContent), &gotContentVal); err != nil {
		return false, fmt.Sprintf("cannot unmarshal obtained contents: %v; %q", err, gotContent)
	}

	if ok, err := DeepEqual(gotContentVal, expectContentVal); !ok {
		return false, err.Error()
	}
	return true, ""
}
