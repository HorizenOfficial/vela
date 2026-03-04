package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	DeployModeArtifactRef = "artifact_ref"
	artifactIDPrefix      = "sha256:"
	sha256HexLength       = 64
)

// DeployDescriptor defines the v1 deploy payload contract stored in Request.Payload.
type DeployDescriptor struct {
	Mode          string            `json:"mode"`
	ApplicationID ApplicationIdType `json:"applicationId"`
	ArtifactID    string            `json:"artifactId"`
	WasmSHA256    string            `json:"wasmSha256"`
	WasmSize      uint64            `json:"wasmSize"`
}

// DecodeDeployDescriptorStrict decodes deploy descriptor JSON with unknown-field rejection and full validation.
func DecodeDeployDescriptorStrict(payload []byte) (*DeployDescriptor, error) {
	if len(bytes.TrimSpace(payload)) == 0 {
		return nil, errors.New("deploy descriptor payload is empty")
	}

	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()

	var descriptor DeployDescriptor
	if err := decoder.Decode(&descriptor); err != nil {
		return nil, fmt.Errorf("invalid deploy descriptor JSON: %w", err)
	}

	// Enforce that payload is exactly one JSON object.
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		return nil, errors.New("invalid deploy descriptor JSON: trailing data")
	}

	if err := descriptor.Validate(); err != nil {
		return nil, err
	}

	return &descriptor, nil
}

// Validate checks descriptor-level schema and consistency rules.
func (d *DeployDescriptor) Validate() error {
	if d == nil {
		return errors.New("deploy descriptor is nil")
	}

	if d.Mode != DeployModeArtifactRef {
		return fmt.Errorf("invalid deploy descriptor mode %q", d.Mode)
	}

	if d.ApplicationID == 0 {
		return errors.New("invalid deploy descriptor applicationId: must be > 0")
	}

	if !isLowerHexFixedSize(d.WasmSHA256, sha256HexLength) {
		return errors.New("invalid deploy descriptor wasmSha256: must be lowercase hex length 64")
	}

	artifactSHA, err := ParseArtifactID(d.ArtifactID)
	if err != nil {
		return fmt.Errorf("invalid deploy descriptor artifactId: %w", err)
	}

	if artifactSHA != d.WasmSHA256 {
		return errors.New("invalid deploy descriptor: artifactId hash does not match wasmSha256")
	}

	if d.WasmSize == 0 {
		return errors.New("invalid deploy descriptor wasmSize: must be > 0")
	}

	return nil
}

// ValidateApplicationID checks that descriptor applicationId matches the expected request applicationId.
func (d *DeployDescriptor) ValidateApplicationID(expected ApplicationIdType) error {
	if err := d.Validate(); err != nil {
		return err
	}
	if d.ApplicationID != expected {
		return fmt.Errorf("invalid deploy descriptor applicationId: expected %d got %d", expected, d.ApplicationID)
	}
	return nil
}

// BuildArtifactID builds a content-addressed artifactId from a validated sha256 hex string.
func BuildArtifactID(wasmSHA256 string) (string, error) {
	if !isLowerHexFixedSize(wasmSHA256, sha256HexLength) {
		return "", errors.New("invalid wasm sha256: must be lowercase hex length 64")
	}
	return artifactIDPrefix + wasmSHA256, nil
}

// ParseArtifactID validates and extracts the sha256 hex from an artifact id.
func ParseArtifactID(artifactID string) (string, error) {
	if !strings.HasPrefix(artifactID, artifactIDPrefix) {
		return "", errors.New("must start with sha256:")
	}

	sha := strings.TrimPrefix(artifactID, artifactIDPrefix)
	if !isLowerHexFixedSize(sha, sha256HexLength) {
		return "", errors.New("sha256 part must be lowercase hex length 64")
	}

	return sha, nil
}

func isLowerHexFixedSize(v string, size int) bool {
	if len(v) != size {
		return false
	}
	for i := 0; i < len(v); i++ {
		c := v[i]
		if (c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') {
			continue
		}
		return false
	}
	return true
}
