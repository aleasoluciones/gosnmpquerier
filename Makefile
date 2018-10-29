all: test

jenkins: install_dep_tool install_go_linter production_restore_deps test

install_dep_tool:
	go get github.com/tools/godep

install_go_linter:
	go get -u -v golang.org/x/lint/golint

initialize_deps:
	go get -d -v ./...
	go get -d -v github.com/stretchr/testify/assert
	go get -v golang.org/x/lint/golint
	godep save

update_deps:
	godep go install -v ./...
	godep go install -v github.com/stretchr/testify/assert
	godep go install -v golang.org/x/lint/golint
	godep update ./...

test:
	golint ./...
	godep go vet ./...
	godep go test -v

production_restore_deps:
	godep restore

.PHONY: install_dep_tool install_go_linter initialize_deps update_deps test production_restore_deps

