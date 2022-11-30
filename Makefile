export GIT_SHA ?= $(shell git rev-parse HEAD)
export GIT_REF ?= HEAD

DAGGER=dagger
CUEM = go run ./cmd/cuem


INTERNAL_FORK = go run ./tool/internalfork

fork.go.internal:
	$(INTERNAL_FORK) \
		-p cmd/go/internal/modload \
		-p cmd/go/internal/modfetch \
		-p internal/godebug \
		./internal

install:
	$(DAGGER) do go archive
	mv ./build/output/cuem_$(shell go env GOOS)_$(shell go env GOARCH)/cuem ${GOPATH}/bin/cuem

fmt:
	goimports -l -w .

tidy:
	go mod tidy

test:
	go test -failfast ./pkg/...

dep:
	go get -u -t ./cmd/...

cuem.fmt:
	$(CUEM) fmt -l -w ./...

cuem.get.u:
	$(CUEM) get -u ./...

cuem.get:
	$(CUEM) get -i=go k8s.io/api k8s.io/apimachinery
	$(CUEM) get github.com/innoai-tech/runtime@main

