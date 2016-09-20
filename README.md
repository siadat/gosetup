# gosetup

Gosetup is an simple interactive commandline tool to setup a Go workspace.

This tool is inspired by [The Go Workbench](https://docs.google.com/presentation/d/1mUIX3btCiGPguqJOE4h9HDoOW3VyhU2-tzXufz9PqQ0/edit#slide=id.p) slides by adg.

## Description

Gosetup performs the following actions:

- set `GOPATH`
- add `GOPATH/bin` to `PATH`
- create the following directories
  - `GOPATH/src`
  - `GOPATH/bin`
  - `GOPATH/pkg`
- create a hello world package with the following file
  - `GOPATH/src/github.com/USERNAME/hello/hello.go`

## Output

This is the output of a simple session:

    >>> Adding the following lines to /home/sina/.bashrc:
    >>> 
    >>>     export GOPATH="/home/sina/go"
    >>>     export PATH="$PATH${PATH:+:}$GOPATH/bin"
    >>> 
    Continue [Y,n,e,?]? y
    >>> Done. Changes will be reflected next time a terminal is started.
    
    >>> Creating a hello world program with the following file:
    >>> 
    >>>     /home/sina/go/src/github.com/$USERNAME/hello/hello.go
    >>> 
    Enter your GitHub username or exit: siadat
    >>> Done.
    
    >>> Run this program with the following commands:
    >>> 
    >>>     cd "/home/sina/go/src/github.com/siadat/hello"
    >>>     go run hello.go
    >>>     go install
    >>>     hello
    >>> 
