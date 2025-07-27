lint:
	golangci-lint run -v ./cmd
.PHONY: lint

generate-doc:
	go install github.com/swaggo/swag/cmd/swag@latest
	mkdir -p docs
	swag init --generalInfo cmd/main.go --output docs
.PHONY: generate-doc
