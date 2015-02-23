# Description

This code is my sketch of RESTish API in Go code

I've learned here Go's testing techniques, interfaces and some net/http tricks.

Also I've gotten my feet wet with Bolt DB (spoiler: mention -file flag)

More detailed documentation is here:
https://docs.google.com/document/d/18tENtluzSS3bKygIJ4MWBv-GkuTh8V0hG95bCFu5qWI/edit?usp=sharing

# Where to test

Working example is accessible here: http://resm-sketch.tk
It is powered with Google Compute Engine, Docker and semi-free .TK domain

# How to build binary manually

## Go way

1. Install Go (brew, apt-get or manually http://golang.org/doc/install)
2. Set GOHOME i.e. to ~/gohome
3. Fetch and compile binary: `go get github.com/nordicdyno/resm-sketch`
4. Run it ~/gohome/bin/resm -limit=5 -verbose

You also need git and mercurial

## Make way

1-2 Steps the same as in "Go way"
3. Get sources by git `git clone git@github.com:nordicdyno/resm-sketch.git`
4. Build and resolve dependencies `cd resm-sketch && go get .`
5. Rebuild and run binary with `make`

Hint: To run in persistent mode use target `run_bolt`: `make run_bolt`

# How to test

 1. Make sure you have done build stage
 2. Run tests in sources root: `make test`


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

_I am using supervisor as system runner and process manager for resm.
Main reason for supervisor – I don't want to write init.d scripts.
Because it's clumsy, error prone and not fun.
Another reason - I just know supervisor well enough._

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