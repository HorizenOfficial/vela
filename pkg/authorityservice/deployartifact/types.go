package deployartifact

import (
	"errors"
	"time"
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
	WasmSize   uint64 `json:"wasmSize"`
}

type Metadata struct {
	ArtifactID string    `json:"artifactId"`
	WasmSHA256 string    `json:"wasmSha256"`
	WasmSize   uint64    `json:"wasmSize"`
	CreatedAt  time.Time `json:"createdAt"`
}
