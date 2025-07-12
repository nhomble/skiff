PROG 		:= skiff
LDFLAGS		:= -w -s
GOOS		?= $(shell go env GOOS)
GOARCH		?= $(shell go env GOARCH)
CGO_ENABLED	?= 0

# Output file with platform suffix if cross-compiling
ifeq ($(GOOS),windows)
    OUTPUT = $(PROG)-$(GOOS)-$(GOARCH).exe
else
    OUTPUT = $(PROG)-$(GOOS)-$(GOARCH)
endif

# Use simple name for native build
ifeq ($(GOOS),$(shell go env GOOS))
ifeq ($(GOARCH),$(shell go env GOARCH))
    OUTPUT = $(PROG)
endif
endif

.PHONY: build image test/unit clean install fmt help test/policy test benchmark example demo build-all

all: build

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build -ldflags="$(LDFLAGS)" -o $(OUTPUT) ./cmd/skiff

image:
	docker build -t $(PROG) .

fmt:
	go fmt ./...

clean:
	rm -f $(PROG) $(PROG)-*

test/unit:
	go test ./...

test/policy: build
	./test/test-policies/test-policy.sh

test: build test/unit test/policy
	@echo "All tests passed"

benchmark: build
	./test/performance/benchmark.sh