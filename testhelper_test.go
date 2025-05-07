// Gocheck - A rich testing framework for Go
//
// Copyright (c) 2010-2013 Gustavo Niemeyer <gustavo@niemeyer.net>
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
//    list of conditions and the following disclaimer.
// 2. Redistributions in binary form must reproduce the above copyright notice,
//    this list of conditions and the following disclaimer in the documentation
//    and/or other materials provided with the distribution.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
// ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
// (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
// ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
// SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package tc_test

import (
	"flag"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"

	check "github.com/juju/tc"
)

var (
	helperRunFlag   = flag.String("helper.run", "", "Run helper suite")
	helperPanicFlag = flag.String("helper.panic", "", "")
)

func TestHelperSuite(t *testing.T) {
	if helperRunFlag == nil || *helperRunFlag == "" {
		t.SkipNow()
	}
	switch *helperRunFlag {
	case "FailHelper":
		check.Run(t, &FailHelper{})
	case "SuccessHelper":
		check.Run(t, &SuccessHelper{})
	case "FixtureHelper":
		suite := &FixtureHelper{}
		if helperPanicFlag != nil {
			suite.panicOn = *helperPanicFlag
		}
		check.Run(t, suite)
	case "integrationTestHelper":
		check.Run(t, &integrationTestHelper{})
	default:
		t.Skip()
	}
}

type helperResult []string

var (
	testRunLine    = regexp.MustCompile(`^=== (?:RUN|CONT)\s+([0-9A-Za-z/]+)$`)
	testStatusLine = regexp.MustCompile(`^\s*--- ([A-Z]+): ([0-9A-Za-z/]+) \(\d+\.\d+s\)$`)
)

func (result helperResult) Status(test string) string {
	for _, line := range result {
		match := testStatusLine.FindStringSubmatch(line)
		if match == nil {
			continue
		}
		if match[2] == "TestHelperSuite/"+test {
			return match[1]
		}
	}
	return ""
}

func (result helperResult) Logs(test string) string {
	var lines []string
	var inTest bool
	for _, line := range result {
		if inTest {
			// Log messages are all indented
			if strings.HasPrefix(line, " ") {
				lines = append(lines, line)
				continue
			}
			inTest = false
		}
		match := testRunLine.FindStringSubmatch(line)
		if match != nil && match[1] == "TestHelperSuite/"+test {
			inTest = true
		}
	}
	return strings.Join(lines, "\n")
}

func runHelperSuite(name string, args ...string) (code int, output helperResult) {
	args = append([]string{"-test.v", "-test.run", "TestHelperSuite", "-helper.run", name}, args...)
	cmd := exec.Command(os.Args[0], args...)
	data, err := cmd.Output()
	output = strings.Split(string(data), "\n")
	if execErr, ok := err.(*exec.ExitError); ok {
		code = execErr.ExitCode()
		err = nil
	}
	if err != nil {
		panic(err)
	}
	return code, output
}
