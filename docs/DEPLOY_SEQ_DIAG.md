# Deploy Sequence Diagram (v2 - Multi-App)

Application deployment flow with system-assigned application IDs, off-chain artifact upload, and on-chain deploy descriptor.

## End-to-End Flow

```
┌──────────┐     ┌───────────────┐     ┌──────────────┐     ┌─────────┐     ┌──────────┐       ┌─────────────┐
│  Wallet  │     │   Authority   │     │    Smart     │     │ Manager │     │ Executor │       │    WASM     │
│   CLI    │     │   Service     │     │  Contract    │     │         │     │ (Enclave)│       │   Runtime   │
└────┬─────┘     └──────┬────────┘     └──────┬───────┘     └────┬────┘     └────┬─────┘       └──────┬──────┘
     │                  │                     │                  │               │                    │
     │  1. POST /deploy/upload (WASM bytes)   │                  │               │                    │
     │─────────────────>│                     │                  │               │                    │
     │                  │                     │                  │               │                    │
     │                  │  SaveWASM():        │                  │               │                    │
     │                  │  write tmp file,    │                  │               │                    │
     │                  │  compute sha256,    │                  │               │                    │
     │                  │  rename to          │                  │               │                    │
     │                  │  blobs/<sha>.wasm   │                  │               │                    │
     │                  │                     │                  │               │                    │
     │  {artifactId, wasmSha256}              │                  │               │                    │
     │<─────────────────│                     │                  │               │                    │
     │                  │                     │                  │               │                    │
     │  2. Verify local SHA256 matches        │                  │               │                    │
     │──┐               │                     │                  │               │                    │
     │<─┘               │                     │                  │               │                    │
     │                  │                     │                  │               │                    │
     │  3. Build JSON deploy descriptor       │                  │               │                    │
     │──┐  {mode: "artifact_ref",             │                  │               │                    │
     │<─┘   artifactId, wasmSha256}           │                  │               │                    │
     │                  │                     │                  │               │                    │
     │  4. submitDeployRequest(protocolVersion, descriptor)      │               │                    │
     │───────────────────────────────────────>│                  │               │                    │
     │                  │                     │                  │               │                    │
     │                  │          5. Contract checks:           │               │                    │
     │                  │             hasRole(DEPLOYER_ROLE)?    │               │                    │
     │                  │             availableDeploySlots > 0?  │               │                    │
     │                  │             queue not full?            │               │                    │
     │                  │             fee >= minFeePerRequest?   │               │                    │
     │                  │                     │                  │               │                    │
     │                  │             Derives appID =            │               │                    │
     │                  │             uint64(bytes8(requestId))  │               │                    │
     │                  │             Enqueues deploy request    │               │                    │
     │                  │             Emits DeployRequestSubmitted               │                    │
     │                  │                     │                  │               │                    │
     │  (appID, requestID, blockNumber)       │                  │               │                    │
     │<───────────────────────────────────────│                  │               │                    │
     │                  │                     │                  │               │                    │
     │  6. Print "Assigned ApplicationID: N"  │                  │               │                    │
     │──┐               │                     │                  │               │                    │
     │<─┘               │                     │                  │               │                    │
     │                  │                     │                  │               │                    │
     .                  .                     .  (polling)       .               .                    .
     │                  │                     │                  │               │                    │
     │                  │                     │  7. GetNextPendingRequest()      │                    │
     │                  │                     │<─────────────────│               │                    │
     │                  │                     │                  │               │                    │
     │                  │                     │  (deployReq, stateRoot)          │                    │
     │                  │                     │─────────────────>│               │                    │
     │                  │                     │                  │               │                    │
     │                  │         8. Manager processDeployApp(req):              │                    │
     │                  │            Parse deploy descriptor     │               │                    │
     │                  │            Extract SHA from artifactId │               │                    │
     │                  │            Read WASM from              │               │                    │
     │                  │            ARTIFACTS_PATH/blobs/<sha>.wasm             │                    │
     │                  │                     │                  │               │                    │
     │                  │                     │                  │  9. SendDeployApp(req, state, wasm)│
     │                  │                     │                  │──────────────>│                    │
     │                  │                     │                  │               │                    │
     │                  │                     │                  │               │  10. Validate:     │
     │                  │                     │                  │               │      - not already │
     │                  │                     │                  │               │        deployed    │
     │                  │                     │                  │               │      - descriptor  │
     │                  │                     │                  │               │        JSON valid  │
     │                  │                     │                  │               │      - WASM not    │
     │                  │                     │                  │               │        empty       │
     │                  │                     │                  │               │      - SHA256      │
     │                  │                     │                  │               │        matches     │
     │                  │                     │                  │               │                    │
     │                  │                     │                  │               │  11. Compile WASM, │
     │                  │                     │                  │               │      call          │
     │                  │                     │                  │               │      load_module   │
     │                  │                     │                  │               │      (appId)       │
     │                  │                     │                  │               │──────────────────> │
     │                  │                     │                  │               │                    │
     │                  │                     │                  │               │  (initial state,   │
     │                  │                     │                  │               │   fuel)            │
     │                  │                     │                  │               │<───────────────────│
     │                  │                     │                  │               │                    │
     │                  │                     │                  │  (signed updatePayload + appState) │
     │                  │                     │                  │<──────────────│                    │
     │                  │                     │                  │               │                    │
     │                  │                     │  12. Store state + WASM in LevelDB                    │
     │                  │                     │                  │──┐            │                    │
     │                  │                     │                  │<─┘            │                    │
     │                  │                     │                  │               │                    │
     │                  │                     │  13. SubmitStateUpdate(appID, prevRoot, newRoot, ...) │
     │                  │                     │<─────────────────│               │                    │
     │                  │                     │                  │               │                    │
     │                  │       14. Contract sets                │               │                    │
     │                  │           applicationStateRoots[appID] │               │                    │
     │                  │           = newRoot                    │               │                    │
     │                  │           (implicit app registration)  │               │                    │
     │                  │           Emits DeployRequestCompleted │               │                    │
     │                  │                     │                  │               │                    │
     │  15. Poll subgraph: GetDeployRequestCompletedByID(requestID)              │                    │
     │───────────────────────────────────────>│                  │               │                    │
     │                  │                     │                  │               │                    │
     │  (status: OK)    │                     │                  │               │                    │
     │<───────────────────────────────────────│                  │               │                    │
     │                  │                     │                  │               │                    │
     │  16. SaveApplicationID(wallet.conf, appID)                │               │                    │
     │──┐               │                     │                  │               │                    │
     │<─┘               │                     │                  │               │                    │
     │                  │                     │                  │               │                    │
```

## Key Design Points

- **System-assigned application ID:** The caller does not choose the application ID. The contract derives it deterministically from the first 8 bytes of the request hash (`uint64(bytes8(requestId))`). This prevents ID collisions and removes the need for a separate registration step.
- **Implicit registration:** The application is considered registered when `stateUpdate` sets `applicationStateRoots[appID]` to a non-zero value (step 14). There is no separate registration transaction.
- **Separate contract function:** Deploy requests use `submitDeployRequest()`, not `submitRequest()`. The latter rejects `DEPLOYAPP` request types. This enforces `DEPLOYER_ROLE` authorization and the `availableDeploySlots` limit.
- **Off-chain artifact, on-chain descriptor:** The WASM binary is uploaded to the Authority Service and stored content-addressed as `blobs/<sha256>.wasm`. Only a lightweight JSON descriptor (referencing the artifact by hash) goes on-chain. The manager resolves the artifact from the shared filesystem.
- **Wallet persistence:** After successful deploy, the wallet automatically saves the assigned application ID to `wallet.conf`, commenting out any previous value.

## Error Paths

### At step 4 (on-chain, smart contract)

| Condition | Result |
|---|---|
| Sender lacks `DEPLOYER_ROLE` | Revert `DeployerNotAllowed` (tx fails, request never queued) |
| `availableDeploySlots == 0` | Revert `MaxNumOfApplicationsExceeded` |
| Queue full | Revert `QueueThresholdExceeded` |
| Fee below minimum | Revert `FeeValueBelowMinimum` |

### At step 8 (manager resolveDeployWASM)

Manager artifact resolution failures do NOT mark the request failed.
The manager logs a warning and forwards `wasmModule = nil` to the executor,
which then produces the appropriate signed error payload.

| Condition | Internal Error | Behavior |
|---|---|---|
| Empty/non-JSON payload | `ARTIFACT_DESCRIPTOR_INVALID` | Forward nil to executor |
| No `artifactId` or `wasmSha256` extractable | `ARTIFACT_DESCRIPTOR_INVALID` | Forward nil to executor |
| Artifact file not found | `ARTIFACT_NOT_FOUND` | Forward nil to executor |
| I/O read error | `ARTIFACT_LOAD_FAILED` | Forward nil to executor |

### At step 9-11 (executor HandleDeployApp) -- signed error payloads

All validation failures produce a signed error payload that the manager submits on-chain.

| Condition | Error Code | Error Message |
|---|---|---|
| Wrong protocol version | *(hard error, not payload)* | "protocol version N is not admitted" |
| Fee below minimum | *(hard error, not payload)* | "request fee is below minimum fee" |
| Wrong request type | `CodeInternalFallback` | "failed to deploy application" |
| Application already deployed | `CodeApplicationAlreadyDeployed` | "application N was already deployed" |
| Invalid descriptor JSON | `CodeInternalFallback` | "failed to deploy application" |
| WASM module nil/empty | `CodeFailedLoadingOrGettingModule` | "failed to load or get module" |
| WASM hash mismatch | `CodeFailedLoadingOrGettingModule` | "failed to load or get module" |
| LoadModule fails | `CodeFailedLoadingOrGettingModule` | "failed to load or get module" |
| Insufficient fuel | `CodeInsufficientFuel` | "insufficient fuel: required X wei, provided Y wei" |
