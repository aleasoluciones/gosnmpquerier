all: update_deps build test

deps:
	go get -d -v ./...
	go get -d -v github.com/stretchr/testify/assert
	go get -d -v github.com/axw/gocov/gocov
	go get -d -v github.com/mattn/goveralls
	go get -v github.com/golang/lint/golint

update_deps:
	go get -d -v -u ./...
	go get -d -v -u github.com/stretchr/testify/assert
	go get -d -v -u github.com/axw/gocov/gocov
	go get -d -v -u github.com/mattn/goveralls
	go get -v -u github.com/golang/lint/golint

test:
	golint ./...
	go vet ./...
	go test -v

build:
	go build -v .

.PHONY: deps update_deps test build

