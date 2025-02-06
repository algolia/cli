# Build the binary
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
ARG VERSION=docker
RUN apk update && apk add --no-cache curl
RUN go mod download
RUN go install github.com/go-task/task/v3/cmd/task@latest
RUN task download-spec-file && VERSION=${VERSION} task build

FROM alpine
RUN apk update && apk add --no-cache ca-certificates 
COPY --from=builder /app/algolia /bin/algolia
ENTRYPOINT ["/bin/algolia"]
