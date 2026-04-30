package common

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeDeployDescriptorStrict_Valid(t *testing.T) {
	sha := strings.Repeat("a", 64)
	payload := []byte(`{
		"mode":"artifact_ref",
		"artifactId":"sha256:` + sha + `",
		"wasmSha256":"` + sha + `"
	}`)

	descriptor, err := DecodeDeployDescriptorStrict(payload)
	require.NoError(t, err)
	require.Equal(t, DeployModeArtifactRef, descriptor.Mode)
	require.Equal(t, "sha256:"+sha, descriptor.ArtifactID)
	require.Equal(t, sha, descriptor.WasmSHA256)
}

func TestDecodeDeployDescriptorStrict_RejectsEmptyPayload(t *testing.T) {
	_, err := DecodeDeployDescriptorStrict([]byte(" \n\t"))
	require.ErrorContains(t, err, "payload is empty")
}

func TestDecodeDeployDescriptorStrict_RejectsNonJSON(t *testing.T) {
	_, err := DecodeDeployDescriptorStrict([]byte("not-json"))
	require.ErrorContains(t, err, "invalid deploy descriptor JSON")
}

func TestDecodeDeployDescriptorStrict_RejectsUnknownField(t *testing.T) {
	sha := strings.Repeat("a", 64)
	payload := []byte(`{
		"mode":"artifact_ref",
		"artifactId":"sha256:` + sha + `",
		"wasmSha256":"` + sha + `",
		"version":1
	}`)

	_, err := DecodeDeployDescriptorStrict(payload)
	require.ErrorContains(t, err, "unknown field")
}

func TestDecodeDeployDescriptorStrict_RejectsLegacyApplicationIDField(t *testing.T) {
	sha := strings.Repeat("a", 64)
	payload := []byte(`{
		"mode":"artifact_ref",
		"applicationId":1,
		"artifactId":"sha256:` + sha + `",
		"wasmSha256":"` + sha + `"
	}`)

	_, err := DecodeDeployDescriptorStrict(payload)
	require.ErrorContains(t, err, "unknown field")
}

func TestDecodeDeployDescriptorStrict_RejectsInvalidMode(t *testing.T) {
	sha := strings.Repeat("a", 64)
	payload := []byte(`{
		"mode":"inline",
		"artifactId":"sha256:` + sha + `",
		"wasmSha256":"` + sha + `"
	}`)

	_, err := DecodeDeployDescriptorStrict(payload)
	require.ErrorContains(t, err, "invalid deploy descriptor mode")
}

func TestDecodeDeployDescriptorStrict_RejectsMissingArtifactID(t *testing.T) {
	sha := strings.Repeat("a", 64)
	payload := []byte(`{
		"mode":"artifact_ref",
		"wasmSha256":"` + sha + `"
	}`)

	_, err := DecodeDeployDescriptorStrict(payload)
	require.ErrorContains(t, err, "invalid deploy descriptor artifactId")
}

func TestDecodeDeployDescriptorStrict_RejectsInvalidArtifactIDFormat(t *testing.T) {
	sha := strings.Repeat("a", 64)
	payload := []byte(`{
		"mode":"artifact_ref",
		"artifactId":"invalid:` + sha + `",
		"wasmSha256":"` + sha + `"
	}`)

	_, err := DecodeDeployDescriptorStrict(payload)
	require.ErrorContains(t, err, "invalid deploy descriptor artifactId")
}

func TestDecodeDeployDescriptorStrict_RejectsArtifactHashMismatch(t *testing.T) {
	artifactSHA := strings.Repeat("a", 64)
	wasmSHA := strings.Repeat("b", 64)
	payload := []byte(`{
		"mode":"artifact_ref",
		"artifactId":"sha256:` + artifactSHA + `",
		"wasmSha256":"` + wasmSHA + `"
	}`)

	_, err := DecodeDeployDescriptorStrict(payload)
	require.ErrorContains(t, err, "does not match wasmSha256")
}

func TestDecodeDeployDescriptorStrict_RejectsUppercaseWasmHash(t *testing.T) {
	sha := strings.Repeat("A", 64)
	payload := []byte(`{
		"mode":"artifact_ref",
		"artifactId":"sha256:` + strings.Repeat("a", 64) + `",
		"wasmSha256":"` + sha + `"
	}`)

	_, err := DecodeDeployDescriptorStrict(payload)
	require.ErrorContains(t, err, "invalid deploy descriptor wasmSha256")
}

func TestBuildArtifactID(t *testing.T) {
	sha := strings.Repeat("a", 64)
	artifactID, err := BuildArtifactID(sha)
	require.NoError(t, err)
	require.Equal(t, "sha256:"+sha, artifactID)
}

func TestBuildArtifactID_InvalidHash(t *testing.T) {
	_, err := BuildArtifactID("invalid")
	require.ErrorContains(t, err, "invalid wasm sha256")
}

func TestParseArtifactID(t *testing.T) {
	sha := strings.Repeat("a", 64)
	got, err := ParseArtifactID("sha256:" + sha)
	require.NoError(t, err)
	require.Equal(t, sha, got)
}

func TestParseArtifactID_InvalidFormat(t *testing.T) {
	_, err := ParseArtifactID("sha:abc")
	require.ErrorContains(t, err, "must start with sha256")
}

func TestDecodeDeployDescriptorStrict_WithConstructorParams(t *testing.T) {
	sha := strings.Repeat("a", 64)
	payload := []byte(`{
		"mode":"artifact_ref",
		"artifactId":"sha256:` + sha + `",
		"wasmSha256":"` + sha + `",
		"constructorParams":{"allowedTokens":["0xdead"]}
	}`)

	descriptor, err := DecodeDeployDescriptorStrict(payload)
	require.NoError(t, err)
	require.NotNil(t, descriptor.ConstructorParams)
	require.Contains(t, string(descriptor.ConstructorParams), "allowedTokens")
}

func TestDecodeDeployDescriptorStrict_WithoutConstructorParams(t *testing.T) {
	sha := strings.Repeat("a", 64)
	payload := []byte(`{
		"mode":"artifact_ref",
		"artifactId":"sha256:` + sha + `",
		"wasmSha256":"` + sha + `"
	}`)

	descriptor, err := DecodeDeployDescriptorStrict(payload)
	require.NoError(t, err)
	require.Nil(t, descriptor.ConstructorParams)
}
