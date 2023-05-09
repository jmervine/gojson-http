PORT ?= 3333
LISTEN ?= 0.0.0.0
CWD=$(shell pwd)

bin/gojson-http:
	go mod download
	go mod tidy
	go mod verify
	go build -o bin/gojson-http main.go

.PHONY: clean
clean:
	rm -rvf bin/gojson-http

.PHONY: start
start: bin/gojson-http
	$(CWD)/bin/gojson-http -port $(PORT) -listen $(LISTEN) -template $(CWD)/index.html
