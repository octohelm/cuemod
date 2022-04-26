export GIT_SHA ?= $(shell git rev-parse HEAD)
export GIT_REF ?= HEAD

DAGGER = dagger --log-format=plain -p ./dagger
CUEM = go run ./cmd/cuem


tar:
	mkdir -p build/tar
	$(foreach n,$(shell ls build/output),\
		tar -czf "$(PWD)/build/tar/$(n).tar.gz" -C $(PWD)/build/output/$(n) .;)
.PHONY: tar

push:
	$(DAGGER) do push

build:
	$(DAGGER) do build
.PHONY: build

dagger.dep:
	$(CUEM) get ./dagger/...

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

dep:
	go get -u -t ./pkg/...

gen-deepcopy:
	deepcopy-gen \
		--output-file-base zz_generated.deepcopy \
		--go-header-file ./hack/boilerplate.go.txt \
		--input-dirs $(PKG)/pkg/apis/release/v1alpha1

cuem.k.show.pager:
	$(CUEM) k show ./__examples__/clusters/demo/nginx.cue

cuem.k.show:
	$(CUEM) k show -o _output/nginx0.yaml ./__examples__/clusters/demo/nginx.cue
	$(CUEM) k show -o _output/nginx1.yaml ./__examples__/components/nginx '{ #values: image: tag: "latest", #context: "crpe-test" }'

cuem.k.apply:
	$(CUEM) k apply ./__examples__/clusters/demo/nginx.cue

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
	$(CUEM) get -i=go k8s.io/api k8s.io/apimachinery
	$(CUEM) get ./...
