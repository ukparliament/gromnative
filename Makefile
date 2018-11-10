.PHONY: build-and-test build install install-go install-ruby test test-go test-ruby check-built proto setup

build-and-test: install setup test

build:
	@echo "-- Building library"
	go build -buildmode=c-shared -o ./ext/gromnative.so ./ext/gromnative.go

install: install-go install-ruby

install-go:
	@echo "-- Installing Go Dependencies"
	dep ensure

install-ruby:
	@echo "-- Install Ruby Dependencies"
	bundle install

test: test-go test-ruby

test-go:
	@echo "-- Testing library"
	ginkgo ./ext/...

test-ruby: check-built
	@echo "-- Testing gem"
	bundle exec rake

check-built:
	@echo "-- Checking for built library"
	@[ -f ./ext/gromnative.so ] && echo "Lib found, not building" || make build

proto:
	@echo "-- Building Protobuf files"
	protowrap -I. --go_out `go env GOPATH`/src ./**/**/*.proto

setup:
	@echo "-- Installing testing framework"
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega/...