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
		local    string
		upstream string
	}{
		{"AdminResponseMessage", adminResponseMessage, string(admin.AdminResponseMessage)},
		{"AdminErrorMessage", adminErrorMessage, string(admin.AdminErrorMessage)},
		{"KeyAttestationRequestMessage", keyAttestationRequest, string(admin.KeyAttestationRequestMessage)},
		{"GetVersionRequestMessage", getVersionRequest, string(admin.GetVersionRequestMessage)},
		{"SetLogLevelRequestMessage", setLogLevelRequest, string(admin.SetLogLevelRequestMessage)},
		{"GetLogLevelRequestMessage", getLogLevelRequest, string(admin.GetLogLevelRequestMessage)},
	}
	for _, c := range checks {
		if c.local != c.upstream {
			t.Errorf("%s: admincli has %q, pkg/admin has %q", c.name, c.local, c.upstream)
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
	local, _ := json.Marshal(setLogLevelReq{Level: "debug"})
	upstream, _ := json.Marshal(admin.SetLogLevelRequest{Level: "debug"})

	var localMap, upstreamMap map[string]any
	json.Unmarshal(local, &localMap)
	json.Unmarshal(upstream, &upstreamMap)

	if !reflect.DeepEqual(localMap, upstreamMap) {
		t.Errorf("SetLogLevelReq JSON mismatch:\n  admincli:  %s\n  pkg/admin: %s", local, upstream)
	}
}

// TestAdminMessageJSONShape verifies that the local adminMessage produces
// the same JSON field names as admin.AdminMessage.
func TestAdminMessageJSONShape(t *testing.T) {
	local, _ := json.Marshal(adminMessage{
		Type:   setLogLevelRequest,
		Target: "manager",
		Data:   json.RawMessage(`{"level":"debug"}`),
	})
	upstream, _ := json.Marshal(admin.AdminMessage{
		Type:   admin.SetLogLevelRequestMessage,
		Target: "manager",
		Data:   json.RawMessage(`{"level":"debug"}`),
	})

	var localMap, upstreamMap map[string]any
	json.Unmarshal(local, &localMap)
	json.Unmarshal(upstream, &upstreamMap)

	if !reflect.DeepEqual(localMap, upstreamMap) {
		t.Errorf("AdminMessage JSON mismatch:\n  admincli:  %s\n  pkg/admin: %s", local, upstream)
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
