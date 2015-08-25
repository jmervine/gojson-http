package readable

import (
	"testing"

	. "github.com/jmervine/GoT"
)

type TestStruct struct {
	A string
	B int
	C bool
}

func TestReadable_KeyValue(T *testing.T) {
	str := KeyValue("foo", "bar", 1)
	exp := "foo=bar 1"
	Go(T).AssertEqual(str, exp)

	str = KeyValue("foo", "bar", true)
	exp = "foo=bar true"
	Go(T).AssertEqual(str, exp)

	str = KeyValue("foo", "bar", "args", []string{"foo", "bar"})
	exp = "foo=bar args=[foo bar]"
	Go(T).AssertEqual(str, exp)

	str = KeyValue("foo", "bar", "struct", TestStruct{"foo", 9, false})
	exp = "foo=bar struct={A:foo B:9 C:false}"
	Go(T).AssertEqual(str, exp)

	// weird one
	obj := TestStruct{"foo", 9, false}
	str = KeyValue(obj, obj)
	exp = "{A:foo B:9 C:false}={A:foo B:9 C:false}"
}

func TestReadable_Join(T *testing.T) {
	str := Join("foo", "bar", 1)
	exp := "foo: bar 1"
	Go(T).AssertEqual(str, exp)

	str = Join("foo", "bar", true)
	exp = "foo: bar true"
	Go(T).AssertEqual(str, exp)

	str = Join("foo", "bar", "args", []string{"foo", "bar"})
	exp = "foo: bar args: [foo bar]"
	Go(T).AssertEqual(str, exp)

	str = Join("foo", "bar", "struct", TestStruct{"foo", 9, false})
	exp = "foo: bar struct: {A:foo B:9 C:false}"
	Go(T).AssertEqual(str, exp)

	// weird one
	obj := TestStruct{"foo", 9, false}
	str = KeyValue(obj, obj)
	exp = "{A:foo B:9 C:false}: {A:foo B:9 C:false}"
}
