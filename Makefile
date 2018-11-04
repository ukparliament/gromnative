.PHONY: build

build:
	go build -buildmode=c-shared -o ./ext/native.so ./ext/native.go