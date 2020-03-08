package readable

import (
	"fmt"
	"strings"
)

// KeyValue formats like so:
//
//  readable.Log("package name", "server", "listern", ":3000", "extra stuff")
//  //=> 2015/08/21 20:01:48 "package name"=server listener=:3000 "extra stuff"
//
func KeyValue(parts ...interface{}) (line string) {

	segments := []string{}

	quote := func(i interface{}) string {
		s := fmt.Sprintf("%+v", i)
		if strings.Contains(s, " ") {
			return fmt.Sprintf("%q", s)
		}
		return s
	}

	if len(parts) == 1 {
		return fmt.Sprintf("%+v", parts[0])
	}

	for i := 0; i < len(parts); i++ {
		if i%2 == 0 {
			if i == len(parts)-1 {
				segments = append(segments, quote(parts[i]))
			} else {
				segments = append(segments, fmt.Sprintf("%s=%s", parts[i], quote(parts[i+1])))
			}
		}

	}

	return strings.Join(segments, " ")
}

// Join is a simple formatter which formats like so:
//
//  readable.SetFormatter(readable.Join)
//  readable.Log("package", "server", "listern", ":3000")
//  //=> "2015/08/21 20:01:48 package server listener :3000"
func Join(parts ...interface{}) string {
	segments := []string{}

	for i := 0; i < len(parts); i++ {
		if i%2 == 0 {
			if i == len(parts)-1 {
				segments = append(segments, fmt.Sprintf("%+v", parts[i]))
			} else {
				segments = append(segments, fmt.Sprintf("%+v: %+v", parts[i], parts[i+1]))
			}
		}

	}

	return strings.Join(segments, " ")
}
