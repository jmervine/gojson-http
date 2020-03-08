full:
	go test . -bench=. -cover -v

travis:
	go test .

test:
	go test . -count 500

benchmark:
	go test . -bench=.

cover:
	go test . -cover -race

docker/test:
	go get -v github.com/jmervine/GoT
	go test .

test/versions: test/1.5 test/1.4 test/1.3

test/1.3:
	docker run --rm -it \
		-v $(shell pwd):/go/src/github.com/jmervine/readable \
		-w /go/src/github.com/jmervine/readable \
		golang:1.3 \
		make docker/test

test/1.4:
	docker run --rm -it \
		-v $(shell pwd):/go/src/github.com/jmervine/readable \
		-w /go/src/github.com/jmervine/readable \
		golang:1.4 \
		make docker/test

test/1.5:
	docker run --rm -it \
		-v $(shell pwd):/go/src/github.com/jmervine/readable \
		-w /go/src/github.com/jmervine/readable \
		golang:1.5 \
		make docker/test

