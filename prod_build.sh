#!/bin/bash

VERSION="0.0.1"

# Creates production ready build
go build ./main/main.go -o webserv

cd ../


DIR_NAME="webserv_${VERSION}"
#Remove if the same version was made 
rm -r DIR_NAME > /dev/null

mkdir "$DIR_NAME"
cd "$DIR_NAME"

//Copy the binary
mkdir bin
mv ../main/webserv ./bin/webserv
#Copy the serup script
cp ../webserv_setup.sh .
cp ../readme.md .

mkdir webserv