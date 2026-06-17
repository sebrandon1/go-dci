APP_NAME=go-dci
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo dev)

.PHONY: vet build lint test clean run

vet:
	go vet ./...

build:
	go build -ldflags "-X github.com/sebrandon1/go-dci/cmd.Version=$(VERSION)" -o $(APP_NAME)

lint:
	golangci-lint run ./cmd ./lib

test:
	go test -v ./cmd ./lib

coverage:
	go test -coverprofile=coverage.out ./cmd ./lib
	go tool cover -func=coverage.out

clean:
	rm -f $(APP_NAME) coverage.out

run: build
	./$(APP_NAME)

check-swagger-alignment:
	@echo "Checking API alignment with DCI API spec..."
	@go run ./scripts/check-swagger-alignment.go \
		--endpoints-file="./scripts/dci-endpoints.yaml" \
		--lib-path="./lib" \
		--base-url-var="DCIURL|BaseURL"

.PHONY: vet build lint test coverage clean run check-swagger-alignment
