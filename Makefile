
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

test:
	go test -v ./...

opsamples:
	cd gen && go run *.go