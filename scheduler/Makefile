

.PHONY: build
build: dist
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/bookrawl

.PHONY: build-arm
build-arm: dist
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o dist/bookrawl

.PHONY: build-amd
build-amd: dist
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/bookrawl-fetcher-amd64

.PHONY: build-amd-docker
build-amd-docker: build-amd
	docker build -f docker/Dockerfile-amd64 -t sopplet/bookrawl-fetcher-amd64 .

dist:
	mkdir -p dist
