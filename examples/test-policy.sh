#!/bin/bash

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

# Build yspec
cd "$ROOT_DIR"
make build > /dev/null 2>&1

echo "Testing yspec policy integration..."
echo

echo "Test 1: Replica change (2 → 5 replicas)"
./yspec examples/test-cases/replica-change-before.yaml examples/test-cases/replica-change-after.yaml | conftest test --policy examples/policy.rego -

echo
echo "Test 2: Mixed changes"
./yspec examples/test-cases/mixed-changes-before.yaml examples/test-cases/mixed-changes-after.yaml | conftest test --policy examples/policy.rego -

echo
echo "Test 3: No changes"
./yspec examples/test-cases/identical-before.yaml examples/test-cases/identical-after.yaml | conftest test --policy examples/policy.rego -