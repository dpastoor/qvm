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

# Create a new tag
release:
    #!/usr/bin/env bash
    NEXT=$(svu n)
    git tag $NEXT
    echo $NEXT
    git push origin --tags

# Test a package (internal helper)
_test-pkg Platform Image Cmd:
    docker run --platform linux/{{Platform}} --rm --workdir /tmp -v $PWD/dist:/tmp {{Image}} sh -c '{{Cmd}} && qvm --version'

# Test rpm packages
test-rpm:
    just _test-pkg amd64 fedora "rpm --nodeps -ivh qvm-*.x86_64.rpm"

# Test deb packages
test-deb:
    just _test-pkg amd64 ubuntu "dpkg -i qvm*_amd64.deb"

# Test apk packages
test-apk:
    just _test-pkg amd64 alpine "apk add --allow-untrusted -U qvm*_x86_64.apk"

# Test all built linux packages
test-packages: test-apk test-deb test-rpm

# Run GoReleaser either in snapshot or release mode
goreleaser: build
    #!/usr/bin/env bash
    SNAPSHOT=""
    if [[ $GITHUB_REF != refs/tags/v* ]]; then
        SNAPSHOT="--snapshot"
    fi
    goreleaser release --rm-dist $SNAPSHOT