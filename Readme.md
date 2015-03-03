# Description

This code is my sketch of RESTish API in Go language.

I've trained Go's testing techniques, interfaces and some net/http tricks
while implementing this.

Also I've gotten my feet wet with Bolt DB (spoiler: mention -file flag), FPM,
Google Compute Engine and domain names registration process.

More detailed documentation is here:
https://docs.google.com/document/d/18tENtluzSS3bKygIJ4MWBv-GkuTh8V0hG95bCFu5qWI/edit?usp=sharing

# Where to try

Working example is accessible here: http://resm-sketch.tk

It is powered with Google Compute Engine, Docker and semi-free .TK domain

# How to build binary manually

## Go way

1. Install Go (brew, apt-get or manually http://golang.org/doc/install)
2. Set GOPATH. Primer: `mkdir ~/ghome; export GOPATH=~/ghome`
3. Fetch and compile binary: `go get github.com/nordicdyno/resm-sketch/resm`
4. Run: `$GOPATH/bin/resm -limit=5 -verbose`


_Hint: you also need to have git and mercurial installed_

## 'Make in local dir'-way

1. Get sources by git `git clone git@github.com:nordicdyno/resm-sketch.git`
2. `cd resm-sketch`
3. Build and run binary with `make`

_Hint: To run in persistent mode use target `run_bolt`: `make run_bolt`_

# How to test

Just run tests in sources root: `make test`

# Run locally

Get help:

    $GOHOME/bin/resm -h

In memory storage example:

    resm -bind=":9090" -limit=5

Persistent storage example:

    resm -bind=":9090" -limit=5 -file=resources.db

Test with `httpie` (https://github.com/jakubroztocil/httpie):

    http -v http://localhost:9090/allocate/me
    http -v http://localhost:9090/list/

# How to build debian package

    make docker_build_deb

*Requires Docker 1.5*

Deb package should be created in `/root` dir in container `resm-deb-builder`

FPM util used for deb creation, here are useful topic links:
* https://github.com/jordansissel/fpm
* https://github.com/jordansissel/fpm/wiki
* https://github.com/jordansissel/fpm/wiki/Debuild-to-fpm


# Example how to run binary on Debian with process manger

    make docker_supervisor

*Requires Docker 1.5*

_I use supervisor as a system runner and process manager for resm.
Main reason for supervisor – I don't want to write init.d scripts.
Because it's clumsy, error prone and not fun.
Another reasons:
* I just know supervisor well enough
* it's a generally good idea to use process manager for Go services like that._

# TODO

## Code improvements

### remove gorilla deps

* Use goji's sessions approach
* Use lightweight route lib

Related links:

* https://elithrar.github.io/article/map-string-interface/
* https://justinas.org/writing-http-middleware-in-go/
* https://github.com/julienschmidt/httprouter

### Refactor tests

* join repetitive code for storages
* add more edge cases
* add benchmarks?
* more static tools?

### Better docs

just do it for great good :)

##