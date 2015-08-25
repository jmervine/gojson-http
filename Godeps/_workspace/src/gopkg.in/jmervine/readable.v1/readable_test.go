package readable

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"
	"testing"

	. "github.com/jmervine/GoT"
)

// HELPERS:

func reset() {
	std = New()
}

func stubLogger() (*bytes.Buffer, *Readable) {
	b := new(bytes.Buffer)
	r := new(Readable)
	r.logger = log.New(b, "", 0)
	r.formatter = KeyValue
	return b, r
}

var basicFormatter = func(parts ...interface{}) string {
	return fmt.Sprintln(parts...)
}

var noopFormatter = func(parts ...interface{}) string {
	return ""
}

// TESTS

func TestReadable_WithDebug(T *testing.T) {
	r := New()

	Go(T).AssertEqual(r.WithDebug().debug, true)
	Go(T).AssertEqual(r.debug, false)
}

func TestReadable_SetDebug(T *testing.T) {
	r := New()
	Go(T).AssertEqual(r.debug, false)

	r.SetDebug(true)
	Go(T).AssertEqual(r.debug, true)

	r.SetDebug(false)
	Go(T).AssertEqual(r.debug, false)
}

func TestSetDebug(T *testing.T) {
	reset()

	Go(T).AssertEqual(std.debug, false)

	SetDebug(true)
	Go(T).AssertEqual(std.debug, true)

	SetDebug(false)
	Go(T).AssertEqual(std.debug, false)
}

func TestReadable_WithFormatter(T *testing.T) {
	r := New()

	Go(T).AssertEqual(
		r.WithFormatter(basicFormatter).formatter,
		basicFormatter,
	)
	Go(T).AssertEqual(r.formatter, KeyValue)
}

func TestReadable_SetFormatter(T *testing.T) {
	r := New()
	Go(T).AssertEqual(r.formatter, KeyValue)
	r.SetFormatter(noopFormatter)
	Go(T).AssertEqual(r.formatter, noopFormatter)
}

func TestSetFormatter(T *testing.T) {
	reset()

	Go(T).AssertEqual(std.formatter, KeyValue)
	SetFormatter(noopFormatter)
	Go(T).AssertEqual(std.formatter, noopFormatter)
}

func TestReadable_WithPrefix(T *testing.T) {
	r := New()

	Go(T).AssertEqual(r.WithPrefix("prefix").prefix, "prefix")
	Go(T).AssertNil(r.prefix)
}

func TestReadable_SetPrefix(T *testing.T) {
	r := New()
	r.SetPrefix("prefix")
	Go(T).AssertEqual(r.prefix, "prefix")
}

func TestSetPrefix(T *testing.T) {
	reset()

	SetPrefix("prefix")

	Go(T).AssertEqual(std.prefix, "prefix")
}

func TestReadable_WithOutput(T *testing.T) {
	r := New()

	var b1 = new(bytes.Buffer)
	var b2 = new(bytes.Buffer)
	r.logger = log.New(b1, "", 0)

	r.WithOutput(b2).logger.Print("foo")

	Go(T).AssertEqual(b1.Len(), 0)
	Go(T).RefuteEqual(b2.Len(), 0)
}

func TestReadable_SetOutput(T *testing.T) {
	r := New()

	var b = new(bytes.Buffer)
	r.SetOutput(b)
	r.logger.Print("foo")

	Go(T).RefuteEqual(b.Len(), 0)
}

func TestSetOutput(T *testing.T) {
	reset()

	var b = new(bytes.Buffer)
	SetOutput(b)
	std.logger.Print("foo")

	Go(T).RefuteEqual(b.Len(), 0)
}

func TestReadable_WithFlags(T *testing.T) {
	b1 := new(bytes.Buffer)
	b2 := new(bytes.Buffer)

	r1 := New()
	r1.logger = log.New(b1, "", log.LstdFlags)

	r2 := r1.WithFlags(0)
	r2.logger = log.New(b2, "", r2.flags)

	r1.logger.Print("foo")
	r2.logger.Print("foo")

	Go(T).AssertEqual(b1.Len(), 24) // date len being 20
	Go(T).AssertEqual(b2.String(), "foo\n")
}

func TestReadable_SetFlags(T *testing.T) {
	b := new(bytes.Buffer)

	r := New()
	r.logger = log.New(b, "", log.LstdFlags)

	r.SetFlags(0)
	r.logger.Print("foo=bar")

	Go(T).AssertEqual(b.String(), "foo=bar\n")
}

func TestSetFlags(T *testing.T) {
	reset()

	b := new(bytes.Buffer)
	std.logger = log.New(b, "", log.LstdFlags)
	SetFlags(0)

	std.logger.Print("foo")
	Go(T).AssertEqual(b.String(), "foo\n")
}

func TestLog(T *testing.T) {
	reset()

	b := new(bytes.Buffer)
	SetOutput(b)

	Log("foo", "bar")
	Go(T).AssertEqual(b.Len(), 28)
}

func TestReadable_Log(T *testing.T) {
	buf, logger := stubLogger()

	logger.Log("foo", "bar")
	Go(T).AssertEqual(buf.String(), "foo=bar\n")
}

func TestReadable_LogSetPrefix(T *testing.T) {
	buf, logger := stubLogger()

	logger.SetPrefix("logger")
	logger.Log("foo", "bar")
	Go(T).AssertEqual(buf.String(), "logger foo=bar\n")
}

// PRIVATES

func TestReadable_clone(T *testing.T) {
	r1 := New()
	r2 := r1.clone()

	r2.logger = log.New(os.Stdout, "", 1)
	r2.prefix = "foo"

	Go(T).RefuteEqual(r1.logger, r2.logger)
	Go(T).RefuteEqual(r1.prefix, r2.prefix)
}

func TestReadable_prepLine(T *testing.T) {
	logger := New()
	str := logger.prepLine("foo", "bar", 1)
	exp := "foo=bar 1"
	Go(T).AssertEqual(str, exp)
}

// MISC

func TestReadable_ThreadSaftyOne(T *testing.T) {
	reset()

	var b = new(bytes.Buffer)

	var wait sync.WaitGroup

	// setup std
	SetOutput(b)
	SetFlags(0)

	var exp int

	// save previous GOMAXPROCS
	procs := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(4)
	defer runtime.GOMAXPROCS(procs)

	for i := 0; i < 10000; i++ {
		wait.Add(1)
		exp = exp + 8
		go func() {
			defer wait.Done()
			Log("foo", "bar")
		}()
	}

	wait.Wait()
	Go(T).AssertEqual(b.Len(), exp)
}

func TestReadable_LogSafe_ThreadSafety(T *testing.T) {
	reset()

	var b = new(bytes.Buffer)

	var wait sync.WaitGroup
	var r = New()

	// setup std
	r.SetOutput(b)
	r.SetFlags(0)

	var exp int

	// save previous GOMAXPROCS
	procs := runtime.GOMAXPROCS(0)
	runtime.GOMAXPROCS(4)
	defer runtime.GOMAXPROCS(procs)

	for i := 0; i < 10000; i++ {
		wait.Add(1)
		exp = exp + 15
		go func() {
			defer wait.Done()
			// With{*} is known to not be thread safe when using Log
			r.WithPrefix("prefix").LogSafe("foo", "bar")
		}()
	}

	wait.Wait()
	Go(T).AssertEqual(b.Len(), exp)
}

// EXAMPLES
func ExampleReadable_New() {
	logger := New().SetPrefix("server").SetOutput(os.Stdout)
	logger.Log("listener", ":3000")
	// 2015/08/21 20:01:48 server listener=:3000
}

func ExampleReadable() {
	SetPrefix("server")
	Log("listener", ":3000")
	// 2015/08/21 20:01:48 server listener=:3000
}

func ExampleReadable_Log() {
	// setting some default to ensure that examples work correctly
	logger := New().SetFlags(0).SetOutput(os.Stdout)

	logger.Log("foo", "bar")
	// Output:
	// foo=bar
}

func ExampleReadable_SetFormatter() {
	logger := New().SetFlags(0).SetOutput(os.Stdout)

	// create custom formatter
	logger.SetFormatter(func(parts ...interface{}) string {
		return fmt.Sprintln(parts...)
	})

	logger.Log("foo", "bar")
	// Output:
	// foo bar
}

func ExampleReadable_SetPrefix() {
	logger := New().SetFlags(0).SetOutput(os.Stdout)

	// set prefix
	logger.SetPrefix("prefix")

	logger.Log("foo", "bar")
	// Output:
	// prefix foo=bar
}
