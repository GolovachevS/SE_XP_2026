# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build

WORKDIR /src

# Install CA certs for module downloads
RUN apk add --no-cache ca-certificates

COPY go.mod ./

# If go.sum exists later, this layer will still be correct.
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/chat ./cmd/chat

FROM gcr.io/distroless/static:nonroot

WORKDIR /app
COPY --from=build /out/chat /app/chat

USER nonroot:nonroot
ENTRYPOINT ["/app/chat"]
