# Makefile for Protocol Buffer formatting and validation

PROTO_FILES := $(wildcard *.proto)

# Install dependencies
.PHONY: install
install:
	@echo "Installing dependencies..."
	@brew install protobuf
	@brew install bufbuild/buf/buf
	@echo "Dependencies installed successfully!"

# Format proto files
.PHONY: format
format:
	@echo "Formatting .proto files..."
	@buf format -w .
	@echo "Formatting complete."

# Validate proto files
.PHONY: validate
validate:
	@echo "Validating .proto files..."
	@for file in $(PROTO_FILES); do \
		echo "Validating $$file..."; \
		protoc --proto_path=. --descriptor_set_out=/dev/null $$file || exit 1; \
	done
	@echo "Validation complete."

# Format and validate all
.PHONY: all
all: format validate
	@echo "All operations completed successfully!"
