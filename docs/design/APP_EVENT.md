# Application Event

WASM applications can produce a new type of event: **AppEvent**.

- Not directed to a specific user — data is **not encrypted** by the executor
- Has an indexed `eventSubType` field: copied as-is in the emitted event, no seed override in the executor
- Supported in both `ProcessRequest` and `DepositFunds` operations

## Smart Contract Changes

### IProcessorEndpoint.sol

Add the new event:

```solidity
event AppEvent(
    uint64 indexed applicationId,
    bytes32 indexed requestId,
    string indexed eventSubType,
    bytes data
);
```

Note: `string indexed` is stored as `keccak256(value)` in log topics (same as `UserEvent`).

### ProcessorEndpoint.sol

In function `stateUpdate`, add new parameters:

```solidity
bytes[] calldata appEvents,
string[] calldata appEventSubTypes
```

Validation (same rules as `UserEvent`):
- `appEvents.length == appEventSubTypes.length` — otherwise `revert InvalidPayload()`
- When `errorCode != NO_ERROR`, `appEvents` must be empty — otherwise `revert InvalidPayload()`

Emit events in a loop (same pattern as `UserEvent`):

```solidity
i = 0;
while (i < appEventsLength) {
    emit AppEvent(applicationId, processedRequestId, appEventSubTypes[i], appEvents[i]);
    unchecked { ++i; }
}
```

### Structs.sol

Add `appEvents` and `appEventSubTypes` as separate fields in `SignatureParams` (same approach as `events` / `eventSubTypes`):

```solidity
struct SignatureParams {
    // ... existing fields ...
    bytes[] events;
    string[] eventSubTypes;
    bytes[] appEvents;          // new
    string[] appEventSubTypes;  // new
    // ... existing fields ...
}
```

### AbstractTeeAuthenticator.sol

Include `appEvents` and `appEventSubTypes` in the signed hash within `checkSignature`.

## Go Changes

### pkg/common/types.go

Add a new struct (distinct from `PlainEvent`, no `UserID` field since AppEvent is not user-directed):

```go
type PlainAppEvent struct {
    EventSubType string `json:"eventSubType"`
    Data         []byte `json:"data"`
}
```

### pkg/wasm/common/types.go

Add `AppEvents` field to both result structs:

```go
type ProcessResult struct {
    State       []byte              `json:"state"`
    Events      []common.PlainEvent `json:"events"`
    AppEvents   []common.PlainAppEvent `json:"appEvents"`
    Withdrawals []common.Withdrawal `json:"withdrawals"`
    Report      []byte              `json:"report,omitempty"`
    Fuel        *common.Big         `json:"fuel"`
    Error       string              `json:"error,omitempty"`
}

type DepositResult struct {
    State     []byte                 `json:"state"`
    Events    []common.PlainEvent    `json:"events"`
    AppEvents []common.PlainAppEvent `json:"appEvents"`
    Fuel      *common.Big            `json:"fuel"`
    Error     string                 `json:"error,omitempty"`
}
```

### pkg/executor/

- Propagate `AppEvents` from WASM result through the execution pipeline
- AppEvent data is **not encrypted** (unlike UserEvent)
- Include `appEvents` and `appEventSubTypes` in the signed payload for `stateUpdate`

### pkg/manager/

- Pass `appEvents` and `appEventSubTypes` to `stateUpdate` contract call

## WASM Guest Changes

### app/simple/

Update the simple app to demonstrate AppEvent emission by populating `AppEvents` in `ProcessResult` and `DepositResult`.

## Subgraph Changes

### subgraphs/hcce/schema.graphql

Add new entity (same structure as `UserEvent`, with `data` instead of `encryptedData`):

```graphql
type AppEvent @entity(immutable: true) {
    id: ID!
    applicationId: BigInt!
    requestId: Bytes!
    eventSubType: Bytes!
    data: Bytes!
    blockNumber: BigInt!
    logIndex: BigInt!
    sortKey: BigInt!
    blockTimestamp: BigInt!
}
```

No explicit relationships with other entities — linked via query on `applicationId` / `requestId` (same as `UserEvent`).

### subgraphs/hcce/subgraph.yaml

Add event mapping:

```yaml
- event: AppEvent(indexed uint64,indexed bytes32,indexed string,bytes)
  handler: handleAppEvent
```

### subgraphs/hcce/src/processor-endpoint.ts

Add handler following the same pattern as `handleUserEvent`:

```typescript
export function handleAppEvent(event: AppEventEvent): void {
    const id = event.transaction.hash.concatI32(event.logIndex.toI32()).toHex();
    let entity = new AppEvent(id);
    entity.applicationId = event.params.applicationId;
    entity.requestId = event.params.requestId;
    entity.eventSubType = event.params.eventSubType;
    entity.data = event.params.data;
    entity.blockNumber = event.block.number;
    entity.logIndex = event.logIndex;
    entity.sortKey = event.block.number.times(SORT_BASE).plus(event.logIndex);
    entity.blockTimestamp = event.block.timestamp;
    entity.save();
}
```

## Test Changes

### contracts/test/

Update Hardhat tests to cover:
- `stateUpdate` with AppEvents (happy path, single and multiple events)
- Validation: mismatched `appEvents` / `appEventSubTypes` array lengths → `InvalidPayload`
- Validation: non-empty `appEvents` with `errorCode != NO_ERROR` → `InvalidPayload`
- Signature verification includes `appEvents` and `appEventSubTypes`
- Empty `appEvents` arrays (backward-compatible calls)

### Go tests (pkg/)

- **pkg/executor/**: test that AppEvents from WASM result are propagated without encryption in the update payload
- **pkg/manager/**: test that `appEvents` and `appEventSubTypes` are correctly passed to `stateUpdate` contract call
- **pkg/wasm/**: test deserialization of `ProcessResult` and `DepositResult` with `AppEvents` field

### subgraphs/hcce/

Update subgraph tests to verify `AppEvent` entity creation and field mapping.

### tests/ (integration)

Update system integration tests to cover end-to-end AppEvent flow.

## Event Ordering

Multiple AppEvents within a single execution maintain order through the array index in `stateUpdate`. In the subgraph, ordering is determined by `logIndex` and `sortKey` (same as `UserEvent`).
