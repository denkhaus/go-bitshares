
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


install_secp256k1:
	@echo "######################## -> install secp256k1"
	@rm -rf /tmp/secp256k1
	git clone https://github.com/bitcoin-core/secp256k1.git /tmp/secp256k1
	cd /tmp/secp256k1 && ./autogen.sh 
	cd /tmp/secp256k1 && ./configure --enable-module-recovery
	cd /tmp/secp256k1 && make && sudo make install

init: install_secp256k1
	@echo "######################## -> install dependencies"
	@go get -u github.com/pquerna/ffjson
	@go get -u golang.org/x/tools/cmd/stringer
	@go get -u github.com/mitchellh/reflectwalk
	@go get -u github.com/stretchr/objx
	@go get -u github.com/cespare/reflex

test_api: 
	go test -v ./api -run ^TestCommon$
	go test -v ./api -run ^TestSubscribe$

test_operations:
	go test -v ./operations -run ^TestOperations$

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
