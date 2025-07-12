#!/bin/bash

set -e  # Exit on any error

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"

# Get the directory where the script is located
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "$(dirname "$SCRIPT_DIR")")"

# Ensure skiff is in the current directory or PATH
if ! command -v ./skiff &> /dev/null && ! command -v skiff &> /dev/null; then
    echo "skiff executable not found. Please ensure it is in the current directory or your PATH."
    exit 1
fi

# Ensure conftest is in the PATH
if ! command -v conftest &> /dev/null; then
    echo "conftest executable not found. Please ensure it is in your PATH."
    echo "You can download it from https://www.conftest.dev/install/"
    exit 1
fi

echo "Testing skiff policy integration..."
echo

# Test 1: Should fail (replica violations)
echo "Test 1: Replica change (2 → 5 replicas) - expect failures"
EXIT_CODE=0
./skiff test/test-cases/replica-change-before.yaml test/test-cases/replica-change-after.yaml | conftest test --policy test/test-policies/policy.rego - || EXIT_CODE=$?
if [ $EXIT_CODE -eq 1 ]; then
    echo "✅ Policy correctly caught violations"
else
    echo "❌ Expected policy failures but got exit code: $EXIT_CODE"
fi

echo
# Test 2: Should fail (mixed changes with violations)
echo "Test 2: Mixed changes - expect failures"
EXIT_CODE=0
./skiff test/test-cases/mixed-changes-before.yaml test/test-cases/mixed-changes-after.yaml | conftest test --policy test/test-policies/policy.rego - || EXIT_CODE=$?
if [ $EXIT_CODE -eq 1 ]; then
    echo "✅ Policy correctly caught violations"
else
    echo "❌ Expected policy failures but got exit code: $EXIT_CODE"
fi

echo
# Test 3: Should pass (no changes)
echo "Test 3: No changes - expect pass"
EXIT_CODE=0
./skiff test/test-cases/identical-before.yaml test/test-cases/identical-after.yaml | conftest test --policy test/test-policies/policy.rego - || EXIT_CODE=$?
if [ $EXIT_CODE -eq 0 ]; then
    echo "✅ Policy correctly passed clean diff"
else
    echo "❌ Expected clean pass but got exit code: $EXIT_CODE"
fi

echo
# Test 4: Should pass (allowed change - image update, 1 replica)
echo "Test 4: Allowed change (image update, 1 replica) - expect pass"
EXIT_CODE=0
./skiff test/test-cases/allowed-change-before.yaml test/test-cases/allowed-change-after.yaml | conftest test --policy test/test-policies/policy.rego - || EXIT_CODE=$?
if [ $EXIT_CODE -eq 0 ]; then
    echo "✅ Policy correctly allowed safe change"
else
    echo "❌ Expected policy to allow safe change but got exit code: $EXIT_CODE"
fi

echo
# Test 5: Should warn (ConfigMap change)
echo "Test 5: ConfigMap change - expect warnings"
EXIT_CODE=0
./skiff test/test-cases/configmap-change-before.yaml test/test-cases/configmap-change-after.yaml | conftest test --policy test/test-policies/policy.rego - || EXIT_CODE=$?
if [ $EXIT_CODE -eq 0 ]; then
    echo "✅ Policy correctly warned about ConfigMap change"
else
    echo "❌ Expected policy to warn about ConfigMap change but got exit code: $EXIT_CODE"
fi

echo
# Test 6: Complex nested changes (HPA)
echo "Test 6: Complex HPA changes - expect failure and warning"
EXIT_CODE=0
./skiff test/test-cases/hpa-before.yaml test/test-cases/hpa-after.yaml | conftest test --policy test/test-policies/complex-policy.rego - || EXIT_CODE=$?
if [ $EXIT_CODE -eq 1 ]; then
    echo "✅ Policy correctly caught HPA violations and warnings"
else
    echo "❌ Expected HPA policy violations but got exit code: $EXIT_CODE"
fi

echo
# Test 7: Service changes - expect warning
echo "Test 7: Service complex changes - expect warning"
EXIT_CODE=0
./skiff test/test-cases/complex-nested-before.yaml test/test-cases/complex-nested-after.yaml | conftest test --policy test/test-policies/complex-policy.rego - || EXIT_CODE=$?
if [ $EXIT_CODE -eq 0 ]; then
    echo "✅ Policy correctly warned about Service changes"
else
    echo "❌ Expected Service policy warning but got exit code: $EXIT_CODE"
fi

echo
echo "🎯 Policy integration tests completed"