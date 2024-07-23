CURRENT = $(shell pwd)
COVERAGE_FILE = coverage.out
COVERAGE_HTML = coverage.html
# Generate Go files from .proto files
generate:
	@echo "Generating Go files for all .proto files in the $(PROTO_DIR) directory..."
	@protoc -I=$(CURRENT) --go_out=$(CURRENT) --go-grpc_out=$(CURRENT)  $(CURRENT)/protos/*.proto
	@echo "Generated Go files for all .proto files in the $(PROTO_DIR) directory."

# Phony targets
.PHONY: generate

build:
	@go build -o bin/fs

run: build
	@./bin/fs

test:
	@echo "Running project tests..."
	@go test ./tests/...

coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=$(COVERAGE_FILE) ./...
	@go tool cover -html=$(COVERAGE_FILE) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILDDIR) $(COVERAGE_FILE) $(COVERAGE_HTML)