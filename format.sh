#!/bin/sh

# Basic format-all script. Does error. Hooray.
#find . -type f | grep \.go | grep -v '\./carbon.go' | xargs -0 go fmt
#go fmt carbon.go
gofmt -w .
