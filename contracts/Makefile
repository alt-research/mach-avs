# Makefile for Mach AVS contracts

.PHONY: storage-report

storage-report:
	@echo "📝 Generating storage layout report..."
	@./lib/storage-layout-reporter/report.sh --output=docs/storage --source="src" --exclude="src/interfaces src/error src/core/groth16"
	@echo "✅ Storage report generated in docs/storage directory"