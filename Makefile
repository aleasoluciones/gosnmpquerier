build: deps
	go build -v .

test: deps
	go test -v

deps:
	go get -d -v .
	go get github.com/stretchr/testify/assert
	go get github.com/axw/gocov/gocov
	go get github.com/mattn/goveralls
	go get golang.org/x/tools/cmd/cover

.PHONY: test build deps
