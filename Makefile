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
	@mockgen -source=$(MAIN_PACKAGE)/nodes/base_node.go -destination=$(MOCK_PACKAGE)/mock_base_node.go -package=mocks
	@mockgen -source=$(MAIN_PACKAGE)/storage/storage.go -destination=$(MOCK_PACKAGE)/mock_storage.go -package=mocks
	@mockgen -source=$(MAIN_PACKAGE)/server/managers/base.go -destination=$(MOCK_PACKAGE)/mock_authentication_manager.go -package=mocks
	@mockgen -source=$(MAIN_PACKAGE)/pb/files_node_grpc.pb.go -destination=$(MOCK_PACKAGE)/mock_files_node_service.go -package=mocks
	@mockgen -source=$(MAIN_PACKAGE)/nodes/fsNode/client.go -destination=$(MOCK_PACKAGE)/mock_file_node_client.go -package=mocks
	@echo "Finished generating test mocks"

test:
	@echo "Running tests"
	@gotestsum --format testname ./tests/... -v

coverage:
	@echo "Running tests with coverage"
	@gotestsum --format testname ./tests/... -v -coverpkg=$(COVER_PACKAGE) -coverprofile=$(COVERAGE_FILE)


clean:
	@echo "Cleaning up..."
	@$(RM_RF) $(BUILDDIR) $(COVERAGE_FILE)