# JSON to Terraform with CDKTF

Go application that converts JSON configurations into Terraform using CDKTF.

## What This Does

```
config.json (JSON) → main.go (Go + CDKTF) → cdktf CLI → Terraform files
```

Developers write simple JSON. Go code reads JSON and generates Terraform using the CDKTF library. The cdktf CLI orchestrates the process.

## Quick Start

```bash
flox activate
make deps
make synth
cat cdktf.out/stacks/my-app-dev-stack/cdk.tf.json
```

## Input Format

**File: `config.json`**
```json
{
  "project": "my-app",
  "environment": "dev",
  "region": "us-west-2",
  "storage": {
    "bucket_name": "my-app-data",
    "enable_versioning": true
  }
}
```

## Commands

```bash
make deps      # Install Go dependencies
make synth     # Generate Terraform
make list      # List stacks
make diff      # Show changes
make deploy    # Deploy to AWS
make destroy   # Destroy infrastructure
make clean     # Remove generated files
```

## How It Works

1. `cdktf synth` reads `cdktf.json` which specifies `"app": "go run ."`
2. cdktf CLI spawns Go subprocess
3. Go reads `config.json`, creates CDKTF constructs, calls `app.Synth()`
4. CDKTF library outputs Terraform JSON to stdout
5. cdktf CLI captures stdout and writes to `cdktf.out/stacks/*/cdk.tf.json`
6. Terraform commands (init/apply/destroy) run against generated files

## Project Structure

```
tf-cdk/
├── config.json          # Developer input
├── main.go              # Go app (reads JSON, generates Terraform)
├── go.mod/go.sum        # Dependencies
├── cdktf.json           # cdktf CLI configuration
├── Makefile             # Commands
├── README.md            # This file
├── CLAUDE.md            # AI assistant guide
└── cdktf.out/           # Generated Terraform (after synth)
```
