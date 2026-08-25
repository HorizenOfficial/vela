package executor

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/big"

	"github.com/HorizenOfficial/vela/pkg/common"
	"github.com/HorizenOfficial/vela/pkg/common/appdata"
	"github.com/HorizenOfficial/vela/pkg/common/apperrors"
	cryptotypes "github.com/HorizenOfficial/vela/pkg/common/crypto"
	"github.com/HorizenOfficial/vela/pkg/crypto"
)

// requestOutcome is the result of executing one request.
// Exactly one of the following holds:
//   - success: payload.ErrorCode == 0 and newSerialized is the serialized
//     app data after the request, never nil (report optionally set)
//   - soft failure: payload is an error payload (ErrorCode != 0) and
//     newSerialized is nil
//
// Callers discriminate on payload.ErrorCode; a success outcome with a nil
// newSerialized is a bug and must be treated as a hard failure.
type requestOutcome struct {
	payload       *common.UpdatePayload
	newSerialized []byte
	report        *common.DeanonymizationReport
}

// executeRequest runs a single request against the given app data. It is the
// shared core of HandleProcessRequest and HandleBatchProcessRequest: it does
// not decrypt/encrypt the application state and does not sign the payload —
// each caller handles those steps according to its own flow.
//
// Returns:
//   - (outcome with error payload, nil) on soft failure — state unchanged
//   - (outcome with success payload + new serialized app data, nil) on success
//   - (nil, err) on hard failure — nothing can be submitted for this request
//
// workData may be partially mutated on soft failure; the caller discards it.
func (e *StatelessExecutor) executeRequest(ctx context.Context, req *common.Request, workData *appdata.AppData, prevStateRoot [32]byte, wasmModule []byte) (*requestOutcome, error) {
	softFailure := func(failure *apperrors.RequestFailure) (*requestOutcome, error) {
		// One class of runtime failure must NOT become an error payload: a guest
		// interrupted because the HOST stopped waiting for it (executor shutdown, a
		// cancelled context, an expired caller execution budget). Nothing is known to
		// be wrong with the request, so signing it would charge the user
		// MinFeePerRequest for our own decision to stop, and would settle a request
		// that was never actually judged. Reporting it as a hard failure instead
		// leaves it pending: the single-request path returns a plain error and the
		// manager retries on the next poll, while the batch path stops here and
		// submits whatever it had already collected.
		//
		// This is the ONLY place the distinction is applied, so it holds for both
		// paths at once. Contrast a guest that exhausted its own wall-clock bound
		// (CodeGuestExecutionTimeout), which stays a soft failure and IS signed: the
		// request queue is FIFO and only advances when an update is submitted, so
		// retrying a non-terminating guest forever would stall every application.
		if transientErr, ok := e.transientFailure(req, failure); ok {
			return nil, transientErr
		}
		e.log.Info("Executor: Returning error payload for request %s (error code: %d)", req.RequestID, failure.Category())
		return &requestOutcome{payload: e.buildUnsignedErrorPayload(req, prevStateRoot, failure)}, nil
	}

	if err := e.validateRequest(req); err != nil {
		return nil, err
	}

	if req.RequestType != common.Process && req.RequestType != common.AssociateKey && req.RequestType != common.Deanonymize && req.RequestType != common.TrustProcess {
		return nil, fmt.Errorf("unsupported request type: %s", req.RequestType)
	}

	// If the request contains a deposit, handle it first
	var tempState = workData.GetAppState()
	var depositEvents []common.PlainEvent
	var depositAppEvents []common.AppEvent
	var totalFuel *big.Int = big.NewInt(0)
	if req.AssetAmount.ToInt().Sign() > 0 {
		newState, depEvents, depAppEvents, reqFuel, failure := e.runtime.Deposit(ctx, req.ApplicationID, req.Sender, req.TokenAddress, req.AssetAmount.ToInt(), tempState, wasmModule)
		if failure != nil {
			return softFailure(failure)
		}
		tempState = newState
		depositEvents = depEvents
		depositAppEvents = depAppEvents
		totalFuel = totalFuel.Add(totalFuel, reqFuel)
		e.log.Info("Executor: Successfully processed deposit for request %s", req.RequestID)
	}

	applicationFee := new(big.Int).Mul(totalFuel, e.config.FuelPricePerUnit)
	if req.MaxFeeValue.ToInt().Cmp(applicationFee) < 0 {
		return softFailure(apperrors.New(apperrors.CodeInsufficientFuel,
			fmt.Sprintf("insufficient fuel: required %s wei, provided %s wei",
				applicationFee.String(),
				req.MaxFeeValue.ToInt().String(),
			)))
	}

	var events []common.PlainEvent
	var appEvents []common.AppEvent
	var withdrawals []common.Withdrawal
	var reportData []byte

	if req.RequestType == common.AssociateKey {
		//request of type associate key: the payload contains the new key (plaintext) and optionally an encrypted seed
		e.log.Info("Associating new key - RequestID %s", req.RequestID)

		const keyOnlyPayloadSize = 133
		// encrypted seed = AES-GCM nonce (12) + seed (65) + tag (16) = 93 bytes
		const encryptedSeedSize = 12 + appdata.SeedStore_ValSize + 16
		const keyWithEncryptedSeedPayloadSize = keyOnlyPayloadSize + encryptedSeedSize // 226 bytes
		if len(req.Payload) != keyOnlyPayloadSize && len(req.Payload) != keyWithEncryptedSeedPayloadSize {
			return nil, fmt.Errorf("invalid payload length")
		}

		keyToAssociate, err := cryptotypes.NewPublicKeyP521(req.Payload[:keyOnlyPayloadSize])
		if err != nil {
			e.log.Error("Executor: request %s (app %d) failed (%s): %v",
				req.RequestID, req.ApplicationID, apperrors.CodeParsingKeyError.Code, err)
			return softFailure(apperrors.New(apperrors.CodeParsingKeyError,
				fmt.Sprintf("failed to parse keyP521 in request payload: %v", err)))
		}

		workData.AddKey(req.Sender, *keyToAssociate)

		if len(req.Payload) == keyWithEncryptedSeedPayloadSize {
			// Decrypt the encrypted seed using ECDH(enclave_priv_P521, user_pub_P521)
			encryptedSeed := req.Payload[keyOnlyPayloadSize:keyWithEncryptedSeedPayloadSize]
			seed, err := crypto.Decrypt(keyToAssociate, &e.keySet.CommunicationKey, encryptedSeed)
			if err != nil {
				e.log.Error("Executor: request %s (app %d) failed (%s): %v",
					req.RequestID, req.ApplicationID, apperrors.CodeParsingKeyError.Code, err)
				return softFailure(apperrors.New(apperrors.CodeParsingKeyError,
					fmt.Sprintf("seed decryption failed: %v", err)))
			}
			if err := workData.AddSeed(req.Sender, seed); err != nil {
				return nil, fmt.Errorf("failed to add seed for request %s: %w", req.RequestID, err)
			}
		}

		totalFuel = totalFuel.Add(totalFuel, big.NewInt(10))
	} else {
		//any other case: forward the payload to the WASM to obtain the new state.
		//Process and Deanonymize payloads are encrypted toward the enclave and must be
		//decrypted first; TrustProcess payloads are sent in clear text and forwarded as-is.

		wasmPayload := req.Payload
		if req.RequestType != common.TrustProcess {
			decryptedPayload, failure := e.decryptPayload(&e.keySet.CommunicationKey, req.Payload, req.Sender, workData.GetKeyStore())
			if failure != nil {
				return softFailure(failure)
			}
			wasmPayload = decryptedPayload
		}

		// Invoke WASM method to process the request
		newState, reqEvents, reqAppEvents, reqWithdrawals, reqReportData, reqFuel, failure := e.runtime.ProcessRequest(ctx, req.ApplicationID, req.Sender, req.RequestType, wasmPayload, tempState, wasmModule)
		if failure != nil {
			return softFailure(failure)
		}
		tempState = newState
		events = reqEvents
		appEvents = reqAppEvents
		withdrawals = reqWithdrawals
		totalFuel = totalFuel.Add(totalFuel, reqFuel)

		// Validate report generation rules:
		// 1. Reports must only be generated for Deanonymize requests
		// 2. Deanonymize requests must always generate a report
		if len(reqReportData) > 0 {
			if req.RequestType != common.Deanonymize {
				return softFailure(apperrors.New(apperrors.CodeRequestFuncFailed,
					fmt.Sprintf("WASM module generated unexpected report for request type %s", req.RequestType)))
			}
			reportData = reqReportData
		} else if req.RequestType == common.Deanonymize {
			return softFailure(apperrors.New(apperrors.CodeRequestFuncFailed,
				"WASM module failed to generate report for Deanonymize request"))
		}
	}

	// Check if there is enough ETH to cover the fuel costs.
	// TRUSTPROCESS requests have maxFeeValue=0 (enqueued by the on-chain trigger
	// contract with no user-provided fee). Skip fee adequacy checks for them;
	// their authenticity is established on-chain.
	if req.RequestType == common.TrustProcess {
		applicationFee = big.NewInt(0)
	} else {
		applicationFee = new(big.Int).Mul(totalFuel, e.config.FuelPricePerUnit)

		// Application fee must be minimum fee at least
		if applicationFee.Cmp(e.config.MinFeePerRequest) < 0 {
			applicationFee = new(big.Int).Set(e.config.MinFeePerRequest)
		}

		if req.MaxFeeValue.ToInt().Cmp(applicationFee) < 0 {
			return softFailure(apperrors.New(apperrors.CodeInsufficientFuel,
				fmt.Sprintf("insufficient fuel: required %s wei, provided %s wei",
					applicationFee.String(),
					req.MaxFeeValue.ToInt().String(),
				)))
		}
	}

	// Compute refundAmount = req.MaxFeeValue - applicationFee
	refundAmount := new(big.Int).Sub(req.MaxFeeValue.ToInt(), applicationFee)

	//set the updated state
	workData.SetAppState(tempState)

	//increment appNonce
	workData.IncrementNonce()

	//serialize the new app data
	newAppData, err := workData.Serialize()
	if err != nil {
		e.log.Error("Executor: request %s (app %d) failed (%s): %v",
			req.RequestID, req.ApplicationID, apperrors.CodeAppDataSerializationFailure.Code, err)
		return softFailure(apperrors.New(apperrors.CodeAppDataSerializationFailure,
			fmt.Sprintf("failed to serialize new app data: %v", err)))
	}

	// Encrypt events if they are not empty
	events = append(depositEvents, events...)
	appEvents = append(depositAppEvents, appEvents...)
	encryptedEvents, failure, err := e.encryptEvents(ctx, events, req.ApplicationID, &e.keySet.CommunicationKey, e.server, workData.GetKeyStore(), workData.GetEventSeedStore())
	if failure != nil {
		return softFailure(failure)
	}
	if err != nil {
		e.log.Error("Executor: error encrypting events: %v", err)
		return nil, err
	}

	// Create appdata root hash
	newStateRoot := sha256.Sum256(newAppData)

	// Create the update payload (unsigned — the caller signs it)
	updatePayload := &common.UpdatePayload{
		ApplicationID:  req.ApplicationID,
		RequestID:      req.RequestID,
		PrevStateRoot:  prevStateRoot,
		NewStateRoot:   newStateRoot,
		Events:         encryptedEvents,
		AppEvents:      appEvents,
		Withdrawals:    withdrawals,
		RefundAmount:   common.ToBig(refundAmount),
		ApplicationFee: common.ToBig(applicationFee),
	}

	// If a report was generated, encrypt it and create the DeanonymizationReport
	var deanonymizationReport *common.DeanonymizationReport
	if len(reportData) > 0 {
		encryptedReport, failure, err := e.encryptDeanonymizationReport(
			req.ApplicationID,
			req.RequestID,
			&e.keySet.CommunicationKey,
			req.Sender,
			reportData,
			workData.GetKeyStore(),
		)
		if failure != nil {
			return softFailure(failure)
		}
		if err != nil {
			e.log.Error("Executor: error encrypting deanonymization report: %v", err)
			return nil, err
		}

		deanonymizationReport = &common.DeanonymizationReport{
			ApplicationID:   req.ApplicationID,
			ReportID:        req.RequestID,
			EncryptedReport: encryptedReport,
			Authority:       req.Sender,
		}
		e.log.Info("Executor: Successfully generated deanonymization report %s", req.RequestID)
	}

	e.log.Info("Executor: Successfully processed request %s", req.RequestID)
	return &requestOutcome{
		payload:       updatePayload,
		newSerialized: newAppData,
		report:        deanonymizationReport,
	}, nil
}
