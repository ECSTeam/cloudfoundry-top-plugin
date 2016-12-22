#!/bin/bash

./compile.sh

RESULT=$?
if [ $RESULT -eq 0 ]; then
  echo compile success
  cf top "$@"
else
  echo compile failed
fi

