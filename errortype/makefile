all: format lint test

fieldAlignment:
	fieldalignment -fix github.com/twothicc/common-go/errortype

format:
	gofmt -s -w $$(find . -type f -name '*.go'| grep -v "/vendor/")

lint:
	golangci-lint run

test:
	go test -v