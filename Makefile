STEAMPIPE_INSTALL_DIR ?= ~/.steampipe
BUILD_TAGS             = netgo
PLUGIN_PATH            = hub.steampipe.io/plugins/ipgeolocation/ipgeolocation@latest
PLUGIN_BIN             = $(STEAMPIPE_INSTALL_DIR)/plugins/$(PLUGIN_PATH)/steampipe-plugin-ipgeolocation.plugin

CONFIG_DIR  = $(STEAMPIPE_INSTALL_DIR)/config
CONFIG_FILE = $(CONFIG_DIR)/ipgeolocation.spc

.PHONY: install install-config setup fmt lint test clean

## Build and install the plugin binary in one step
install:
	@mkdir -p $(dir $(PLUGIN_BIN))
	go build -o $(PLUGIN_BIN) -tags "${BUILD_TAGS}" *.go
	@echo "✓ Plugin installed → $(PLUGIN_BIN)"

## Copy the sample config (skips if already exists so you don't lose edits)
install-config:
	@mkdir -p $(CONFIG_DIR)
	@if [ -f "$(CONFIG_FILE)" ]; then \
		echo "⚠  Config already exists at $(CONFIG_FILE) — skipping."; \
		echo "   To reset it, delete the file and re-run: make install-config"; \
	else \
		cp config/ipgeolocation.spc $(CONFIG_FILE); \
		echo "✓ Config installed → $(CONFIG_FILE)"; \
		echo "  Edit the file to add your API key, or set IPGEOLOCATION_API_KEY."; \
	fi

## Full setup: dependencies + binary + config in one step
setup:
	@echo "→ Fetching Go dependencies..."
	go mod tidy
	@echo "→ Building and installing plugin..."
	@$(MAKE) install
	@echo "→ Installing config..."
	@$(MAKE) install-config
	@echo ""
	@echo "Done! Next steps:"
	@echo "  1. Edit $(CONFIG_FILE)  (or set IPGEOLOCATION_API_KEY)"
	@echo "  2. steampipe service restart"
	@echo "  3. steampipe query \"select * from ipgeolocation_ip where ip = '8.8.8.8'\""

## Format all Go source files
fmt:
	gofmt -w .

## Lint (requires golangci-lint)
lint:
	golangci-lint run ./...

## Run tests
test:
	go test ./...

## Remove built plugin binary
clean:
	@rm -f $(PLUGIN_BIN)
	@echo "✓ Cleaned"