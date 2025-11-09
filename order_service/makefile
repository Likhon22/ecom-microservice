# =============================================================================
# ENVIRONMENT SETUP
# =============================================================================
ifeq (,$(wildcard .env))
$(error ".env file not found")
endif

include .env
export

# =============================================================================
# MIGRATION COMMANDS
# =============================================================================
new-migration:
	@echo "Creating new migration: $(name)..."
	@sql-migrate new $(name) -env=development

migrate-up:
	@echo "Running migrations UP..."
	@sql-migrate up -env=development

migrate-down:
	@echo "Running migrations DOWN..."
	@sql-migrate down -env=development

migrate-status:
	@echo "Checking migration STATUS..."
	@sql-migrate status -env=development

# =============================================================================
# SERVICE COMMANDS
# =============================================================================
run: migrate-up
	@echo "🚀 Starting server with Air..."
	@air

.PHONY: new-migration migrate-up migrate-down migrate-status run