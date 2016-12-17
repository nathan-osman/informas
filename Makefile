CWD = $(shell pwd)

all: dist/informas

dist/informas:
	docker run \
	    --rm \
	    -e CGO_ENABLED=0 \
	    -v ${CWD}:/go/src/github.com/nathan-osman/informas \
	    -v ${CWD}/dist:/go/bin \
	    -w /go/src/github.com/nathan-osman/informas \
	    golang:latest
	    go get ./...

clean:
	@rm dist/informas

.PHONY: clean
