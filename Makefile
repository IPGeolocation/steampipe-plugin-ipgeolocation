PLUGIN_NAME    = steampipe-plugin-ipgeolocation
PLUGIN_DIR     = $(HOME)/.steampipe/plugins/local/ipgeolocation
PLUGIN_BIN     = $(PLUGIN_DIR)/ipgeolocation.plugin
CONFIG_DIR     = $(HOME)/.steampipe/config
CONFIG_FILE    = $(CONFIG_DIR)/ipgeolocation.spc

.PHONY: all build install install-config setup fmt lint test clean

## Default: build the plugin binary
all: build

## Compile the plugin binary
build:
	go build -o $(PLUGIN_NAME) .

## Install binary into Steampipe's local plugin path
install: build
	@mkdir -p $(PLUGIN_DIR)
	@cp $(PLUGIN_NAME) $(PLUGIN_BIN)
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
	@echo "→ Building plugin..."
	@$(MAKE) install
	@echo "→ Installing config..."
	@$(MAKE) install-config
	@echo ""
	@echo "Done! Next steps:"
	@echo "  1. Edit ~/.steampipe/config/ipgeolocation.spc  (or set IPGEOLOCATION_API_KEY)"
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

## Remove built binary
clean:
	@rm -f $(PLUGIN_NAME)
	@echo "✓ Cleaned"