all: format lint test

fieldAlignment:
	fieldalignment -fix github.com/twothicc/common-go/grpcclient
	fieldalignment -fix github.com/twothicc/common-go/grpcclient/pool

format:
	gofmt -s -w $$(find . -type f -name '*.go'| grep -v "/vendor/")

lint:
	golangci-lint run

test:
	go test -v