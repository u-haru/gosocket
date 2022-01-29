GOOS	  = windows
GOARCH	  = amd64

TARGET	  = server client

GO		  = go
LDFLAGS   = -w -s

BINDIR    = .
BINS     := $(TARGET:%=$(BINDIR)/%)

GO_FILES:=$(shell find . -type f -name '*.go' -print)

.PHONY: all clean

all: $(BINS)
	@echo FINISHED!

$(BINS): $(GO_FILES)
	$(GO) build -ldflags='$(LDFLAGS)' $(@:%=./cmd/%)

client:TARGET=client

server:TARGET=server

clean:
	$(GO) clean
	rm -f $(BINS)