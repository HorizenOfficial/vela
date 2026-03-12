package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/HorizenOfficial/vela/pkg/common"
)

const (
	artifactBlobsFolder = "blobs"
)

type DeployErrorKind string

const (
	DeployErrorDescriptorInvalid  DeployErrorKind = "ARTIFACT_DESCRIPTOR_INVALID"
	DeployErrorArtifactLoadFailed DeployErrorKind = "ARTIFACT_LOAD_FAILED"
	DeployErrorArtifactNotFound   DeployErrorKind = "ARTIFACT_NOT_FOUND"
)

type deployResolutionError struct {
	kind  DeployErrorKind
	cause error
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

func newDeployResolutionError(kind DeployErrorKind, cause error) *deployResolutionError {
	return &deployResolutionError{
		kind:  kind,
		cause: cause,
	}
}

func (m *SecureProcessorManager) resolveDeployWASM(_ context.Context, req *common.Request) ([]byte, error) {
	artifactSHA, err := extractDeployArtifactSHA(req.Payload)
	if err != nil {
		return nil, newDeployResolutionError(DeployErrorDescriptorInvalid, err)
	}

	artifactPath := filepath.Join(m.config.ArtifactsPath, artifactBlobsFolder, artifactSHA+".wasm")
	wasmBytes, err := os.ReadFile(artifactPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, newDeployResolutionError(DeployErrorArtifactNotFound, err)
		}
		return nil, newDeployResolutionError(DeployErrorArtifactLoadFailed, err)
	}

	return wasmBytes, nil
}

type deployArtifactLocator struct {
	ArtifactID string `json:"artifactId"`
	WasmSHA256 string `json:"wasmSha256"`
}

func extractDeployArtifactSHA(payload []byte) (string, error) {
	if len(bytes.TrimSpace(payload)) == 0 {
		return "", errors.New("deploy descriptor payload is empty")
	}

	var locator deployArtifactLocator
	if err := json.Unmarshal(payload, &locator); err != nil {
		return "", fmt.Errorf("invalid deploy descriptor JSON: %w", err)
	}

	if sha, err := common.ParseArtifactID(strings.TrimSpace(locator.ArtifactID)); err == nil {
		return sha, nil
	}

	wasmSHA := strings.TrimSpace(locator.WasmSHA256)
	if _, err := common.BuildArtifactID(wasmSHA); err == nil {
		return wasmSHA, nil
	}

	return "", errors.New("artifact locator not found in deploy payload")
}
