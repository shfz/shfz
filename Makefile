pre:
	@make fmt
	@make lint
	@make build
	@make test

build:
	go build

fmt:
	gofmt -s -w .

lint:
	golangci-lint run --tests --enable=goimports

test:
	go test -v -race -cover ./...

up:
	rm go.sum
	go get -u
	go mod tidy
	@make pre
