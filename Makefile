APP_NAME=go-dci
GO_FILES=$(shell find . -type f -name '*.go' -not -path "./vendor/*")

.PHONY: vet build lint test clean run

vet:
	go vet $(GO_FILES)

build:
	go build -o $(APP_NAME)

lint:
	golangci-lint run ./cmd ./lib

test:
	go test -v ./cmd ./lib

clean:
	rm -f $(APP_NAME)

run: build
	./$(APP_NAME)

check-swagger-alignment:
	@echo "Checking API alignment with DCI API spec..."
	@go run ./scripts/check-swagger-alignment.go \
		--endpoints-file="./scripts/dci-endpoints.yaml" \
		--lib-path="./lib" \
		--base-url-var="DCIURL|BaseURL"

.PHONY: vet build lint test clean run check-swagger-alignment
