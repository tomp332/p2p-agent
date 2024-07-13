# Define the directories
PROTO_DIR := protos
GEN_DIR := src/protocol

# Ensure the generated directory exists
$(GEN_DIR):
	mkdir -p $(GEN_DIR)

# Find all .proto files
PROTO_FILES := $(shell find $(PROTO_DIR) -name '*.proto')

# Generate Go files from .proto files
generate: $(GEN_DIR)
	@echo "Generating Go files for all .proto files in the $(PROTO_DIR) directory..."
	@for proto in $(PROTO_FILES); do \
		protoc --go_out=$(GEN_DIR) --go-grpc_out=$(GEN_DIR) --proto_path=$(PROTO_DIR) \
		--go_opt=paths=source_relative --go-grpc_opt=paths=source_relative \
		$$proto; \
	done
	@echo "Generated Go files for all .proto files in the $(PROTO_DIR) directory."

# Phony targets
.PHONY: generate




build:
	@go build -o bin/fs

run: build
	@./bin/fs

test:
	@go test ./...

