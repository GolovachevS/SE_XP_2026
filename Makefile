.PHONY: fmt fmt-check lint test build proto

BINARY_NAME ?= chat
PROTO_FILES := api/chat/v1/chat.proto

fmt:
	gofmt -w .

ifeq ($(OS),Windows_NT)
fmt-check:
	@powershell -NoProfile -Command "$$files = gofmt -l .; if ($$files) { Write-Host 'gofmt check failed; run: gofmt -w .'; $$files; exit 1 }"
else
fmt-check:
	@sh -ec 'files="$$(gofmt -l .)"; if [ -n "$$files" ]; then echo "gofmt check failed; run: gofmt -w ."; echo "$$files"; exit 1; fi'
endif

lint:
	golangci-lint run --timeout=5m

test:
	go test ./...

build:
	go build -o bin/$(BINARY_NAME) ./cmd/chat

proto:
	PATH="$(PATH):$$(go env GOPATH)/bin" protoc -I . \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
