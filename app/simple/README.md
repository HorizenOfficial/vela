# Simple App

## Purpose

The Simple App is a sample WebAssembly (WASM) module that demonstrates a basic application. It manages user accounts with simple deposit and withdrawal functionalities. It also includes a feature to compare the balances of two users. This app serves as an example of how to build and interact with a WASM-based application on the platform.

## Testing

The Simple App is tested at three different levels: unit, integration, and system tests.

### Unit Tests (`app/app_test.go`)

Unit tests focus on the core application logic within the `app` package. They test individual functions in isolation to ensure they behave as expected.

**Key aspects tested:**

-   **`LoadModule`:** Verifies that the module initializes with a correct default state.
-   **`DepositFunds`:** Checks that deposits are correctly added to new and existing accounts.
-   **`ProcessRequest`:**
    -   Tests successful and unsuccessful withdrawals (e.g., insufficient balance).
    -   Verifies the balance comparison logic for all cases (richer, poorer, equal).
    -   Ensures correct handling of invalid instructions and payloads.
-   **`GenerateDeanonymizationReport`:** Confirms that the deanonymization report is generated correctly.

These tests do not require the WASM runtime and are executed using the standard `go test` command.

### Integration Tests (`integration_test.go`)

Integration tests verify the interaction between the Go application logic and the WASM runtime. They ensure that the compiled WASM module can be loaded, executed, and that data is correctly passed between the host and the module.

**Test flow:**

1.  **Build the WASM module:** The test starts by compiling the Go code into a `.wasm` file using the `make build` command.
2.  **Load the module:** It uses a `Wasmtime` runtime to load the compiled WASM module.
3.  **Simulate user interactions:**
    -   It simulates deposits from two different users.
    -   It processes a withdrawal request.
    -   It processes a balance comparison request.
    -   It generates a deanonymization report.
4.  **Assert outcomes:** After each step, the test asserts that the application state is updated correctly and that the events and withdrawals are generated as expected.

These tests provide confidence that the application works as a whole, including the WASM compilation and execution steps.

### System Tests (`tests/system/simple_app_system_test.go`)

System tests cover the end-to-end flow of deploying and interacting with the Simple App in a simulated production environment. These tests involve all the components of the system, including the executor, manager, and the blockchain.

**Key scenarios tested:**

-   **`TestDeploySimpleApp`:**
    1.  Builds the WASM module.
    2.  Starts the executor and manager.
    3.  Adds user keys to the key registry.
    4.  Submits a `deploy` request with the WASM bytecode.
    5.  Waits for the application state to be created in both the local database and the blockchain.
    6.  Asserts that the deployment request is marked as completed.

-   **`TestWasmtimeRuntimeSimpleAppFullSystemFlow`:**
    -   This test executes a full end-to-end user flow, including deploying the app, depositing funds, and processing requests. It uses a test utility (`ExecTestAppFullSystemFlow`) to run through a standard sequence of interactions.