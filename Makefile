lint:
	golangci-lint run -v ./cmd
.PHONY: lint

generate-doc:
	mkdir -p docs
	swag init --generalInfo cmd/main.go --output docs
.PHONY: generate-doc

generate-mocks:
	mkdir -p cmd/tests
	mkdir -p cmd/tests/mocks
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
	go test ./cmd/... -coverprofile=coverage.out
	go tool cover -func=coverage.out
.PHONY: unit-test

integration-test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit --exit-code-from testapp
.PHONY: integration-test
