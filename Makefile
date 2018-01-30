all: test

jenkins: install_dep_tool install_go_linter production_restore_deps test

install_dep_tool:
	go get github.com/tools/godep

install_go_linter:
	go get -u -v github.com/golang/lint/golint

initialize_deps:
	go get -d -v ./...
	go get -d -v github.com/stretchr/testify/assert
	go get -v github.com/golang/lint/golint
	godep save

update_deps:
	godep go get -d -v ./...
	godep go get -d -v github.com/stretchr/testify/assert
	godep go get -v github.com/golang/lint/golint
	godep update ./...

test:
	golint ./...
	godep go vet ./...
	godep go test -v

production_restore_deps:
	godep restore

.PHONY: install_dep_tool install_go_linter initialize_deps update_deps test production_restore_deps

