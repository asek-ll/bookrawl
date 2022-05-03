.PHONY: build-arm
build-arm: dist
	GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -o dist/bookrawl-bot-arm64

.PHONY: build-amd
build-amd: dist
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o dist/bookrawl-bot-amd64

.PHONY: build-amd-docker
build-amd-docker: build-amd
	docker build -f docker/Dockerfile-amd64 -t sopplet/bookrawl-bot-amd64 .

dist:
	mkdir -p dist

cert:
	openssl req -x509 -newkey rsa:2048 -keyout dist/key.pem -out dist/cert.pem -days 3560 -subj "/C=US/ST=New York/L=Brooklyn/O=Example Brooklyn Company/CN=${BOOKRAWL_BOT_CN}" -nodes
