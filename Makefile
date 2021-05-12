CUEM = go run ./cmd/cuem -v -p ./__examples__

cuem.k.show:
	$(CUEM) k show ./__examples__/clusters/demo/nginx.cue

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
	$(CUEM) eval -w -o _output/nginx.cue ./__examples__/components/nginx
	$(CUEM) eval -o nginx.yaml ./__examples__/components/nginx
	cue eval _output/nginx.cue
	$(CUEM) eval -o nginx.yaml _output/nginx.cue

cuem.fmt:
	$(CUEM) fmt -l -w ./...

cuem.get:
	$(CUEM) get ./...

build: download
	goreleaser build --snapshot --rm-dist

fmt:
	goimports -l -w .
	gofmt -l -w .

test:
	go test -v ./pkg/...

cover:
	go test -v -coverprofile=coverage.txt -covermode=atomic ./pkg/...

install: build
	mv ./bin/cuemod_$(shell go env GOOS)_$(shell go env GOARCH)/cuem ${GOPATH}/bin/cuem

dep:
	go get -u ./...
	go mod tidy

download:
	go mod download -x

setup:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/go-courier/husky/cmd/husky@latest
	#go install github.com/goreleaser/goreleaser/cmd@latest

debug:
	go test -v ./pkg/cuemod
	#tree ./pkg/cuemod/testdata/b/cue.mod
