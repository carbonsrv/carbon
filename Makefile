# Compilers and stuff
GO?=go
GOFMT?=gofmt
GO_BINDATA?=go-bindata
STRIP?=strip --strip-all
UPX?=upx --lzma -9
GIT?=git

# Vars
GLUE_DIRS=$(shell find ./builtin -type d | grep -v ".git" | grep -v "spec")
GLUE_FILES=$(shell find ./builtin -type f | grep -v ".git" | grep -v "spec")
GLUE_OUTPUT=modules/glue/generated_glue.go

all: carbon
carbon: submodules $(GLUE_OUTPUT) fmt
	$(GO) build -o $@ -v ./carbon.go

$(GLUE_OUTPUT): submodules $(GLUE_FILES)
	$(GO_BINDATA) -nomemcopy -o $(GLUE_OUTPUT) -pkg=glue -prefix "./builtin" $(GLUE_DIRS)
	$(GOFMT) -w -s modules/glue

submodules: .gitmodules builtin/libs/vfs
	$(GIT) submodule init
	$(GIT) submodule foreach "git checkout master && git pull"

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

glue: $(GLUE_OUTPUT)

dist: carbon
	$(STRIP) carbon
	$(UPX) carbon
