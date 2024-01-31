BINARY_NAME=webserv
VERSION := $(shell git describe --tags)
BUILD_TIMESTAMP := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ')
COMMIT_HASH := $(shell git rev-parse HEAD)

build:
	GOARCH=amd64 GOOS=darwin go build -ldflags "-X static/cli.Version=${VERSION} -X static/cli.BuildTimestamp=${BUILD_TIMESTAMP} -X static/cli.CommitHash=${COMMIT_HASH}" -o ${BINARY_NAME}_${VERSION}-darwin 
	GOARCH=amd64 GOOS=linux go build -ldflags "-X static/cli.Version=${VERSION} -X static/cli.BuildTimestamp=${BUILD_TIMESTAMP} -X static/cli.CommitHash=${COMMIT_HASH}" -o ${BINARY_NAME}_${VERSION}-linux 
	GOARCH=amd64 GOOS=windows go build -ldflags "-X static/cli.Version=${VERSION} -X static/cli.BuildTimestamp=${BUILD_TIMESTAMP} -X static/cli.CommitHash=${COMMIT_HASH}" -o ${BINARY_NAME}_${VERSION}-windows 

run: build

clean:
	go clean
	rm ${BINARY_NAME}_${VERSION}-darwin
	rm ${BINARY_NAME}_${VERSION}-linux
	rm ${BINARY_NAME}_${VERSION}-windows