
all: generate

clean: 
#	-rm -rf objects/ffjson-inception*
#	-rm objects/*_ffjson_expose.go

generate: clean	
	-go generate ./...