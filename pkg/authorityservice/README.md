# Signature scheme for `/getreport`

This section specifies exactly how the client must build and sign the message used in the `POST /getreport` request, and how the signature must be encoded in JSON.

**TL;DR:**  
The client signs `keccak256(chain_id || app_id || report_id || nonce)` as **raw bytes** (no EIP-191 / personal_sign prefix), using the authority’s Ethereum private key.  
The signature is sent as a **65-byte R||S||V hex string**.

---

## Service configuration (blockchain connectivity)

The authority service now verifies on-chain completion before serving a report. Provide these settings via environment variables or `authorityservice.conf`:

- `CHAIN_ID`: expected chain id.
- `CHAIN_RPC_PROTOCOL`, `CHAIN_RPC_ADDRESS`, `CHAIN_RPC_PORT`: RPC endpoint to the node.
- `CHAIN_PROCESSOR_ADDRESS`: ProcessorEndpoint contract address.
- `AUTHORITY_SERVICE_SUBGRAPH_URL`: subgraph endpoint to read `RequestCompleted`/`UserEvent`.
- `MANAGER_REPORTS_FOLDER`: path to the report files shared with the manager (default `/tmp/horizen-pes-data/manager_reports`).
- `DEPLOY_ARTIFACTS_PATH`: path where `POST /deploy/upload` stores uploaded WASM artifacts (default `/tmp/horizen-pes-data/deploy_artifacts`).
- `DEPLOY_ARTIFACTS_MAX_SIZE_MB`: upload limit in MB (default `50`, set `0` only to explicitly disable limits).

Logging/TLS options are documented in `pkg/authorityservice/config.go`.

---

## Deploy upload endpoint (`POST /deploy/upload`)

- Request: `multipart/form-data` with required part `wasm`.
- Success response: `{ artifactId, wasmSha256 }`.
- Error responses:
  - `400` invalid multipart or missing wasm part.
  - `413` upload exceeds `DEPLOY_ARTIFACTS_MAX_SIZE_MB`.
  - `500` internal storage/persistence error.

---

## 1. Data to be signed

For a `POST /getreport` request, the following fields participate in the signature:

- `chain_id` (uint64)
- `app_id` (uint64)
- `report_id` (hex string, 32 bytes)
- `nonce` (hex string, 32 bytes)

The `nonce` is obtained from the previous `GET /nonce` call; it is **not** re-computed by the client.

The `salt` and `timestamp` are used by the server to validate the nonce, but they **are not included** in the signed message.

---

## 2. Message format (binary layout)

The message to sign is the concatenation of the following fields, as **raw bytes**:

1. `chain_id` – 8 bytes, unsigned, big-endian  
2. `app_id` – 8 bytes, unsigned, big-endian  
3. `report_id` – 32 bytes  
4. `nonce` – 32 bytes  

Pseudo-code:
```bash
message = uint64_be(chain_id)
|| uint64_be(app_id)
|| report_id_bytes // 32 bytes
|| nonce_bytes // 32 bytes
```

Where:

- `report_id_bytes` is the 32-byte value corresponding to the hex string `report_id`.
- `nonce_bytes` is the 32-byte value returned by `GET /nonce`.

---

## 3. Hashing

Hash the message using **Keccak-256** (Ethereum standard hash):

```bash
digest = keccak256(message)
```

**Important:**  
There is **no** EIP-191 / personal_sign prefix.  
The client must sign the raw hash, not:
```bash
"\x19Ethereum Signed Message:\n32" || digest
```

Do **not** use `personal_sign` / `eth_sign` that automatically prepend the prefix.  
Use a function that signs a raw digest (`signDigest`, `signHash`, `Sign`, etc., depending on the library).

---

## 4. Signature format

The client signs `digest` using the authority’s Ethereum private key (secp256k1).

Expected on the wire:

- **65 bytes:** `R || S || V`
  - `R`: 32 bytes  
  - `S`: 32 bytes  
  - `V`: 1 byte (recovery ID)

Server accepts:

- `27` / `28`, or  
- `0` / `1`  

Internally the server normalizes `V` to `0` or `1` before calling `ecrecover`.

The signature must be sent as a **hex string of 65 bytes**:

- Length: **130 hex chars**
- **No `0x` prefix**

Example:  
`"7b3f...<130 hex chars>..."`

---

## 5. JSON payload

Example of a valid `POST /getreport` body:

```json
{
  "chain_id": 1337,
  "app_id": 1,
  "report_id": "04a8e88fc6569a430a13e5657aecf11b8ff633d57a2ad6f4d774b05a468fbafb",
  "salt": "c1d4a94c8aa8c5bd9a2b32f01e7e63d1",
  "nonce": "5f7d1516e21ed7d64f083f60f181e68fdfc08b0a27bb93ef7b8db69153abe59b",
  "timestamp": 1737550000,
  "signature": "7b3f5f1a...<130 hex chars>..."
}
```
Where:

- `report_id` is the 32-byte report identifier encoded as a hex string.
- `nonce` and `timestamp` must be used exactly as returned by the `GET /nonce` endpoint (the client must not modify or regenerate them).
- `signature` is the 65-byte `R || S || V` value encoded as a 130-character hex string, as specified above.
