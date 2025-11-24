.PHONY: help deps synth deploy plan destroy diff list outputs clean watch version

help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

deps: ## Download Go dependencies
	go mod download
	go mod tidy

synth: deps ## Generate Terraform using cdktf CLI
	cdktf synth

list: ## List all stacks
	cdktf list

diff: ## Show infrastructure changes
	cdktf diff

plan: synth ## Show Terraform plan (manual check)
	cd cdktf.out/stacks/my-app-dev-stack && terraform init && terraform plan

deploy: ## Deploy infrastructure to AWS
	@echo "⚠️  WARNING: This will create real AWS resources!"
	@echo "Make sure AWS credentials are configured."
	@read -p "Continue? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	cdktf deploy

destroy: ## Destroy all infrastructure
	@echo "⚠️  WARNING: This will destroy all infrastructure!"
	@read -p "Continue? (y/N): " confirm && [ "$$confirm" = "y" ] || exit 1
	cdktf destroy

clean: ## Clean generated files
	rm -rf cdktf.out
