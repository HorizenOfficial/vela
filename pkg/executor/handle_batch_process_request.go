package executor

import (
	"context"
	"fmt"

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
// It also returns processedCount: how many of the input requests were handled
// (successfully or with an error payload). If processedCount < len(requests) a
// hard failure stopped the batch at request processedCount.
func (e *StatelessExecutor) HandleBatchProcessRequest(ctx context.Context, requests []*common.Request, appState *common.ApplicationState, wasmModule []byte) ([]*common.UpdatePayload, []byte, *common.ApplicationState, []*common.DeanonymizationReport, int, error) {
	if len(requests) == 0 {
		return nil, nil, nil, nil, 0, fmt.Errorf("empty batch")
	}
	e.log.Info("Executor: Processing batch of %d requests for application %d", len(requests), requests[0].ApplicationID)

	// State presence, wasm presence, state root and wasm fingerprint checks.
	// Any failure here is a hard failure: the whole batch stays pending.
	// The state is decrypted once for the whole batch.
	_, decryptedState, err := e.loadVerifiedAppData(appState, wasmModule, requests[0].ApplicationID, "batch")
	if err != nil {
		return nil, nil, nil, nil, 0, err
	}

	// currentSerialized always holds the serialized app data of the last
	// successful request; error payloads chain from it unchanged.
	currentSerialized := decryptedState
	currentStateRoot := appState.StateRoot

	var results []*common.UpdatePayload
	var reports []*common.DeanonymizationReport
	// processed counts the input requests that were handled (a payload was
	// produced). It tracks progress through the input queue and is kept
	// distinct from len(results) so the manager can advance its on-chain
	// cursor by requests consumed, not by payloads emitted.
	processed := 0

	for i, req := range requests {
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
		processed++
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
		return nil, nil, nil, nil, 0, nil
	}

	// One signature covering all entry hashes
	batchSignature, err := e.signBatch(results)
	if err != nil {
		return nil, nil, nil, nil, 0, fmt.Errorf("batch signing failed, discarding batch: %w", err)
	}

	// Encrypt the final state once for the whole batch
	finalState, err := e.buildEncryptedApplicationState(appState.ApplicationID, currentStateRoot, currentSerialized)
	if err != nil {
		return nil, nil, nil, nil, 0, fmt.Errorf("failed to encrypt final state, discarding batch: %w", err)
	}

	e.log.Info("Executor: Successfully processed %d/%d batch requests for application %d", processed, len(requests), appState.ApplicationID)
	return results, batchSignature, finalState, reports, processed, nil
}

// signBatch signs the hash covering all entry hashes with the TEE signing key.
func (e *StatelessExecutor) signBatch(payloads []*common.UpdatePayload) ([]byte, error) {
	batchHash, err := e.BuildBatchMsgHash(payloads)
	if err != nil {
		return nil, fmt.Errorf("failed to create batch message to sign: %w", err)
	}
	return e.signHash(batchHash)
}
