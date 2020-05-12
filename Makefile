.PHONY: all test build clean

all: clean test build

build: 
	mkdir -p build
	go build -o build ./...

test:
	go test -v -coverprofile=tests/results/cover.out ./...

cover:
	go tool cover -html=tests/results/cover.out -o tests/results/cover.html

clean:
	rm -rf build/*
	go clean ./...

container:
	podman build -t nexus-registry-nexus.apps.aws2-dev.ocp.14west.io/trackmate-couchbase-analytics:1.14.2 .

push:
	podman push nexus-registry-nexus.apps.aws2-dev.ocp.14west.io/trackmate-couchbase-analytics:1.14.2 
