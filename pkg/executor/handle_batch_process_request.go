package executor

import (
	"context"
	"fmt"
	"time"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/appdata"
)

// HandleBatchProcessRequest processes a batch of requests for a single
// application in one round-trip: the state is decrypted once, requests run
// sequentially against the in-memory state, and the state is re-encrypted once
// at the end. Payloads are not signed individually — a single batch signature
// covers all entry hashes (see BuildBatchMsgHash).
//
// Failure handling follows docs/design/BATCH_EXECUTION.md section 7:
//   - soft failure: an unsigned error payload is appended and the batch
//     continues from the unchanged state
//   - hard failure: the batch stops; results for the already-processed
//     requests are returned (possibly none) with a nil error, and the
//     remaining requests stay pending on-chain
//   - batch signing or final state encryption failure: the entire batch is
//     discarded and an error is returned
//
// One payload is returned per handled request, in input order, so len(payloads)
// is how many of the input requests were consumed: len(payloads) < len(requests)
// means a hard failure stopped the batch at request len(payloads).
func (e *StatelessExecutor) HandleBatchProcessRequest(ctx context.Context, requests []*common.Request, appState *common.ApplicationState, wasmModule []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, error) {
	if len(requests) == 0 {
		return nil, nil, nil, nil, fmt.Errorf("empty batch")
	}
	// The batch head identifies the application whose state is loaded below, so
	// it has to exist before anything else. BatchProcessRequestData.Validate
	// rejects nil requests on the wire, but this handler is exported and must
	// not rely on the caller having gone through the protocol layer.
	if requests[0] == nil {
		return nil, nil, nil, nil, fmt.Errorf("nil request at batch head")
	}
	e.log.Info("Executor: Processing batch of %d requests for application %d", len(requests), requests[0].ApplicationID)

	// State presence, wasm presence, state root and wasm fingerprint checks.
	// Any failure here is a hard failure: the whole batch stays pending.
	// The state is decrypted once for the whole batch.
	_, decryptedState, err := e.loadVerifiedAppData(appState, wasmModule, requests[0].ApplicationID, "batch")
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// currentSerialized always holds the serialized app data of the last
	// successful request; error payloads chain from it unchanged.
	currentSerialized := decryptedState
	currentStateRoot := appState.StateRoot

	// results holds one payload per handled request, in input order, so its
	// length doubles as the count of input requests consumed by the batch.
	var results []*common.UpdatePayload
	var reports []*common.DeanonymizationReport

	guestBound := e.config.guestExecutionBound()

	for i, req := range requests {
		// The whole batch shares one caller budget, so stop before starting a request
		// the remaining budget cannot cover. Without this the guest would run until
		// the budget expired, be interrupted, classified as host-side abandonment and
		// discarded — burning a full guest bound of execution and evicting the module
		// (an interrupt is a trap) only for the work to be repeated on the next poll.
		// Stopping early settles exactly the same requests for free.
		//
		// The threshold is one bound, not two: a request that also carries a deposit
		// makes two guest calls and may still be cut short, which simply falls back to
		// the behaviour above. Requiring two would stop earlier than necessary and cost
		// throughput on the common single-call case.
		if deadline, ok := ctx.Deadline(); ok && time.Until(deadline) < guestBound {
			e.log.Info("Executor: batch stopping after %d/%d requests: %v left, less than the %v guest bound",
				len(results), len(requests), time.Until(deadline), guestBound)
			break
		}

		// A nil request cannot be executed and cannot be reported on-chain, so
		// it is a hard failure like any other: stop and keep what was done.
		if req == nil {
			e.log.Warn("Executor: batch stopped at request %d: nil request", i)
			break
		}

		// A batch is scoped to one application; a request for another
		// application inside the batch is evidence of tampering.
		if req.ApplicationID != appState.ApplicationID {
			e.log.Warn("Executor: batch stopped at request %d (%s): applicationId %d does not match batch application %d",
				i, req.RequestID, req.ApplicationID, appState.ApplicationID)
			break
		}

		// Each request works on a fresh copy of the last successful state, so
		// partial mutations from a failed request never leak into the next one.
		workData, err := appdata.DeserializeAppData(currentSerialized)
		if err != nil {
			e.log.Error("Executor: batch stopped at request %d (%s): failed to deserialize state: %v", i, req.RequestID, err)
			break
		}

		e.log.Info("Executor: Processing batch request %s for application %d", req.RequestID, req.ApplicationID)
		outcome, err := e.executeRequest(ctx, req, workData, currentStateRoot, wasmModule)
		if err != nil {
			// Hard failure: nothing can be submitted for this request and the
			// application queue is FIFO, so the batch stops here.
			e.log.Warn("Executor: batch stopped at request %d (%s): %v", i, req.RequestID, err)
			break
		}

		// Success outcomes always carry the serialized app data (Serialize never
		// returns empty bytes); a nil value here is a bug in executeRequest.
		// Hard failure: stop the batch before including this payload.
		if outcome.payload.ErrorCode == 0 && outcome.newSerialized == nil {
			e.log.Error("Executor: batch stopped at request %d (%s): success outcome without serialized app data", i, req.RequestID)
			break
		}

		results = append(results, outcome.payload)
		if outcome.payload.ErrorCode == 0 {
			currentSerialized = outcome.newSerialized
			currentStateRoot = outcome.payload.NewStateRoot
			if outcome.report != nil {
				reports = append(reports, outcome.report)
			}
		}
	}

	if len(results) == 0 {
		// Hard failure on the very first request: nothing to submit, the
		// manager retries on the next poll.
		return nil, nil, nil, nil, nil
	}

	// One signature covering all entry hashes
	batchSignature, err := e.signBatch(results)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("batch signing failed, discarding batch: %w", err)
	}

	// Encrypt the final state once for the whole batch
	finalState, err := e.buildEncryptedApplicationState(appState.ApplicationID, currentStateRoot, currentSerialized)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to encrypt final state, discarding batch: %w", err)
	}

	e.log.Info("Executor: Successfully processed %d/%d batch requests for application %d", len(results), len(requests), appState.ApplicationID)
	return results, batchSignature, finalState, reports, nil
}

// signBatch signs the hash covering all entry hashes with the TEE signing key.
func (e *StatelessExecutor) signBatch(payloads []*common.UpdatePayload) ([]byte, error) {
	batchHash, err := e.BuildBatchMsgHash(payloads)
	if err != nil {
		return nil, fmt.Errorf("failed to create batch message to sign: %w", err)
	}
	return e.signHash(batchHash)
}
