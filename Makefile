build:
	@echo "Building CLI..."
	@go build -o dol cmd/cli/main.go

# Run the TUI application
run::
	@go run cmd/main.go

# Run the TUI with debugging enabled
# aka run TUI while using less on a log file
debug:
	@go run cmd/main.go 2> debug.log & tail -f debug

# generate sqlc files
sqlc:
	@sqlc generate -f internal/data/db/sqlc.yaml

# WARNING: This will delete all your local state!
state-reset:
	@rm -rf ~/.g7c
