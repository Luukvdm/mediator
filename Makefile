.PHONY: vet
vet:
	go vet ./...

.PHONY: build
build:
	go build ./...

.PHONY: test
test:
	go test -shuffle=on ./...

.PHONY: coverage
cover:
	go test -coverprofile=coverage.out -covermode=atomic .
	go tool cover -html=coverage.out -o coverage.html

.PHONY: gen
gen:
	go generate ./...

.PHONY: lint
lint:
	golangci-lint run
