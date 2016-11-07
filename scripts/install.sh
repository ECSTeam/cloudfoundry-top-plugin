#!/bin/bash

set -e

(cf uninstall-plugin "TopPlugin" || true) && go build -o top-plugin main.go && cf install-plugin top-plugin
