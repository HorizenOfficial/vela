package manager

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/common/apperrors"
)

const (
	artifactBlobsFolder = "blobs"
)

const (
	deployFailureMsgGeneric = "failed to deploy application"
)

type DeployErrorKind string

const (
	DeployErrorDescriptorInvalid    DeployErrorKind = "ARTIFACT_DESCRIPTOR_INVALID"
	DeployErrorArtifactLoadFailed   DeployErrorKind = "ARTIFACT_LOAD_FAILED"
	DeployErrorArtifactNotFound     DeployErrorKind = "ARTIFACT_NOT_FOUND"
	DeployErrorArtifactHashMismatch DeployErrorKind = "ARTIFACT_HASH_MISMATCH"
)

type deployResolutionError struct {
	kind      DeployErrorKind
	transient bool
	cause     error
}

func (e *deployResolutionError) Error() string {
	if e == nil {
		return ""
	}
	if e.cause == nil {
		return string(e.kind)
	}
	return fmt.Sprintf("%s: %v", e.kind, e.cause)
}

func (e *deployResolutionError) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func newDeployResolutionError(kind DeployErrorKind, transient bool, cause error) *deployResolutionError {
	return &deployResolutionError{
		kind:      kind,
		transient: transient,
		cause:     cause,
	}
}

func (m *SecureProcessorManager) resolveDeployWASM(_ context.Context, req *common.Request) ([]byte, error) {
	descriptor, err := common.DecodeDeployDescriptorStrict(req.Payload)
	if err != nil {
		return nil, newDeployResolutionError(DeployErrorDescriptorInvalid, false, err)
	}

	if err := descriptor.ValidateApplicationID(req.ApplicationID); err != nil {
		return nil, newDeployResolutionError(DeployErrorDescriptorInvalid, false, err)
	}

	artifactPath := filepath.Join(m.config.ArtifactsPath, artifactBlobsFolder, descriptor.WasmSHA256+".wasm")
	wasmBytes, err := os.ReadFile(artifactPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, newDeployResolutionError(DeployErrorArtifactNotFound, true, err)
		}
		return nil, newDeployResolutionError(DeployErrorArtifactLoadFailed, true, err)
	}

	hash := sha256.Sum256(wasmBytes)
	if got := hex.EncodeToString(hash[:]); got != descriptor.WasmSHA256 {
		return nil, newDeployResolutionError(
			DeployErrorArtifactHashMismatch,
			false,
			fmt.Errorf("artifact hash mismatch expected=%s got=%s", descriptor.WasmSHA256, got),
		)
	}

	return wasmBytes, nil
}

func (m *SecureProcessorManager) submitDeterministicDeployFailure(
	ctx context.Context,
	req *common.Request,
	stateRoot [32]byte,
	kind DeployErrorKind,
	cause error,
) error {
	failure := mapDeployErrorToFailure(kind)
	updatePayload, err := m.executorClient.SendBuildErrorPayloadRequest(ctx, req, stateRoot, failure)
	if err != nil {
		return fmt.Errorf("failed to create signed deterministic deploy error payload for kind %s: %w", kind, err)
	}

	m.log.Warn(
		"Manager: deterministic deploy failure requestId=%s applicationId=%d kind=%s error=%v",
		req.RequestID,
		req.ApplicationID,
		kind,
		cause,
	)

	return m.submitStateOnChain(ctx, updatePayload)
}

func stateRootForDeployFailure(state *common.ApplicationState) [32]byte {
	if state == nil {
		return [32]byte{}
	}
	return state.StateRoot
}

func mapDeployErrorToFailure(kind DeployErrorKind) *apperrors.RequestFailure {
	switch kind {
	case DeployErrorDescriptorInvalid:
		return apperrors.New(apperrors.CodeInternalFallback, deployFailureMsgGeneric)
	case DeployErrorArtifactLoadFailed, DeployErrorArtifactNotFound, DeployErrorArtifactHashMismatch:
		return apperrors.New(apperrors.CodeFailedLoadingOrGettingModule, deployFailureMsgGeneric)
	default:
		return apperrors.New(apperrors.CodeInternalFallback, deployFailureMsgGeneric)
	}
}
