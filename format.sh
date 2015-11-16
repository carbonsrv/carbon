#!/bin/sh

# Basic format-all script. Does error. Hooray.
find . -type f | grep \.go | xargs -0 go fmt
