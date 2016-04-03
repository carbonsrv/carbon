# Compilers and stuff
GO?=go
GOFMT?=gofmt

# Vars
LUA_GLUE=$(wildcard builtin/*) $(wildcard builtin/3rdparty/*) $(wildcard builtin/libs/*) $(wildcard builtin/libs/wrappers/*)

all: carbon
carbon: modules/glue/generated_glue.go fmt
	$(GO) build -o $@ -v ./carbon.go

modules/glue/generated_glue.go: $(LUA_GLUE)
	$(GO) generate
	$(GOFMT) -w -s modules/glue

fmt:
	$(GOFMT) -w -s .

test:
	$(GO) test -v ./...

clean:
	rm -f carbon

.PHONY: carbon test

# Convenience stuff
repl: carbon
	./carbon -repl
