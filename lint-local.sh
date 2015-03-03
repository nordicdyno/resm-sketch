#!/bin/bash
set -e
echo "Run lint"
GOPATH_TMP=`pwd`/.gohome
export GOPATH=${GOPATH_TMP}
go get -v github.com/golang/lint/golint
find * -name '*.go' -exec ${GOPATH}/bin/golint {} \;
