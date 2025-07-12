PROG 	:= skiff

.PHONY: build test/unit clean install fmt help test/policy test benchmark example demo

all: build

build:
	go build -o $(PROG) ./cmd/skiff

fmt:
	go fmt ./...

clean:
	rm -f $(PROG)

test/unit:
	go test ./...

test/policy: build
	./test/test-policies/test-policy.sh

test: build test/unit test/policy
	@echo "All tests passed"

benchmark: build
	./test/performance/benchmark.sh