#!/bin/bash

# Quick benchmark script to ensure yspec is fast enough for CI/CD

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"

cd "$ROOT_DIR"

echo "skiff performance benchmark"
echo

# Test with various file sizes
echo "Test 1: Simple files (small)"
time ./skiff test/test-cases/simple-before.yaml test/test-cases/simple-after.yaml > /dev/null

echo
echo "Test 2: Mixed changes (medium)"
time ./skiff test/test-cases/mixed-changes-before.yaml test/test-cases/mixed-changes-after.yaml > /dev/null

echo
echo "Test 3: Complete pipeline (skiff + policy)"
time (./skiff test/test-cases/replica-change-before.yaml test/test-cases/replica-change-after.yaml | conftest test --policy test/test-policies/policy.rego - > /dev/null)

echo
echo "✅ Benchmark complete - skiff should be fast enough for CI/CD"