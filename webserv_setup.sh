#!/bin/bash

#Check if we are running as sudo
if [ "$EUID" -ne 0 ]
  then echo "Install script needs to be run as sudo."
  exit
fi

WEBSERV_DIRECTORY="/var/lib/webserv"

#If data directory is not present, then create it
if [ -d "$WEBSERV_DIRECTORY" ]; then
  echo "Creating directory $WEBSERV_DIRECTORY."
  mkdir "$WEBSERV_DIRECTORY"
fi

echo "Moving static files to $WEBSERV_DIRECTORY"
#Copy static resources to the webserv directory
mv ./resources/* "$WEBSERV_DIRECTORY"

BIN_DIR="/usr/local/bin"

echo "Moving binaries to $BIN_DIR"
#Move the binary to the bin folder for easy setup
mv ./bin/* "$BIN_DIR"

#Path does not contain bin_dir 
echo "PATH does not contain $BIN_DIR => adding it to PATH"
if [[ $PATH != *"$BIN_DIR"* ]]; then
    export PATH=$PATH:$BIN_DIR && 
    echo "export PATH=$PATH:$BIN_DIR" >> ~/.bashrc
fi
