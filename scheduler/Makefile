

.PHONY: build
build: dist
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/bookrawl

dist:
	mkdir -p dist
