package main

import (
	// "gopkg.in/jmervine/readable.v1"
	"../.."
)

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
