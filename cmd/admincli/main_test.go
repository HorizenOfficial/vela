package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/horizen-pes/pkg/admin"
	"github.com/horizen-pes/pkg/communication"
)

// TestMessageTypeConstants verifies that the locally-duplicated admin message
// type constants stay in sync with the canonical values in pkg/admin.
func TestMessageTypeConstants(t *testing.T) {
	checks := []struct {
		name     string
		local    int
		upstream int
	}{
		{"AdminResponseMessage", adminResponseMessage, int(admin.AdminResponseMessage)},
		{"AdminErrorMessage", adminErrorMessage, int(admin.AdminErrorMessage)},
		{"KeyAttestationRequestMessage", keyAttestationRequest, int(admin.KeyAttestationRequestMessage)},
		{"GetVersionRequestMessage", getVersionRequest, int(admin.GetVersionRequestMessage)},
		{"SetLogLevelRequestMessage", setLogLevelRequest, int(admin.SetLogLevelRequestMessage)},
		{"GetLogLevelRequestMessage", getLogLevelRequest, int(admin.GetLogLevelRequestMessage)},
	}
	for _, c := range checks {
		if c.local != c.upstream {
			t.Errorf("%s: admincli has %d, pkg/admin has %d", c.name, c.local, c.upstream)
		}
	}
}

// TestValidTargets verifies that the target list matches pkg/admin constants.
func TestValidTargets(t *testing.T) {
	expected := []string{admin.TargetManager, admin.TargetExecutor, admin.TargetAll}
	if !reflect.DeepEqual(validTargets, expected) {
		t.Errorf("validTargets = %v, want %v", validTargets, expected)
	}
}

// TestValidLogLevels verifies that the log level list matches pkg/admin.SupportedLogLevels.
func TestValidLogLevels(t *testing.T) {
	parts := strings.Split(admin.SupportedLogLevels, ",")
	expected := make([]string, len(parts))
	for i, p := range parts {
		expected[i] = strings.TrimSpace(p)
	}
	if !reflect.DeepEqual(validLogLevels, expected) {
		t.Errorf("validLogLevels = %v, want %v", validLogLevels, expected)
	}
}

// TestSetLogLevelReqJSONShape verifies that the local setLogLevelReq produces
// the same JSON field names as admin.SetLogLevelRequest.
func TestSetLogLevelReqJSONShape(t *testing.T) {
	local, _ := json.Marshal(setLogLevelReq{Level: "debug", Target: "all"})
	upstream, _ := json.Marshal(admin.SetLogLevelRequest{Level: "debug", Target: "all"})

	var localMap, upstreamMap map[string]any
	json.Unmarshal(local, &localMap)
	json.Unmarshal(upstream, &upstreamMap)

	if !reflect.DeepEqual(localMap, upstreamMap) {
		t.Errorf("SetLogLevelReq JSON mismatch:\n  admincli:  %s\n  pkg/admin: %s", local, upstream)
	}
}

// TestGetLogLevelReqJSONShape verifies that the local getLogLevelReq produces
// the same JSON field names as admin.GetLogLevelRequest.
func TestGetLogLevelReqJSONShape(t *testing.T) {
	local, _ := json.Marshal(getLogLevelReq{Target: "executor"})
	upstream, _ := json.Marshal(admin.GetLogLevelRequest{Target: "executor"})

	var localMap, upstreamMap map[string]any
	json.Unmarshal(local, &localMap)
	json.Unmarshal(upstream, &upstreamMap)

	if !reflect.DeepEqual(localMap, upstreamMap) {
		t.Errorf("GetLogLevelReq JSON mismatch:\n  admincli:  %s\n  pkg/admin: %s", local, upstream)
	}
}

// TestErrorDataJSONShape verifies that the local errorData produces the same
// JSON field names as communication.ErrorData.
func TestErrorDataJSONShape(t *testing.T) {
	local, _ := json.Marshal(errorData{Code: "ERR", Message: "fail"})
	upstream, _ := json.Marshal(communication.ErrorData{Code: "ERR", Message: "fail"})

	var localMap, upstreamMap map[string]any
	json.Unmarshal(local, &localMap)
	json.Unmarshal(upstream, &upstreamMap)

	if !reflect.DeepEqual(localMap, upstreamMap) {
		t.Errorf("ErrorData JSON mismatch:\n  admincli:  %s\n  pkg/comm:  %s", local, upstream)
	}
}
