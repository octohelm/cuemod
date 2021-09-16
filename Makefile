PKG = $(shell cat go.mod | grep "^module " | sed -e "s/module //g")
CUEM = go run ./cmd/cuem -v -p ./__examples__

up.operator:
	WATCH_NAMESPACE=default \
		go run ./cmd/cuem-operator/main.go

cuem.k.show:
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

build: gen.modutil
	goreleaser build --snapshot --rm-dist

fmt:
	goimports -l -w .

test:
	go test -v ./pkg/...

cover:
	go test -v -coverprofile=coverage.txt -covermode=atomic ./pkg/...

install: build
	mv ./bin/cuem_$(shell go env GOOS)_$(shell go env GOARCH)/cuem ${GOPATH}/bin/cuem

dep:
	go get -u ./...
	go mod tidy

setup:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/go-courier/husky/cmd/husky@latest
	#go install github.com/goreleaser/goreleaser/cmd@latest

debug:
	go test -v ./pkg/cuemod
	#tree ./pkg/cuemod/testdata/b/cue.mod


gen-deepcopy:
	deepcopy-gen \
		--output-file-base zz_generated.deepcopy \
		--go-header-file ./hack/boilerplate.go.txt \
		--input-dirs $(PKG)/pkg/apis/release/v1alpha1

dockerx:
	$(foreach target,$(TARGETS),\
		$(DOCKER_BUILDX_BUILD) \
		--build-arg=VERSION=$(VERSION) \
		$(foreach namespace,$(NAMESPACES),--tag=$(namespace)/$(target):$(TAG)) \
		--file=cmd/$(target)/Dockerfile . ;\
	)
