BINARY_NAME=webserv
VERSION=$(shell git describe --tags)
BUILD_TIMESTAMP=$(shell date +%FT%T%z)
COMMIT_HASH=$(shell git rev-parse HEAD)


build:
 GOARCH=amd64 GOOS=darwin go build -ldflags "-X main.Version=${VERSION} -X main.BuildTimestamp=${BUILD_TIMESTAMP} -X main.CommitHash=${COMMIT_HASH}" -o ${BINARY_NAME}_${VERSION}-darwin main.go
 GOARCH=amd64 GOOS=linux go build -ldflags "-X main.Version=${VERSION} -X main.BuildTimestamp=${BUILD_TIMESTAMP} -X main.CommitHash=${COMMIT_HASH}" -o ${BINARY_NAME}_${VERSION}-linux main.go
 GOARCH=amd64 GOOS=windows go build -ldflags "-X main.Version=${VERSION} -X main.BuildTimestamp=${BUILD_TIMESTAMP} -X main.CommitHash=${COMMIT_HASH}" -o ${BINARY_NAME}_${VERSION}-windows main.go

run: build
 ./${BINARY_NAME}

clean:
 go clean
 rm ${BINARY_NAME}-darwin
 rm ${BINARY_NAME}-linux
 rm ${BINARY_NAME}-window