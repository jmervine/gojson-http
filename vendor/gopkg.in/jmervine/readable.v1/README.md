# readable

`readable` is a simple Golang logger based loosely on the 12 Factor Logs (http://12factor.net/logs) and backed by Goalng's built in `log` package.

[![GoDoc](https://godoc.org/gopkg.in/jmervine/readable.v1?status.png)](https://godoc.org/gopkg.in/jmervine/readable.v1) [![Build Status](https://travis-ci.org/jmervine/readable.svg?branch=master)](https://travis-ci.org/jmervine/readable)

## notes on thread safety

The [Go `log` package](http://godoc.org/log), which is the underlying implementation is suppose to be thread safe,
however, with my testing, when using buffers ([`bytes.Buffer`](http://godoc.org/bytes#Buffer)), it proved not to be in conjunction
with `With{*}` style setters. This shouldn't be an issue when using [`os.Stdout`](http://godoc.org/os#Stdout) and
[`os.Stderr`](http://godoc.org/os#Stderr), as they should be as safe as the standard `log` package.

To be safe, I've added `LogSafe` as a guaranteed thread safe logger, but it's a little slower than the standard logging.
Additionally, this should be used when using `With{*}` style setters in line, to ensure thread safety. Otherwise, there's
a good chance it won't be, as stated above. Example:

```go
// maybe not safe
go func() {
    readable.WithPrefix("foo").Log("bar")
    // do stuff
}()

// guaranteed safe
go func() {
    readable.WithPrefix("foo").LogSafe("bar")
    // do stuff
}()
```

## usage

```go
package main

import "gopkg.in/jmervine/readable.v1"

func main() {
    readable.Log("type", "default", "fn", "Log")
    // 2015/08/23 19:17:38 type=default fn=Log

    readable.SetFlags(0)
    readable.Log("type", "without data stamp", "fn", "Log")
    // type=without data stamp fn=Log

    readable.SetFormatter(readable.Join)
    readable.Log("type", "with Join formatter", "fn", "Log")
    // type: with Join formatter fn: Log

    logger := readable.New().WithPrefix("[INFO]:")
    debugger := readable.New().WithDebug().WithPrefix("[DEBUG]:")

    logger.Log("type", "logger", "fn", "Log")
    logger.Debug("type", "logger", "fn", "Debug")
    logger.WithDebug().Debug("type", "logger", "fn", "WithDebug().Debug")
    // 2015/08/23 19:17:38 [INFO]: type=logger fn=Log
    // 2015/08/23 19:17:38 [INFO]: type=logger fn=WithDebug().Debug

    debugger.Log("type", "debugger", "fn", "Log")
    debugger.Debug("type", "debuger", "fn", "Debug")
    // 2015/08/23 19:17:38 [DEBUG]: type=debugger fn=Log
    // 2015/08/23 19:17:38 [DEBUG]: type=debugger fn=Debug
}
```

## performance

`readable` performs a little slower then the default Go logger. `With{*}` setter
convenience methods, do impact performance, so keep that in mind where performance
is king.

```
go test . -bench=.
PASS
Benchmark_Builtin_Logger-4                       3000000               529 ns/op
Benchmark_KeyValue_Log-4                         1000000              1217 ns/op
Benchmark_KeyValue_LogSafe-4                     1000000              1312 ns/op
Benchmark_Join_Log-4                             1000000              1222 ns/op
Benchmark_Join_LogSafe-4                         1000000              1320 ns/op
Benchmark_KeyValue_WithSETTER_Log-4              1000000              1715 ns/op
Benchmark_KeyValue_WithSETTER_LogSafe-4          1000000              1794 ns/op
Benchmark_KeyValue_WithSETTERs_Log-4             1000000              1860 ns/op
Benchmark_KeyValue_WithSETTERs_LogSafe-4         1000000              1965 ns/op
ok      github.com/jmervine/readable    14.919s
```

## MIT Licence

```
Copyright (c) 2015 Joshua P. Mervine

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
of the Software, and to permit persons to whom the Software is furnished to do
so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
```
