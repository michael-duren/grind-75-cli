build:
	@echo "Building CLI..."
	@go build -o dol cmd/cli/main.go

# Run the TUI application
serve:
	@go run cmd/cli/main.go

# Run the TUI with debugging enabled
# aka run TUI while using less on a log file
debug:
	@go run cmd/cli/main.go 2> debug.log & tail -f debug
