
all: generate

clean: 
	-rm -rf objects/ffjson-inception*
	-rm objects/*_ffjson_expose.go
	-rm objects/*_ffjson.go

generate: clean	
	-go generate ./...