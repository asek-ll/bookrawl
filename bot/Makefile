.PHONY: build-arm
build-arm: dist
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o dist/bookrawl-bot

dist:
	mkdir -p dist
