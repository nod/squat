
default: build/squat

PACKAGES := $(shell go list -f {{.Dir}} ./...)
GOFILES  := $(addsuffix /*.go,$(PACKAGES))
GOFILES  := $(wildcard $(GOFILES))

VER =
ifndef VER
	VER := $(shell ./bin/incr_build ./VERSION)
endif

setuplocal:
	mkdir -p build/tmp

fmt:
	go fmt -x main.go
	go fmt -x squat/*.go

tidy-vendor:
	# we want to ensure that we have the proper libraries in place
	go mod vendor
	# just keep things clean
	go mod tidy

build/squat: setuplocal tidy-vendor main.go
	$(info setting VER to $(VER))
	go build -ldflags "-X main.version=$(VER)" -mod=vendor -o build/squat main.go

run: build/squat
	./bin/run_local

clean: setuplocal
	go clean
	rm -rf build

test: export TMPDIR=build/tmp
test: export CGO_ENABLED=0
test: setuplocal
	go test -v ./...

binaries: setuplocal tidy-vendor \
	build/squat-linux-amd64 build/squat-macos-amd64 build/squat-linux-arm64

build/squat-linux-amd64: $(GOFILES)
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VER)" -mod=vendor -o build/squat-linux-amd64 main.go

build/squat-macos-amd64: $(GOFILES)
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VER)" -mod=vendor -o build/squat-macos-amd64 main.go

build/squat-linux-arm64: $(GOFILES)
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VER)" -mod=vendor -o build/squat-linux-arm64  main.go


tagproj:
	git tag -a v`cat VERSION`
	git push --tags


