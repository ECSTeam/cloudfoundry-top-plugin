#!/bin/bash

GOOS=darwin go build -o top-plugin-osx
if [ $? != 0 ]; then
   printf "Error when executing compile\n"
   exit 1
fi
cf uninstall-plugin TopPlugin
cf install-plugin -f ./top-plugin-osx
