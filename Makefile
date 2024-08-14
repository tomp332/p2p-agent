CURRENT = $(shell pwd)
COVERAGE_FILE = cover.out
COVERAGE_HTML = cover.html
COVER_PACKAGE = ./src/...
MOCK_PACKAGE = ./tests/mocks
# Generate Go file_node from .proto file_node
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

mocks:
	@echo "Generating test mocks"
	@mockgen -source=./src/interfaces.go -destination=$(MOCK_PACKAGE)/mock_interfaces.go -package=mocks
	@mockgen -source=./src/pb/files_node_grpc.pb.go -destination=$(MOCK_PACKAGE)/mock_files_node_service.go -package=mocks
	@echo "Finished generating test mocks"

test:
	@echo "Running tests with coverage"
	@go test ./tests/... -v -coverpkg=$(COVER_PACKAGE) -coverprofile=$(COVERAGE_FILE)
	@echo "Finished running tests and coverage"

clean:
	@echo "Cleaning up..."
	@rm -rf $(BUILDDIR) $(COVERAGE_FILE)