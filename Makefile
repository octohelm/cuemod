export GIT_SHA ?= $(shell git rev-parse HEAD)
export GIT_REF ?= HEAD

DAGGER = dagger --log-format=plain
CUEM = go run ./cmd/cuem

tar:
	mkdir -p build/tar
	$(foreach n,$(shell ls build/output),\
		tar -czf "$(PWD)/build/tar/$(n).tar.gz" -C $(PWD)/build/output/$(n) .;)
.PHONY: tar

push:
	$(DAGGER) do push
.PHONY: push

build:
	$(DAGGER) do build
.PHONY: build

dagger.dep:
	$(CUEM) get ./cuepkg/...

INTERNAL_FORK = go run ./tool/internalfork

fork.go.internal:
	$(INTERNAL_FORK) \
		-p cmd/go/internal/modload \
		-p cmd/go/internal/modfetch \
		-p internal/execabs \
		./pkg/modutil/internal

install:
	$(DAGGER) do build $(shell go env GOOS) $(shell go env GOARCH)
	mv ./build/output/cuem_$(shell go env GOOS)_$(shell go env GOARCH)/cuem ${GOPATH}/bin/cuem

fmt:
	goimports -l -w .

tidy:
	go mod tidy

test:
	go test -v -failfast ./pkg/...

dep:
	go get -u -t ./pkg/...

cuem.eval:
	$(CUEM) eval -w -o _output/nginx.cue ./__examples__/clusters/demo/nginx.cue
	$(CUEM) eval -o nginx.yaml ./__examples__/clusters/demo/nginx.cue
	cue eval _output/nginx.cue
	$(CUEM) eval -o nginx.json _output/nginx.cue

cuem.fmt:
	$(CUEM) fmt -l -w ./...

cuem.get:
	$(CUEM) get -i=go k8s.io/api k8s.io/apimachinery
	$(CUEM) get ./...
