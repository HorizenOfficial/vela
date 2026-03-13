Review the error handling in the diff of this PR (against the `dev` branch) to verify that the correct error type is chosen when a request execution fails.

## Error Handling Rules

In the executor (`pkg/executor/`), there are exactly two ways to handle a failure:

### 1. Execution Error (signed, submitted on-chain)

```go
errorPayload, err := e.processErrorResponse(req, stateRoot, apperrors.New(apperrors.CodeXxx, "message"))
return errorPayload, nil, nil, err
```

Use this **only** when the failure is something the smart contract could NOT have validated — a genuine execution-time failure that the user should be notified about on-chain (and refunded for).

Examples of correct execution errors:
- WASM runtime failures: module load, memory write, function not found, fingerprint mismatch
- Application logic failures: deposit failed, request func failed
- App already deployed / app state not found
- Insufficient fuel (after computing actual fuel cost)
- Key parsing errors, pub key not registered
- Report validation failures (wrong report for request type)
- App data serialization failures

### 2. Transient Error (plain Go error, retry on next poll)

```go
return nil, nil, nil, fmt.Errorf("some error: %w", err)
```

Use this when:
- **Infrastructure failures**: cryptographic operations (encryption, signing, decryption), DB errors, communication failures, state integrity mismatches
- **Tampered or pre-validated fields**: If a request field was already validated on-chain and has an unexpected value at the executor, it means the request was tampered with between the chain and the executor. The executor should NOT execute the request and should NOT submit an error on-chain. Examples:
  - `validateRequest()`: protocol version, applicationId, minimum fee — all checked on-chain
  - Request type mismatches (e.g., non-deploy request reaching `HandleDeployApp`)
  - `AssociateKey` payload length (checked on-chain)
  - State decryption or state root mismatch (`fromEncryptedStateToAppData`)

## The Most Common Mistake

**Returning an execution error (`processErrorResponse` / `apperrors.New`) when it should be a transient error (`fmt.Errorf`).** This happens when a developer adds a new validation check and treats it as a request failure, when in reality the field was already validated on-chain — making the invalid value evidence of tampering, not a user error.

## How to Review

For every new or modified error return path in the diff:

1. **Identify the error type used**: Is it `processErrorResponse`/`apperrors.New` (execution error) or `fmt.Errorf`/plain error (transient)?
2. **Ask**: Could the smart contract have already rejected this condition? If yes, an unexpected value here means tampering → must be transient.
3. **Ask**: Is this a genuine runtime failure that only the executor can detect? If yes → execution error is correct.
4. **Ask**: Is this an infrastructure/system failure (crypto, DB, network)? If yes → must be transient.

## Where to Look

Focus on Go code in:
- `pkg/executor/` — where execution errors vs transient errors are decided
- `pkg/manager/` — where the two error paths diverge (`updatePayload.ErrorCode != 0` vs `err != nil`)
- `pkg/wasm/` — WASM runtime errors that bubble up as `*apperrors.RequestFailure`
- `pkg/common/apperrors/` — new error codes or categories being added

Also check the smart contracts in `contracts/contracts/` (especially `ProcessorEndpoint.sol`) to verify which fields are validated on-chain. If the contract already rejects a condition, the executor seeing that condition means tampering — and the error must be transient, not an execution error.

## Output Format

For each finding, report:
- File and line number
- The error return in question
- Whether it should be execution or transient, and why
- Suggested fix (if the current choice is wrong)

Only report actionable findings. Skip correct error handling — do not list things that are fine.
