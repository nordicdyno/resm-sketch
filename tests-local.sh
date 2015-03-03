#!/bin/bash
set -e
echo "Run resm tests: $@"
GOPATH_TMP=`pwd`/.gohome
export GOPATH=${GOPATH_TMP}
go get golang.org/x/tools/cmd/cover

go test $@ github.com/nordicdyno/resm-sketch/store/inmemory
go test $@ github.com/nordicdyno/resm-sketch/store/inbolt
go test $@ github.com/nordicdyno/resm-sketch/resm
