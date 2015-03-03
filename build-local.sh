#!/bin/bash
echo "Build resm binary"
set -e
GOPATH_TMP=`pwd`/.gohome
echo "set GOPATH=$GOPATH_TMP"
mkdir -p ${GOPATH_TMP}
export GOPATH=${GOPATH_TMP}
mkdir -p ${GOPATH}/src/github.com/nordicdyno/resm-sketch
cp -r * ${GOPATH}/src/github.com/nordicdyno/resm-sketch/
go get -v github.com/nordicdyno/resm-sketch/resm
