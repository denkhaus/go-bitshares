
clean_ffjson_base: 
	@rm -rf types/ffjson-inception* ||:
	@rm -f types/*_ffjson_expose.go ||:
	@rm -rf operations/ffjson-inception* ||:
	@rm -f operations/*_ffjson_expose.go ||:
	@rm -rf api/ffjson-inception* ||:
	@rm -f api/*_ffjson_expose.go ||:

clean_ffjson_gen:
	@rm -rf types/*_ffjson.go ||: 
	@rm -rf operations/*_ffjson.go ||:
	@rm -rf api/*_ffjson.go ||: 

generate: clean_ffjson_base
	@echo "######################## -> generate"
	-@go generate types/gen.go
	-@go generate ./...

generate_new: clean_ffjson_base clean_ffjson_gen
	@echo "######################## -> generate new"		
	-@go generate types/gen.go
	-@go generate ./...

init: 
	@echo "######################## -> install/update dev dependencies"
	@GO111MODULE=on go get -u golang.org/x/tools/cmd/stringer
	@GO111MODULE=on go get -u github.com/cheekybits/genny
	@GO111MODULE=on go get -u github.com/pquerna/ffjson
	@GO111MODULE=on go get -u github.com/mitchellh/reflectwalk
	@GO111MODULE=on go get -u github.com/stretchr/objx
	@GO111MODULE=on go get -u github.com/stretchr/testify
	@GO111MODULE=on go get -u github.com/cespare/reflex
	@GO111MODULE=on go get -u github.com/bradhe/stopwatch

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
	@echo "######################## -> build"
	go build -o /tmp/go-tmpbuild ./operations 

watch:
	reflex -g 'operations/*.go' make test_operations
