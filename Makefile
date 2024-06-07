BINARY_NAME := rootx

build:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o ./release/$(BINARY_NAME)-linux-amd64 ./main.go
	GOOS=darwin GOARCH=amd64 go build -o ./release/$(BINARY_NAME)-darwin-amd64 ./main.go
	GOOS=windows GOARCH=amd64 go build -o ./release/$(BINARY_NAME)-windows-amd64.exe ./main.go

clean:
	rm -f ./release/$(BINARY_NAME)-*

.PHONY: build clean