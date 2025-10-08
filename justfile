# Go environment variables
export GO111MODULE := "on"
export GOPROXY := "https://proxy.golang.org,direct"

# Variables
DOCKER := env_var_or_default("DOCKER", "docker")

# Default recipe
default: ci

# Run the binary
run *args:
    go run main.go {{args}}

# Setup git hooks
dev:
    cp -f scripts/pre-commit.sh .git/hooks/pre-commit

# Install dependencies
setup:
    go mod tidy

# Build the binary
build:
    go build

# Run tests
test TEST_OPTIONS="" SOURCE_FILES="./..." TEST_PATTERN=".":
    #!/usr/bin/env bash
    export LC_ALL=C
    go test {{TEST_OPTIONS}} -failfast -race -coverpkg=./... -covermode=atomic -coverprofile=coverage.txt {{SOURCE_FILES}} -run {{TEST_PATTERN}} -timeout=5m

# Open the cover tool
cover:
    go tool cover -html=coverage.txt

# Format all code with gofumpt
fmt:
    gofumpt -w -l .

# Lint the code with golangci-lint
lint:
    golangci-lint run ./...

# Run all CI steps
ci: setup build test

# Build a snapshot release locally for testing
snapshot:
    goreleaser release --snapshot --clean --skip=publish

# Run GoReleaser in release mode (triggered by CI on tags)
goreleaser:
    goreleaser release --clean

# Create a new tag and push (triggers release workflow)
release:
    #!/usr/bin/env bash
    NEXT=$(svu n)
    git tag $NEXT
    echo $NEXT
    git push origin --tags