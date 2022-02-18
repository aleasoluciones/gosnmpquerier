all: clean build test

update_dep:
	go get $(DEP)
	go mod tidy

update_all_deps:
	go get -u all
	go mod tidy

format:
	go fmt ./...

test:
	go vet ./...
	go clean -testcache
	go test -v ./...

build:
	go build examples/amqpsnmpqueries/amqpsnmpqueries.go
	go build examples/get/get.go
	go build examples/snmphttpserver/snmphttpserver.go
	go build examples/snmpqueries/snmpqueries.go
	go build examples/sync_snmpqueries/sync_snmpqueries.go
	go build examples/walk/walk.go

clean:
	rm -f amqpsnmpqueries
	rm -f get
	rm -f snmphttpserver
	rm -f snmpqueries
	rm -f sync_snmpqueries
	rm -f walk


.PHONY: update_dep update_all_deps format test build clean

