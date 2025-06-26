# Testing the Horizen Privacy Preserving Execution System

This document describes the testing strategy for the Horizen Privacy Preserving Execution System.

## Test Structure

The tests are organized into the following categories:

1. **Unit Tests**: Tests for individual components and functions
   - Located in the directory of each package, e.g., [`pkg/manager/manager_test.go`](pkg/manager/manager_test.go)
2. **System\Integration\E2E Tests**: Tests for the integration of multiple components or system as a whole
   - Located in [`/tests`](tests/README.md)


## Running the Tests

To run the tests, use the following command:

```bash
cd code
go test ./...
```

This will run all the tests in the project. Some tests are skipped because they require external dependencies like Wasmtime-go.
E2E tests might end up being skipped as well, depending on the environment setup.

## Test Coverage

To run the tests with coverage, use the following command:

```bash
cd code
go test ./... -cover
```

This will run all the tests and report the test coverage.