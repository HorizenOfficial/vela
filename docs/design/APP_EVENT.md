# Application Event

WASM applications can produce a new type of event: **AppEvent**.

- Not directed to a specific user — data is **not encrypted** by the executor
- Has an indexed `eventSubType` field of type `bytes32`: stored as-is in log topics (no hashing), readable for short ASCII tags (≤32 bytes). Represented as `[32]byte` end-to-end (host, guest, signature hash, ABI) — no ASCII/hex re-encoding in between.
- Supported in both `ProcessRequest` and `DepositFunds` operations

## Smart Contract Changes

### IProcessorEndpoint.sol

New event (alongside existing `UserEvent`):

```solidity
event AppEvent(
    uint64 indexed applicationId,
    bytes32 indexed requestId,
    bytes32 indexed eventSubType,
    bytes data
);
```

Note: both `UserEvent` and `AppEvent` use `bytes32 indexed eventSubType` (not `string indexed`). This saves gas (no keccak256 on-chain) and preserves the value in log topics for short strings. If an app needs longer messages, it can put the hash in `eventSubType` and the full message in `data`.

### Structs.sol

New struct for grouping event data (used by both user events and app events):

```solidity
struct EventData {
    bytes[] events;
    bytes32[] subTypes;
}
```

`SignatureParams` uses `EventData` for both event types:

```solidity
struct SignatureParams {
    uint64 applicationId;
    bytes32 prevStateRoot;
    bytes32 newStateRoot;
    bytes32 processedRequestId;
    EventData userEvents;
    EventData appEvents;
    WithdrawalRequest[] withdrawalRequests;
    uint256 refundAmount;
    uint256 applicationFee;
    ErrorCode errorCode;
    string errorMsg;
}
```

### ProcessorEndpoint.sol

`stateUpdate` function signature uses `EventData` structs to avoid stack-too-deep:

```solidity
function stateUpdate(
    uint64 applicationId,
    bytes32 prevStateRoot,
    bytes32 newStateRoot,
    bytes32 processedRequestId,
    Structs.EventData calldata userEventData,
    Structs.EventData calldata appEventData,
    Structs.WithdrawalRequest[] calldata withdrawalRequests,
    uint256 refund,
    uint256 applicationFees,
    Structs.ErrorCode errorCode,
    string calldata errorMsg,
    bytes calldata signature
) external;
```

Validation (same rules for both user and app events):
- `events.length == subTypes.length` — otherwise `revert InvalidPayload()`
- When `errorCode != NO_ERROR`, both event arrays must be empty — otherwise `revert InvalidPayload()`

### AbstractTeeAuthenticator.sol

`checkSignature` hashes each `EventData` component separately:

```solidity
bytes32 userEventsHash = keccak256(abi.encode(params.userEvents.events));
bytes32 userEventSubTypesHash = keccak256(abi.encode(params.userEvents.subTypes));
bytes32 appEventsHash = keccak256(abi.encode(params.appEvents.events));
bytes32 appEventSubTypesHash = keccak256(abi.encode(params.appEvents.subTypes));
```

Message hash includes all four hashes between `processedRequestId` and `withdrawalRequestsHash`.

## Go Changes

### pkg/common/types.go

`AppEvent` struct (no `UserID` since it's not user-directed). `EventSubType` is
stored as `[32]byte` on both host and guest sides, matching the on-chain
`bytes32` type exactly — no ASCII-to-bytes32 truncation at the serialization
boundary.

```go
type AppEvent struct {
    EventSubType [32]byte `json:"eventSubType"`
    Data         []byte   `json:"data"`
}
```

`Event` and `PlainEvent` use the same `[32]byte` representation for
`EventSubType`.

The framework provides no built-in string→bytes32 helper. How subtype bytes
are produced is a per-application policy: short readable tags packed with
`copy(b[:], s)` (as SimpleApp does via its own `app.SubTypeFromString`),
keccak256 hashes, enumerated constants, or anti-linkability subtypes derived
from a user seed (see `pkg/executor/subtype.go`) are all valid — apps just
need to be consistent between producer and filter.

`UpdatePayload` includes `AppEvents []AppEvent`.

### pkg/wasm/common/types.go

`AppEvents` field added to `DepositResult` and `ProcessResult`:

```go
type DepositResult struct {
    State     []byte                 `json:"state"`
    Events    []common.PlainEvent    `json:"events"`
    AppEvents []common.AppEvent      `json:"appEvents"`
    Fuel      *common.Big            `json:"fuel"`
    Error     string                 `json:"error,omitempty"`
}
```

### pkg/executor/

- `Runtime` interface: `Deposit` and `ProcessRequest` return `[]common.AppEvent` as additional return value
- AppEvent data is **not encrypted** (passed through as-is, unlike UserEvent)
- `MsgToSignBuilder`: uses `EventSubType` directly as `[32]byte` for ABI encoding as `bytes32[]` — no string/bytes32 conversion at the hashing boundary
- Inline comments on `msgArgs` and `values` slices for ordering verification against Solidity

### pkg/blockchain/client.go

`SubmitStateUpdate` constructs `StructsEventData` for both user events and app events, passing `EventSubType` ([32]byte) straight into the `bytes32[]` ABI field.

### pkg/executor/msgtosign_builder.go

`ethSignStateUpdate` uses `bytes32[]` ABI type for event subtypes (not `string[]`). `appEvents` and `appEventSubTypes` are required parameters (not optional with defaults) to prevent silent signing mismatches.

### pkg/executor/subtype.go

`GenerateSubtype` / `AllSubtypes` / `GenerateRandomSubtype` return the raw 32
bytes of the HMAC-SHA256 digest (`[32]byte`). They used to return `"0x" + hex`
strings, which — when copied into a `[32]byte` on the blockchain boundary —
were silently truncated to 32 ASCII characters, losing ~17 bytes of entropy
and writing the ASCII hex representation on-chain instead of the digest.

## vela-common-go Changes

### wasm/types/common.go

Guest-side structs use `[32]byte` for `EventSubType`, mirroring the host side
(they serialize to the same JSON as `[32]byte` on both sides, so the
host/guest JSON round-trip stays symmetric):

```go
type PlainEvent struct {
    UserID       Address  `json:"userId"`
    EventSubType [32]byte `json:"eventSubType"`
    Data         []byte   `json:"data"`
}

type AppEvent struct {
    EventSubType [32]byte `json:"eventSubType"`
    Data         []byte   `json:"data"`
}
```

The guest framework package also provides no subtype-encoding helper; apps
own that policy. SimpleApp defines its own `app.SubTypeFromString` (in
`app/simple/app/subtype.go`) that packs up to 32 ASCII bytes left-aligned.
Other apps are free to pick a different encoding (e.g. keccak256).

### wasm/types/results.go

`AppEvents []AppEvent` added to `DepositResult` and `ProcessResult`.

## vela-common-ts Changes

- **Typechain types**: regenerated with `bytes32 indexed eventSubType`, `bytes32[] subTypes`, `EventData` struct, and `AppEvent` event
- **`blockchain/client.ts`**: added `getAppEvents(fromBlock, toBlock, applicationId, requestId, eventSubType)` method. `eventSubType` filter must be a bytes32 hex value (`ethers.encodeBytes32String()` to convert)
- **`subgraph/types.ts`**: added `AppEvent` interface and `getAppEvents` method to `SubgraphClient`
- **`subgraph/client.ts`**: implemented `getAppEvents` GraphQL query (same pattern as `getUserEvents`, with `data` instead of `encryptedData`)
- **`subgraph/mock.ts`**: added `withAppEvents` and `getAppEvents` to mock client

## Subgraph Changes

### subgraphs/hcce/schema.graphql

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

### subgraphs/hcce/subgraph.yaml

```yaml
- event: AppEvent(indexed uint64,indexed bytes32,indexed bytes32,bytes)
  handler: handleAppEvent
```

### subgraphs/hcce/src/processor-endpoint.ts

Handler follows the same pattern as `handleUserEvent`.

## Test Changes

### Hardhat tests (234 passing)

- All `stateUpdate` calls updated to use `EventData` structs
- Event subtype assertions use `ethers.encodeBytes32String()` (not `ethers.id()`)
- `ethSignStateUpdate`: `appEvents`/`appEventSubTypes` are required params (before `withdrawalRequests`)
- New tests: appEventData length mismatch, AppEvent emission, appEventData non-empty on error path
- `checkSignature.spec.ts`: tests with non-empty appEvents in payload

### Go unit tests

- `MockRuntime.ProcessRequest` emits an AppEvent on transfer
- `TestHandleProcessRequest_Process_Success`: verifies AppEvents propagate through executor pipeline (2 UserEvents + 1 AppEvent for transfer)
- `TestCheckSignature`: verifies Go-side signature matches Solidity contract (covers AppEvents in hash)

### System integration tests

- `TestSimpleAppDepositEmitsAppEvent`: end-to-end test verifying AppEvent presence in UpdatePayload after deposit

## Event Ordering

Multiple AppEvents within a single execution maintain order through the array index in `stateUpdate`. In the subgraph, ordering is determined by `logIndex` and `sortKey` (same as `UserEvent`).
