package authorityservice

import (
	"os"
	"testing"

	"github.com/HorizenOfficial/vela/pkg/authorityservice/deployartifact"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// loadConfigFromFile writes an authorityservice.conf in a fresh temp working
// directory and runs LoadConfig against it, exercising the real parsing path.
func loadConfigFromFile(t *testing.T, conf string) (*Config, error) {
	t.Helper()
	t.Chdir(t.TempDir())
	require.NoError(t, os.WriteFile(confFileName, []byte(conf), 0600))
	return LoadConfig()
}

// TestLoadConfig_DeployArtifactsMaxSize_AboveOverflowBound_Rejected covers the
// start-up half of the upload-limit guard. Above this bound the megabytes-to-bytes
// conversion overflows int64 and goes negative, which switches the HTTP body limit
// off entirely instead of enlarging it, so the value must be refused here rather
// than clamped silently downstream.
func TestLoadConfig_DeployArtifactsMaxSize_AboveOverflowBound_Rejected(t *testing.T) {
	t.Setenv("DEPLOY_ARTIFACTS_MAX_SIZE_MB", "")

	_, err := loadConfigFromFile(t, "DEPLOY_ARTIFACTS_MAX_SIZE_MB=8796093022208\n")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DEPLOY_ARTIFACTS_MAX_SIZE_MB")
}

// TestLoadConfig_DeployArtifactsMaxSize_AtOverflowBound_Accepted pins the boundary:
// the largest convertible value is legal.
func TestLoadConfig_DeployArtifactsMaxSize_AtOverflowBound_Accepted(t *testing.T) {
	t.Setenv("DEPLOY_ARTIFACTS_MAX_SIZE_MB", "")

	cfg, err := loadConfigFromFile(t, "DEPLOY_ARTIFACTS_MAX_SIZE_MB=8796093022207\n")
	require.NoError(t, err)
	require.Equal(t, int64(deployartifact.MaxUploadSizeMB), cfg.DeployArtifactsMaxSizeMB)
}

// TestLoadConfig_DeployArtifactsMaxSize_Negative_Rejected keeps the existing lower
// bound covered alongside the new upper one.
func TestLoadConfig_DeployArtifactsMaxSize_Negative_Rejected(t *testing.T) {
	t.Setenv("DEPLOY_ARTIFACTS_MAX_SIZE_MB", "")

	_, err := loadConfigFromFile(t, "DEPLOY_ARTIFACTS_MAX_SIZE_MB=-1\n")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "DEPLOY_ARTIFACTS_MAX_SIZE_MB")
}
