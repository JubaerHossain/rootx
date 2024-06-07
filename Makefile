BINARY_NAME := rootx

build:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o ./release/rootx-linux-amd64 ./main.go
	GOOS=darwin GOARCH=amd64 go build -o ./release/rootx-darwin-amd64 ./main.go
	GOOS=windows GOARCH=amd64 go build -o ./release/rootx-windows-amd64.exe ./main.go

clean:
	rm -f ./release/rootx-*

.PHONY: build clean