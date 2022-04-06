PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
VERSION = $(shell cat internal/version/version)
CUEM = go run ./cmd/cuem -v -p ./__examples__
COMMIT_SHA ?= $(shell git rev-parse --short HEAD)
TAG ?= $(VERSION)

cuem.k.show:
	$(CUEM) k show -o _output/nginx.yaml ./__examples__/clusters/demo/nginx.cue
	$(CUEM) k show ./__examples__/clusters/demo/nginx.cue

cuem.k.apply:
	$(CUEM) k apply ./__examples__/clusters/demo/nginx.cue

cuem.k.apply-as-release-template:
	$(CUEM) k apply --as-template ./__examples__/clusters/demo/nginx.cue

cuem.k.prune:
	$(CUEM) k prune ./__examples__/clusters/demo/nginx.cue

cuem.k.delete:
	$(CUEM) k delete ./__examples__/clusters/demo/nginx.cue

cuem.k.export:
	rm -rf ./_output
	$(CUEM) k show -o ./_output ./__examples__/clusters/demo/*.cue

cuem.eval:
	$(CUEM) eval -w -o _output/nginx.cue ./__examples__/clusters/demo/nginx.cue
	$(CUEM) eval -o nginx.yaml ./__examples__/clusters/demo/nginx.cue
	cue eval _output/nginx.cue
	$(CUEM) eval -o nginx.json _output/nginx.cue

cuem.fmt:
	$(CUEM) fmt -l -w ./...

cuem.get:
	$(CUEM) get ./...

gen.modutil:
	go generate ./pkg/modutil/internal

build:
	goreleaser build --snapshot --rm-dist

fmt:
	goimports -l -w .

test:
	go test -v -failfast ./pkg/...

cover:
	go test -v -coverprofile=coverage.txt -covermode=atomic ./pkg/...

install: build
	mv ./bin/cuem_$(shell go env GOOS)_$(shell go env GOARCH)/cuem ${GOPATH}/bin/cuem

tidy:
	go mod tidy

dep:
	go get -u ./...

debug:
	go test -v ./pkg/cuemod
	#tree ./pkg/cuemod/testdata/b/cue.mod

gen-deepcopy:
	deepcopy-gen \
		--output-file-base zz_generated.deepcopy \
		--go-header-file ./hack/boilerplate.go.txt \
		--input-dirs $(PKG)/pkg/apis/release/v1alpha1


PUSH ?= true
NAMESPACES ?= docker.io/octohelm
TARGETS ?= cuem

DOCKER_BUILDX_BUILD = docker buildx build \
	--label=org.opencontainers.image.source=https://github.com/octohelm/cuemod \
	--label=org.opencontainers.image.revision=$(COMMIT_SHA) \
	--platform=linux/arm64,linux/amd64

ifeq ($(PUSH),true)
	DOCKER_BUILDX_BUILD := $(DOCKER_BUILDX_BUILD) --push
endif


dockerx: build
	$(foreach target,$(TARGETS),\
		$(DOCKER_BUILDX_BUILD) \
		--build-arg=VERSION=$(VERSION) \
		$(foreach namespace,$(NAMESPACES),--tag=$(namespace)/$(target):$(TAG)) \
		--file=cmd/$(target)/Dockerfile . ;\
	)
