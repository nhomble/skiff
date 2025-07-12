# Test Structure

This directory contains all test-related files organized by purpose:

## test-cases/
YAML test files used for integration testing and examples:
- **Simple cases**: `simple-before.yaml`, `simple-after.yaml`
- **Replica changes**: `replica-change-before.yaml`, `replica-change-after.yaml`
- **Complex nested**: `hpa-before.yaml`, `hpa-after.yaml` (HPA with behavior policies)
- **Mixed changes**: Various combinations of add/remove/update operations
- **Edge cases**: Empty files, identical files, no namespace, etc.

## test-policies/
Policy files and policy testing scripts:
- **policy.rego**: Basic policy for replica limits and config changes
- **complex-policy.rego**: Advanced policies for HPA and service changes
- **hpa-policy.rego**: HPA policy demonstrating percentage limits
- **test-policy.sh**: Integration test script for all policies
- **demo-format.sh**: Demonstration of format improvements
- **format-improvements.md**: Documentation of new format benefits

## performance/
Performance testing and benchmarking:
- **benchmark.sh**: Performance benchmarks to ensure CI/CD compatibility

## Usage

Run tests from the project root:
```bash
make test/policy          # Run policy integration tests
make benchmark           # Run performance benchmarks  
make demo               # Show Terraform format improvements
make example            # Run simple diff example
```