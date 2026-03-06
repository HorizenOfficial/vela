package deployartifact

import (
	"errors"
)

const (
	artifactIDPrefix = "sha256:"
)

var (
	ErrEmptyWASM = errors.New("wasm payload is empty")
)

type UploadResponse struct {
	ArtifactID string `json:"artifactId"`
	WasmSHA256 string `json:"wasmSha256"`
}
