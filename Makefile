
clean_ffjson_base: 
	@rm -rf types/ffjson-inception* ||:
	@rm -f types/*_ffjson_expose.go ||:
	@rm -rf operations/ffjson-inception* ||:
	@rm -f operations/*_ffjson_expose.go ||:

clean_ffjson_gen:
	@rm -f types/*_ffjson.go ||: 
	@rm -rf operations/*_ffjson.go ||: 

generate: clean_ffjson_base	
	-@go generate types/gen.go
	-@go generate ./...

generate_new: clean_ffjson_base clean_ffjson_gen		
	-@go generate types/gen.go
	-@go generate ./...

init: 
	@echo "######################## -> install/update dev dependencies"
	@go get -u golang.org/x/tools/cmd/stringer
	@go get -u github.com/cheekybits/genny
	@go get -u github.com/pquerna/ffjson
	@go get -u github.com/mitchellh/reflectwalk
	@go get -u github.com/stretchr/objx
	@go get -u github.com/stretchr/testify
	@go get -u github.com/cespare/reflex
	@go get -u github.com/bradhe/stopwatch

test_api: 
	@echo "######################## -> test bitshares api"
	-go test -v ./tests -run ^TestCommon$
	-go test -v ./tests -run ^TestSubscribe$
	-go test -v ./tests -run ^TestWalletAPI$
	-go test -v ./tests -run ^TestWebsocketAPI$
	-go test -v ./types 

test_operations:
	@echo "######################## -> test operations"
	@go test -v ./operations -run ^TestOperations$

test_blocks:
	@echo "this is a long running test, abort with Ctrl + C"
	@go test -v ./tests -timeout 10m -run ^TestBlockRange$

buildgen:
	@echo "######################## -> build btsgen"
	@cd ./gen && go get -u -d
	@cd ./gen && go build -o /tmp/btsgen
	@cp /tmp/btsgen $(GOPATH)/bin

opsamples: buildgen
	@echo "######################## -> exec btsgen"
	@cd ./gen && btsgen

build: generate
	go build -o /tmp/go-tmpbuild ./operations 

watch:
	reflex -g 'operations/*.go' make test_operations
