.PHONY: fmt fmt-check lint test cover build proto

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

cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
ifeq ($(OS),Windows_NT)
	@powershell -NoProfile -Command "Start-Process coverage.html"
else
	@if [ "$$(uname)" = "Darwin" ]; then \
		open coverage.html; \
	else \
		xdg-open coverage.html; \
	fi
endif

build:
	go build -o bin/$(BINARY_NAME) ./cmd/chat

proto:
ifeq ($(OS),Windows_NT)
	@powershell -NoProfile -Command "$$gopath = (go env GOPATH); $$env:PATH = $$env:PATH + ';' + $$gopath + '\\bin'; protoc -I . --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $(PROTO_FILES)"
else
	@sh -ec 'PATH="$$PATH:$(shell go env GOPATH)/bin" protoc -I . --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative $(PROTO_FILES)'
endif
