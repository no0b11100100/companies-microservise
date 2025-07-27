lint:
	golangci-lint run -v ./cmd
.PHONY: lint

generate-doc:
	go install github.com/swaggo/swag/cmd/swag@latest

	go mod tidy

	mkdir -p docs
	swag init --generalInfo cmd/main.go --output docs
.PHONY: generate-doc

generate-mocks:
	go generate ./cmd/internal/server/handlers/createRecordHandler.go
	go generate ./cmd/internal/server/handlers/deleteRecordHandler.go
	go generate ./cmd/internal/server/handlers/getRecordHandler.go
	go generate ./cmd/internal/server/handlers/updateRecordHandler.go
	go generate ./cmd/internal/eventSender/sender.go
	go generate ./cmd/internal/database/database.go
	go generate ./cmd/internal/server/server.go
	go generate ./cmd/app.go
.PHONY: generate-mocks

unit-test:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
.PHONY: unit-test

integration-test:
	# apt-get install -y --no-install-recommends build-essential librdkafka-dev ca-certificates
	# go get github.com/testcontainers/testcontainers-go
	# go get github.com/testcontainers/testcontainers-go/modules/mysql
	# go get github.com/testcontainers/testcontainers-go/modules/kafka
	# go get github.com/stretchr/testify

	# apt-get update && apt-get install -y --no-install-recommends \
	# 	build-essential \
	# 	librdkafka-dev \
	# 	ca-certificates \
	# 	curl

	# # Kafka (confluent)
	# go get github.com/confluentinc/confluent-kafka-go/kafka

	# # Testcontainers core + modules (MySQL, Kafka, etc.)
	# go get github.com/testcontainers/testcontainers-go
	# go get github.com/testcontainers/testcontainers-go/modules/mysql
	# go get github.com/testcontainers/testcontainers-go/modules/kafka

	# # Optional: testify for easier assertions
	# go get github.com/stretchr/testify

	# go mod tidy


	./cmd/tests/test_env.sh

	docker-compose up --build -d

	go test -v ./cmd/tests/...

	docker-compose down -v

.PHONY: generate-doc
