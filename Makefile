# Makefile for the boilerplate application
#

.DEFAULT_GOAL := help

# Go settings
GOCMD=go
GOMAIN=./cmd/main.go
BINARY_NAME=boilerplate

# Phony targets don't represent files
.PHONY: all build clean run-server help

# The `all` target is a convention, often used to build everything.
all: build

# Build the application binary.
# This compiles the application and places the output binary in the root directory.
build:
	@echo "==> Building application..."
	@$(GOCMD) build -o $(BINARY_NAME) $(GOMAIN)

# Run the web server for development.
# This command uses `go run` to compile and run the application without creating a permanent binary.
# It passes the `run` argument to execute the server command.
run-server:
	@echo "==> Starting web server..."
	@$(GOCMD) run $(GOMAIN) run

# Clean removes the built binary.
clean:
	@echo "==> Cleaning up..."
	@if [ -f $(BINARY_NAME) ]; then rm $(BINARY_NAME); fi

# Variables
# Ensure your DATABASE_URL environment variable is set, or dbmate will fail.
# You can optionally set it here if it's static and not sensitive,
# or load it from a .env file (though direct env var is often preferred for dbmate).
# Example: export DATABASE_URL ?= postgresql://user:password@host:port/dbname?sslmode=disable

MIGRATIONS_DIR ?= ./internal/db/migrations
SCHEMA_FILE    ?= ./internal/db/schema.sql

# Base dbmate command with common options
DBMATE_CMD = dbmate --migrations-dir "${MIGRATIONS_DIR}" --schema-file "${SCHEMA_FILE}"

# SQLC command (ensure sqlc is installed and in your PATH, or adjust as needed)
# e.g., for Go projects: SQLC_CMD = go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest
SQLC_CMD ?= sqlc

# --- DBMate Targets ---

help:
	@echo "Makefile for Boilerplate Application"
	@echo ""
	@echo "Usage: make [target] [ARGS...]"
	@echo ""
	@echo "General Targets:"
	@echo "  all           Build the application binary (default)."
	@echo "  build         Compile the application."
	@echo "  run-server    Start the web server."
	@echo "  clean         Remove the compiled binary."
	@echo "  dev           Start the full development environment (backend, frontend)"
	@echo ""
	@echo "DBMate Database Migration Targets:"
	@echo "  Ensure DATABASE_URL environment variable is set for dbmate commands."
	@echo "  db-create                      Create the database (if it doesn't exist)"
	@echo "  db-drop                        Drop the database (USE WITH CAUTION!)"
	@echo "  db-wait                        Wait for the database to become available"
	@echo "  migrate-create                 Create a new migration file (interactive prompt)"
	@echo "  migrate-up                     Apply all pending migrations"
	@echo "  migrate                        Alias for migrate-up"
	@echo "  migrate-down [COUNT=N]         Rollback N migrations (default is 1 if COUNT is not specified, e.g., make migrate-down COUNT=3)"
	@echo "  migrate-rollback               Rollback the last migration (alias for migrate-down with COUNT=1)"
	@echo "  migrate-status                 Show migration status"
	@echo "  migrate-dump                   Dump database schema to ${SCHEMA_FILE}"
	@echo ""
	@echo "SQLC Code Generation Target:"
	@echo "  Ensure sqlc is installed and a sqlc.yaml (or sqlc.yml) configuration file exists."
	@echo "  sqlc-generate                  Generate code using sqlc (runs migrate-dump first)"

# Default target
.DEFAULT_GOAL := help

.PHONY: migrate-create
migrate-create:
	@echo "Creating a new migration..."
	@read -p "Enter migration name (use underscores): " MIGRATION_NAME; \
	if [ -z "$$MIGRATION_NAME" ]; then \
		echo "Error: Migration name cannot be empty."; \
		exit 1; \
	fi; \
	echo "Creating migration: $$MIGRATION_NAME..."; \
	$(DBMATE_CMD) new $$MIGRATION_NAME

.PHONY: migrate-up
migrate-up:
	@echo "Applying pending migrations..."
	@$(DBMATE_CMD) up

.PHONY: migrate
migrate: migrate-up

.PHONY: migrate-down
migrate-down:
	$(eval COUNT_ARG := $(if $(COUNT),$(COUNT),1))
	@echo "Rolling back $(COUNT_ARG) migration(s)..."
	@$(DBMATE_CMD) down $(COUNT_ARG)

.PHONY: migrate-rollback
migrate-rollback:
	@echo "Rolling back the last migration..."
	@$(DBMATE_CMD) rollback

.PHONY: migrate-status
migrate-status:
	@echo "Checking migration status..."
	@$(DBMATE_CMD) status

.PHONY: migrate-dump
migrate-dump:
	@echo "Dumping schema to ${SCHEMA_FILE}..."
	@$(DBMATE_CMD) dump

.PHONY: db-create
db-create:
	@echo "Creating database (if it doesn't exist)..."
	@$(DBMATE_CMD) create

.PHONY: db-drop
db-drop:
	@echo "WARNING: This will drop the entire database!"
	@read -p "Are you sure you want to drop the database? (yes/NO): " CONFIRM; \
	if [ "$${CONFIRM}" = "yes" ]; then \
		echo "Dropping database..."; \
		$(DBMATE_CMD) drop; \
	else \
		echo "Database drop cancelled."; \
	fi

.PHONY: db-wait
db-wait:
	@echo "Waiting for database to become available..."
	@$(DBMATE_CMD) wait

# --- SQLC Target ---

.PHONY: dev
dev:
	@echo "==> Starting development environment..."
	@# Start the Go server in the background and capture its PID
	@make run-server &
	@BE_PID=$!
	@echo "Backend server started with PID: $BE_PID"
	@# Start the frontend server in the background and capture its PID
	@(cd frontend && npm run dev) &
	@FE_PID=$!
	@echo "Frontend server started with PID: $FE_PID"
	@# Set a trap: when this script EXITS (for any reason), kill the PIDs.
	@# The '|| true' prevents errors if the process is already gone.
	@trap 'echo "Stopping all processes..."; kill $BE_PID $FE_PID || true' EXIT
	@echo "==> Dev environment started. Press Ctrl+C to stop all processes."
	@# Wait for any process to exit. The trap will handle cleanup for all.
	@wait -n

.PHONY: sqlc-generate
sqlc-generate: migrate-dump
	@echo "Generating code with sqlc..."
	@echo "Ensure your sqlc configuration (e.g., sqlc.yaml) points to '${SCHEMA_FILE}' for the schema and specifies your query files."
	@${SQLC_CMD} generate

# --- Makefile hygiene ---
# Prevent deletion of intermediate files for chained rules
.SECONDARY:
# Disable built-in rules
.SUFFIXES:
