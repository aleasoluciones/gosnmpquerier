build:
	go build -v .

test:
	go test -v

deps:
	go get -d -v .
	go get github.com/stretchr/testify/assert
	go get github.com/axw/gocov/gocov
	go get github.com/mattn/goveralls

.PHONY: test build deps
