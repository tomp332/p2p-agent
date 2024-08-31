CURRENT = $(shell pwd)
COVERAGE_FILE = cover.out
COVERAGE_HTML = cover.html
MAIN_PACKAGE = ./pkg
COVER_PACKAGE = ./$(MAIN_PACKAGE)/...
MOCK_PACKAGE = ./tests/mocks

# Generate Go fsNode from .proto fsNode
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
	@rm -rf $(MOCK_PACKAGE)
	@mockgen -source=./pkg/pb/files_node_grpc.pb.go -destination ./tests/mocks/mock_fsNodeService.go --package mocks
	@mockgen -source=./pkg/storage/storage.go -destination ./tests/mocks/mock_storage.go --package mocks
	@mockgen -source=./tests/interfaces/fsNodeClient.go -destination ./tests/mocks/mock_fsNodeStreams.go --package mocks
	@mockgen -source=./pkg/nodes/client.go -destination ./tests/mocks/mock_nodeClient.go --package mocks
	@mockgen -source=./pkg/nodes/fsNode/client.go -destination ./tests/mocks/mock_fsNodeClient.go --package mocks
	@mockgen -source=./pkg/server/managers/base.go -destination ./tests/mocks/mock_authenticationManager.go --package mocks
	@echo "Finished generating test mocks"

tests:
	@echo "Running tests"
	@gotestsum --format testname ./tests/... -v

coverage:
	@echo "Running tests with coverage"
	@gotestsum --format testname ./tests/... -v -coverpkg=$(COVER_PACKAGE) -coverprofile=$(COVERAGE_FILE)


clean:
	@echo "Cleaning up..."
	@$(RM_RF) $(BUILDDIR) $(COVERAGE_FILE)