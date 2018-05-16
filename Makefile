
all: generate

clean_ffjson_base: 
	@rm -rf types/ffjson-inception* ||:
	@rm -f types/*_ffjson_expose.go ||:
	@rm -rf operations/ffjson-inception* ||:
	@rm -f operations/*_ffjson_expose.go ||:

clean_ffjson_gen:
	@rm -f types/*_ffjson.go ||: 
	@rm -rf operations/*_ffjson.go ||: 

generate: clean_ffjson_base	
	-go generate ./...

generate_new: clean_ffjson_base clean_ffjson_gen		
	-go generate ./...

#install dependencies
init:
	@go get -u github.com/pquerna/ffjson
	@go get -u golang.org/x/tools/cmd/stringer
	@go get -u github.com/mitchellh/reflectwalk
	@go get -u github.com/stretchr/objx
	@go get -u github.com/cespare/reflex

test_all:
	go test -v ./...

test_operations:
	go test -v ./operations

buildgen:
	@echo "build btsgen"
	@go get -u -d ./gen 
	@go build -o /tmp/btsgen ./gen 
	@cp /tmp/btsgen $(GOPATH)/bin

opsamples: buildgen
	@echo "exec btsgen"
	@cd gen && btsgen

build: generate
	go build -o /tmp/go-tmpbuild ./operations 

watch:
	reflex -g 'operations/*.go' make test_operations
