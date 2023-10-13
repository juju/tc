// Copyright 2013 Canonical Ltd.
// Licensed under the LGPLv3, see LICENCE file for details.

package check_test

import (
	"github.com/juju/loggo"

	. "gopkg.in/check.v2"
)

type SimpleMessageSuite struct{}

var _ = Suite(&SimpleMessageSuite{})

func (s *SimpleMessageSuite) TestSimpleMessageString(c *C) {
	m := SimpleMessage{
		Level: loggo.INFO,
		Message: `hello
world
`,
	}
	c.Check(m.String(), Matches, "INFO hello\nworld\n")
}

func (s *SimpleMessageSuite) TestSimpleMessagesGoString(c *C) {
	m := SimpleMessages{{
		Level:   loggo.DEBUG,
		Message: "debug",
	}, {
		Level:   loggo.ERROR,
		Message: "Error",
	}}
	c.Check(m.GoString(), Matches, `SimpleMessages{
DEBUG debug
ERROR Error
}`)
}

type LogMatchesSuite struct{}

var _ = Suite(&LogMatchesSuite{})

func (s *LogMatchesSuite) TestMatchSimpleMessage(c *C) {
	log := []loggo.Entry{
		{Level: loggo.INFO, Message: "foo bar"},
		{Level: loggo.INFO, Message: "12345"},
	}
	c.Check(log, LogMatches, []SimpleMessage{
		{loggo.INFO, "foo bar"},
		{loggo.INFO, "12345"},
	})
	c.Check(log, LogMatches, []SimpleMessage{
		{loggo.INFO, "foo .*"},
		{loggo.INFO, "12345"},
	})
	// UNSPECIFIED means we don't care what the level is,
	// just check the message string matches.
	c.Check(log, LogMatches, []SimpleMessage{
		{loggo.UNSPECIFIED, "foo .*"},
		{loggo.INFO, "12345"},
	})
	c.Check(log, Not(LogMatches), []SimpleMessage{
		{loggo.INFO, "foo bar"},
		{loggo.DEBUG, "12345"},
	})
}

func (s *LogMatchesSuite) TestMatchSimpleMessages(c *C) {
	log := []loggo.Entry{
		{Level: loggo.INFO, Message: "foo bar"},
		{Level: loggo.INFO, Message: "12345"},
	}
	c.Check(log, LogMatches, SimpleMessages{
		{loggo.INFO, "foo bar"},
		{loggo.INFO, "12345"},
	})
	c.Check(log, LogMatches, SimpleMessages{
		{loggo.INFO, "foo .*"},
		{loggo.INFO, "12345"},
	})
	// UNSPECIFIED means we don't care what the level is,
	// just check the message string matches.
	c.Check(log, LogMatches, SimpleMessages{
		{loggo.UNSPECIFIED, "foo .*"},
		{loggo.INFO, "12345"},
	})
	c.Check(log, Not(LogMatches), SimpleMessages{
		{loggo.INFO, "foo bar"},
		{loggo.DEBUG, "12345"},
	})
}

func (s *LogMatchesSuite) TestMatchStrings(c *C) {
	log := []loggo.Entry{
		{Level: loggo.INFO, Message: "foo bar"},
		{Level: loggo.INFO, Message: "12345"},
	}
	c.Check(log, LogMatches, []string{"foo bar", "12345"})
	c.Check(log, LogMatches, []string{"foo .*", "12345"})
	c.Check(log, Not(LogMatches), []string{"baz", "bing"})
}

func (s *LogMatchesSuite) TestMatchInexact(c *C) {
	log := []loggo.Entry{
		{Level: loggo.INFO, Message: "foo bar"},
		{Level: loggo.INFO, Message: "baz"},
		{Level: loggo.DEBUG, Message: "12345"},
		{Level: loggo.ERROR, Message: "12345"},
		{Level: loggo.INFO, Message: "67890"},
	}
	c.Check(log, LogMatches, []string{"foo bar", "12345"})
	c.Check(log, LogMatches, []string{"foo .*", "12345"})
	c.Check(log, LogMatches, []string{"foo .*", "67890"})
	c.Check(log, LogMatches, []string{"67890"})

	// Matches are always left-most after the previous match.
	c.Check(log, LogMatches, []string{".*", "baz"})
	c.Check(log, LogMatches, []string{"foo bar", ".*", "12345"})
	c.Check(log, LogMatches, []string{"foo bar", ".*", "67890"})

	// Order is important: 67890 advances to the last item in obtained,
	// and so there's nothing after to match against ".*".
	c.Check(log, Not(LogMatches), []string{"67890", ".*"})
	// ALL specified patterns MUST match in the order given.
	c.Check(log, Not(LogMatches), []string{".*", "foo bar"})

	// Check that levels are matched.
	c.Check(log, LogMatches, []SimpleMessage{
		{loggo.UNSPECIFIED, "12345"},
		{loggo.UNSPECIFIED, "12345"},
	})
	c.Check(log, LogMatches, []SimpleMessage{
		{loggo.DEBUG, "12345"},
		{loggo.ERROR, "12345"},
	})
	c.Check(log, LogMatches, []SimpleMessage{
		{loggo.DEBUG, "12345"},
		{loggo.INFO, ".*"},
	})
	c.Check(log, Not(LogMatches), []SimpleMessage{
		{loggo.DEBUG, "12345"},
		{loggo.INFO, ".*"},
		{loggo.UNSPECIFIED, ".*"},
	})
}

func (s *LogMatchesSuite) TestFromLogMatches(c *C) {
	tw := &loggo.TestWriter{}
	_, err := loggo.ReplaceDefaultWriter(tw)
	c.Assert(err, IsNil)
	defer loggo.ResetWriters()
	logger := loggo.GetLogger("test")
	logger.SetLogLevel(loggo.DEBUG)
	logger.Infof("foo")
	logger.Debugf("bar")
	logger.Tracef("hidden")
	c.Check(tw.Log(), LogMatches, []string{"foo", "bar"})
	c.Check(tw.Log(), Not(LogMatches), []string{"foo", "bad"})
	c.Check(tw.Log(), Not(LogMatches), []SimpleMessage{
		{loggo.INFO, "foo"},
		{loggo.INFO, "bar"},
	})
}

func (s *LogMatchesSuite) TestLogMatchesOnlyAcceptsSliceTestLogValues(c *C) {
	obtained := []string{"banana"} // specifically not []loggo.TestLogValues
	expected := SimpleMessages{}
	result, err := LogMatches.Check([]interface{}{obtained, expected}, nil)
	c.Assert(result, Equals, false)
	c.Assert(err, Equals, "Obtained value must be of type []loggo.Entry or SimpleMessage")
}

func (s *LogMatchesSuite) TestLogMatchesOnlyAcceptsStringOrSimpleMessages(c *C) {
	obtained := []loggo.Entry{
		{Level: loggo.INFO, Message: "foo bar"},
		{Level: loggo.INFO, Message: "baz"},
		{Level: loggo.DEBUG, Message: "12345"},
	}
	expected := "totally wrong"
	result, err := LogMatches.Check([]interface{}{obtained, expected}, nil)
	c.Assert(result, Equals, false)
	c.Assert(err, Equals, "Expected value must be of type []string or []SimpleMessage")
}

func (s *LogMatchesSuite) TestLogMatchesFailsOnInvalidRegex(c *C) {
	var obtained interface{} = []loggo.Entry{{Level: loggo.INFO, Message: "foo bar"}}
	var expected interface{} = []string{"[]foo"}

	result, err := LogMatches.Check([]interface{}{obtained, expected}, nil /* unused */)
	c.Assert(result, Equals, false)
	c.Assert(err, Equals, "bad message regexp \"[]foo\": error parsing regexp: missing closing ]: `[]foo`")
}
