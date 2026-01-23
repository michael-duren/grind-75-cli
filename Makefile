build:
	@echo "Building CLI..."
	@go build -o g7c cmd/main.go

# Run the TUI application (compiles and runs)
run::
	@go run cmd/main.go

# Build and run the binary (faster if already built, or just cleaner)
run-bin: build
	@./g7c

# watch the log files
debug:
	@less +G "$$(ls -t ~/.g7c/logs/g7c-*.log | head -n 1)"

# generate sqlc files
sqlc:
	@sqlc generate -f internal/data/db/sqlc.yaml

# WARNING: This will delete all your local state!
state-reset:
	@rm -rf ~/.g7c
