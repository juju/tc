// Copyright 2013 Canonical Ltd.
// Copyright 2014 Cloudbase Solutions SRL
// Licensed under the LGPLv3, see LICENCE file for details.

package check_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	. "gopkg.in/check.v2"
)

type FileSuite struct{}

var _ = Suite(&FileSuite{})

func (s *FileSuite) TestIsNonEmptyFile(c *C) {
	file, err := ioutil.TempFile(c.MkDir(), "")
	c.Assert(err, IsNil)
	fmt.Fprintf(file, "something")
	file.Close()

	c.Assert(file.Name(), IsNonEmptyFile)
}

func (s *FileSuite) TestIsNonEmptyFileWithEmptyFile(c *C) {
	file, err := ioutil.TempFile(c.MkDir(), "")
	c.Assert(err, IsNil)
	file.Close()

	result, message := IsNonEmptyFile.Check([]interface{}{file.Name()}, nil)
	c.Assert(result, IsFalse)
	c.Assert(message, Equals, file.Name()+" is empty")
}

func (s *FileSuite) TestIsNonEmptyFileWithMissingFile(c *C) {
	name := filepath.Join(c.MkDir(), "missing")

	result, message := IsNonEmptyFile.Check([]interface{}{name}, nil)
	c.Assert(result, IsFalse)
	c.Assert(message, Equals, name+" does not exist")
}

func (s *FileSuite) TestIsNonEmptyFileWithNumber(c *C) {
	result, message := IsNonEmptyFile.Check([]interface{}{42}, nil)
	c.Assert(result, IsFalse)
	c.Assert(message, Equals, "obtained value is not a string and has no .String(), int:42")
}

func (s *FileSuite) TestIsDirectory(c *C) {
	dir := c.MkDir()
	c.Assert(dir, IsDirectory)
}

func (s *FileSuite) TestIsDirectoryMissing(c *C) {
	absentDir := filepath.Join(c.MkDir(), "foo")

	result, message := IsDirectory.Check([]interface{}{absentDir}, nil)
	c.Assert(result, IsFalse)
	c.Assert(message, Equals, absentDir+" does not exist")
}

func (s *FileSuite) TestIsDirectoryWithFile(c *C) {
	file, err := ioutil.TempFile(c.MkDir(), "")
	c.Assert(err, IsNil)
	file.Close()

	result, message := IsDirectory.Check([]interface{}{file.Name()}, nil)
	c.Assert(result, IsFalse)
	c.Assert(message, Equals, file.Name()+" is not a directory")
}

func (s *FileSuite) TestIsDirectoryWithNumber(c *C) {
	result, message := IsDirectory.Check([]interface{}{42}, nil)
	c.Assert(result, IsFalse)
	c.Assert(message, Equals, "obtained value is not a string and has no .String(), int:42")
}

func (s *FileSuite) TestDoesNotExist(c *C) {
	absentDir := filepath.Join(c.MkDir(), "foo")
	c.Assert(absentDir, DoesNotExist)
}

func (s *FileSuite) TestDoesNotExistWithPath(c *C) {
	dir := c.MkDir()
	result, message := DoesNotExist.Check([]interface{}{dir}, nil)
	c.Assert(result, IsFalse)
	c.Assert(message, Equals, dir+" exists")
}

func (s *FileSuite) TestDoesNotExistWithSymlink(c *C) {
	dir := c.MkDir()
	deadPath := filepath.Join(dir, "dead")
	symlinkPath := filepath.Join(dir, "a-symlink")
	err := os.Symlink(deadPath, symlinkPath)
	c.Assert(err, IsNil)
	// A valid symlink pointing to something that doesn't exist passes.
	// Use SymlinkDoesNotExist to check for the non-existence of the link itself.
	c.Assert(symlinkPath, DoesNotExist)
}

func (s *FileSuite) TestDoesNotExistWithNumber(c *C) {
	result, message := DoesNotExist.Check([]interface{}{42}, nil)
	c.Assert(result, IsFalse)
	c.Assert(message, Equals, "obtained value is not a string and has no .String(), int:42")
}

func (s *FileSuite) TestSymlinkDoesNotExist(c *C) {
	absentDir := filepath.Join(c.MkDir(), "foo")
	c.Assert(absentDir, SymlinkDoesNotExist)
}

func (s *FileSuite) TestSymlinkDoesNotExistWithPath(c *C) {
	dir := c.MkDir()
	result, message := SymlinkDoesNotExist.Check([]interface{}{dir}, nil)
	c.Assert(result, IsFalse)
	c.Assert(message, Equals, dir+" exists")
}

func (s *FileSuite) TestSymlinkDoesNotExistWithSymlink(c *C) {
	dir := c.MkDir()
	deadPath := filepath.Join(dir, "dead")
	symlinkPath := filepath.Join(dir, "a-symlink")
	err := os.Symlink(deadPath, symlinkPath)
	c.Assert(err, IsNil)

	result, message := SymlinkDoesNotExist.Check([]interface{}{symlinkPath}, nil)
	c.Assert(result, IsFalse)
	c.Assert(message, Equals, symlinkPath+" exists")
}

func (s *FileSuite) TestSymlinkDoesNotExistWithNumber(c *C) {
	result, message := SymlinkDoesNotExist.Check([]interface{}{42}, nil)
	c.Assert(result, IsFalse)
	c.Assert(message, Equals, "obtained value is not a string and has no .String(), int:42")
}

func (s *FileSuite) TestIsSymlink(c *C) {
	file, err := ioutil.TempFile(c.MkDir(), "")
	c.Assert(err, IsNil)
	c.Log(file.Name())
	c.Log(filepath.Dir(file.Name()))
	symlinkPath := filepath.Join(filepath.Dir(file.Name()), "a-symlink")
	err = os.Symlink(file.Name(), symlinkPath)
	c.Assert(err, IsNil)

	c.Assert(symlinkPath, IsSymlink)
}

func (s *FileSuite) TestIsSymlinkWithFile(c *C) {
	file, err := ioutil.TempFile(c.MkDir(), "")
	c.Assert(err, IsNil)
	result, message := IsSymlink.Check([]interface{}{file.Name()}, nil)
	c.Assert(result, IsFalse)
	c.Assert(message, Contains, " is not a symlink")
}

func (s *FileSuite) TestIsSymlinkWithDir(c *C) {
	result, message := IsSymlink.Check([]interface{}{c.MkDir()}, nil)
	c.Assert(result, IsFalse)
	c.Assert(message, Contains, " is not a symlink")
}

func (s *FileSuite) TestSamePathWithNumber(c *C) {
	result, message := SamePath.Check([]interface{}{42, 52}, nil)
	c.Assert(result, IsFalse)
	c.Assert(message, Equals, "obtained value is not a string and has no .String(), int:42")
}

func (s *FileSuite) TestSamePathBasic(c *C) {
	dir := c.MkDir()

	result, message := SamePath.Check([]interface{}{dir, dir}, nil)

	c.Assert(result, IsTrue)
	c.Assert(message, Equals, "")
}

type SamePathLinuxSuite struct{}

var _ = Suite(&SamePathLinuxSuite{})

func (s *SamePathLinuxSuite) SetUpSuite(c *C) {
	if runtime.GOOS == "windows" {
		c.Skip("Skipped Linux-intented SamePath tests on Windows.")
	}
}

func (s *SamePathLinuxSuite) TestNotSamePathLinuxBasic(c *C) {
	dir := c.MkDir()
	path1 := filepath.Join(dir, "Test")
	path2 := filepath.Join(dir, "test")

	result, message := SamePath.Check([]interface{}{path1, path2}, nil)

	c.Assert(result, IsFalse)
	c.Assert(message, Equals, "stat "+path1+": no such file or directory")
}

func (s *SamePathLinuxSuite) TestSamePathLinuxSymlinks(c *C) {
	file, err := ioutil.TempFile(c.MkDir(), "")
	c.Assert(err, IsNil)
	symlinkPath := filepath.Join(filepath.Dir(file.Name()), "a-symlink")
	err = os.Symlink(file.Name(), symlinkPath)

	result, message := SamePath.Check([]interface{}{file.Name(), symlinkPath}, nil)

	c.Assert(result, IsTrue)
	c.Assert(message, Equals, "")
}

type SamePathWindowsSuite struct{}

var _ = Suite(&SamePathWindowsSuite{})

func (s *SamePathWindowsSuite) SetUpSuite(c *C) {
	if runtime.GOOS != "windows" {
		c.Skip("Skipped Windows-intented SamePath tests.")
	}
}

func (s *SamePathWindowsSuite) TestNotSamePathBasic(c *C) {
	dir := c.MkDir()
	path1 := filepath.Join(dir, "notTest")
	path2 := filepath.Join(dir, "test")

	result, message := SamePath.Check([]interface{}{path1, path2}, nil)

	c.Assert(result, IsFalse)
	path1 = strings.ToUpper(path1)
	c.Assert(message, Equals, "GetFileAttributesEx "+path1+": The system cannot find the file specified.")
}

func (s *SamePathWindowsSuite) TestSamePathWindowsCaseInsensitive(c *C) {
	dir := c.MkDir()
	path1 := filepath.Join(dir, "Test")
	path2 := filepath.Join(dir, "test")

	result, message := SamePath.Check([]interface{}{path1, path2}, nil)

	c.Assert(result, IsTrue)
	c.Assert(message, Equals, "")
}

func (s *SamePathWindowsSuite) TestSamePathWindowsFixSlashes(c *C) {
	result, message := SamePath.Check([]interface{}{"C:/Users", "C:\\Users"}, nil)

	c.Assert(result, IsTrue)
	c.Assert(message, Equals, "")
}

func (s *SamePathWindowsSuite) TestSamePathShortenedPaths(c *C) {
	dir := c.MkDir()
	dir1, err := ioutil.TempDir(dir, "Programming")
	defer os.Remove(dir1)
	c.Assert(err, IsNil)
	result, message := SamePath.Check([]interface{}{dir + "\\PROGRA~1", dir1}, nil)

	c.Assert(result, IsTrue)
	c.Assert(message, Equals, "")
}

func (s *SamePathWindowsSuite) TestSamePathShortenedPathsConsistent(c *C) {
	dir := c.MkDir()
	dir1, err := ioutil.TempDir(dir, "Programming")
	defer os.Remove(dir1)
	c.Assert(err, IsNil)
	dir2, err := ioutil.TempDir(dir, "Program Files")
	defer os.Remove(dir2)
	c.Assert(err, IsNil)

	result, message := SamePath.Check([]interface{}{dir + "\\PROGRA~1", dir2}, nil)

	c.Assert(result, Not(IsTrue))
	c.Assert(message, Equals, "Not the same file")

	result, message = SamePath.Check([]interface{}{"C:/PROGRA~2", "C:/Program Files (x86)"}, nil)

	c.Assert(result, IsTrue)
	c.Assert(message, Equals, "")
}
