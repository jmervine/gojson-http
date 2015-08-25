package readable

import (
	"bytes"
	"log"
	"testing"
)

func Benchmark_Builtin_Logger(b *testing.B) {
	buf := new(bytes.Buffer)
	l := log.New(buf, "", log.LstdFlags)
	for n := 0; n < b.N; n++ {
		l.Print("foo bar")
	}
}

func Benchmark_KeyValue_Log(b *testing.B) {
	buf := new(bytes.Buffer)
	r := New().WithOutput(buf)
	for n := 0; n < b.N; n++ {
		r.Log("foo", "bar")
	}
}

func Benchmark_KeyValue_LogSafe(b *testing.B) {
	buf := new(bytes.Buffer)
	r := New().WithOutput(buf)
	for n := 0; n < b.N; n++ {
		r.LogSafe("foo", "bar")
	}
}

func Benchmark_Join_Log(b *testing.B) {
	buf := new(bytes.Buffer)
	r := New().WithOutput(buf).WithFormatter(Join)
	for n := 0; n < b.N; n++ {
		r.Log("foo", "bar")
	}
}

func Benchmark_Join_LogSafe(b *testing.B) {
	buf := new(bytes.Buffer)
	r := New().WithOutput(buf).WithFormatter(Join)
	for n := 0; n < b.N; n++ {
		r.LogSafe("foo", "bar")
	}
}

func Benchmark_KeyValue_WithSETTER_Log(b *testing.B) {
	buf := new(bytes.Buffer)
	r := New().WithOutput(buf)
	for n := 0; n < b.N; n++ {
		r.WithPrefix("prefix").Log("foo", "bar")
	}
}

func Benchmark_KeyValue_WithSETTER_LogSafe(b *testing.B) {
	buf := new(bytes.Buffer)
	r := New().WithOutput(buf)
	for n := 0; n < b.N; n++ {
		r.WithPrefix("prefix").LogSafe("foo", "bar")
	}
}

func Benchmark_KeyValue_WithSETTERs_Log(b *testing.B) {
	buf := new(bytes.Buffer)
	r := New()
	for n := 0; n < b.N; n++ {
		r.WithPrefix("prefix").WithOutput(buf).Log("foo", "bar")
	}
}

func Benchmark_KeyValue_WithSETTERs_LogSafe(b *testing.B) {
	buf := new(bytes.Buffer)
	r := New()
	for n := 0; n < b.N; n++ {
		r.WithPrefix("prefix").WithOutput(buf).LogSafe("foo", "bar")
	}
}
