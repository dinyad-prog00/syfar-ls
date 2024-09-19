LANG_SERVER_NAME := syfar-ls

BIN_DIR := ../vscode/extension/syfar/bin

# Define the Go build command
GO_BUILD := go build -o

# Platforms to build for
PLATFORMS := \
	darwin-amd64 \
	darwin-arm64 \
	linux-amd64 \
	windows-amd64

# Platform-specific binaries
bin/darwin-amd64/$(LANG_SERVER_NAME): 
	GOOS=darwin GOARCH=amd64 $(GO_BUILD) $(BIN_DIR)/darwin-amd64/$(LANG_SERVER_NAME)

bin/darwin-arm64/$(LANG_SERVER_NAME): 
	GOOS=darwin GOARCH=arm64 $(GO_BUILD) $(BIN_DIR)/darwin-arm64/$(LANG_SERVER_NAME)

bin/linux-amd64/$(LANG_SERVER_NAME): 
	GOOS=linux GOARCH=amd64 $(GO_BUILD) $(BIN_DIR)/linux-amd64/$(LANG_SERVER_NAME)

bin/windows-amd64/$(LANG_SERVER_NAME).exe: 
	GOOS=windows GOARCH=amd64 $(GO_BUILD) $(BIN_DIR)/windows-amd64/$(LANG_SERVER_NAME).exe

# Build all platforms
.PHONY: all
all: $(PLATFORMS)

# Define the build targets
darwin-amd64: bin/darwin-amd64/$(LANG_SERVER_NAME)
darwin-arm64: bin/darwin-arm64/$(LANG_SERVER_NAME)
linux-amd64: bin/linux-amd64/$(LANG_SERVER_NAME)
windows-amd64: bin/windows-amd64/$(LANG_SERVER_NAME).exe

# Clean up binaries
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)
