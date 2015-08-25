package main

import (
	// "gopkg.in/jmervine/readable.v1"
	"../.."
	"fmt"
	"net/http"
)

var logger = readable.New().WithPrefix("server")

func handler(w http.ResponseWriter, r *http.Request) {
	defer logger.Log("fn", "handler", "path", r.URL.Path, "method", r.Method)
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path)
}

func main() {
	http.HandleFunc("/", handler)
	logger.Log("fn", "main", "listener", ":8080")
	logger.Fatal("fn", "main", "error", http.ListenAndServe(":8080", nil))
}
