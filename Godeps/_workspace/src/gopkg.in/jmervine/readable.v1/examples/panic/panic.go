package main

import (
	"../.."
	"fmt"
)

func main() {
	readable.SetPrefix("example")
	readable.Log("fn", "main", "example", "panic")
	readable.Panic("fn", "main", "error", fmt.Errorf("example of panic"))
}
