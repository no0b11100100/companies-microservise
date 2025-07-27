FROM golang:1.24-bullseye AS builder

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    librdkafka-dev \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN sh ./cmd/scripts/install_dependencies.sh

RUN go mod tidy

# RUN make lint
# RUN make generate-mocks
# RUN make unit-test

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o main ./cmd

FROM debian:bullseye-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
    librdkafka1 \
    ca-certificates \
 && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8080

CMD ["./main"]