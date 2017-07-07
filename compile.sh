#!/bin/bash

GOOS=darwin go build -o top-plugin-osx
#GOOS=linux go build -o top-plugin-linux
#GOOS=windows GOARCH=amd64 go build -o top-plugin.exe
#GOOS=linux GOARCH=arm GOARM=6 go build -o top-plugin-linux-arm6
if [ $? != 0 ]; then
   printf "Error when executing compile\n"
   exit 1
fi
cf uninstall-plugin top
cf install-plugin -f ./top-plugin-osx
