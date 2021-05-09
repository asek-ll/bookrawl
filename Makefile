

.PHONY: build
build: dist
	cd app && GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ../dist/bookrawl

dist:
	mkdir -p dist
