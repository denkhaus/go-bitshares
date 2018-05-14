
all: generate

clean_ffjson_base: 
	@rm -rf types/ffjson-inception* ||:	
	@rm types/*_ffjson_expose.go ||:	
	@rm -rf operations/ffjson-inception* ||:	
	@rm operations/*_ffjson_expose.go ||:	

clean_ffjson_gen:
	@rm types/*_ffjson.go ||:	
	@rm operations/*_ffjson.go ||:	

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

test:
	go test -v ./...

buildgen:
	@echo "build btsgen"
	@go get -u -d ./gen 
	@go build -o /tmp/btsgen ./gen 
	@cp /tmp/btsgen $(GOPATH)/bin

opsamples: buildgen
	@echo "exec btsgen"
	@cd gen && btsgen
