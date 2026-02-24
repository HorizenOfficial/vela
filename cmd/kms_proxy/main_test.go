package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func setEnv(t *testing.T, key, value string) {
	t.Helper()
	old, ok := os.LookupEnv(key)
	if value == "" {
		_ = os.Unsetenv(key)
	} else {
		_ = os.Setenv(key, value)
	}
	t.Cleanup(func() {
		if ok {
			_ = os.Setenv(key, old)
		} else {
			_ = os.Unsetenv(key)
		}
	})
}

func TestParseOperation(t *testing.T) {
	op, err := parseOperation("KMS.GenerateDataKey")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if op != "generatedatakey" {
		t.Fatalf("expected generatedatakey, got %q", op)
	}
}

func TestParseOperationInvalid(t *testing.T) {
	if _, err := parseOperation("InvalidTarget"); err == nil {
		t.Fatalf("expected error for invalid target")
	}
	if _, err := parseOperation("KMS."); err == nil {
		t.Fatalf("expected error for empty operation")
	}
}

func TestHasRecipient(t *testing.T) {
	if hasRecipient(map[string]any{}) {
		t.Fatalf("expected false for missing recipient")
	}
	if hasRecipient(map[string]any{"Recipient": "nope"}) {
		t.Fatalf("expected false for invalid recipient type")
	}
	if hasRecipient(map[string]any{"Recipient": map[string]any{"AttestationDocument": "x"}}) {
		t.Fatalf("expected false when KeyEncryptionAlgorithm missing")
	}
	if !hasRecipient(map[string]any{"Recipient": map[string]any{
		"AttestationDocument":    "x",
		"KeyEncryptionAlgorithm": "RSAES_OAEP_SHA_256",
	}}) {
		t.Fatalf("expected true for valid recipient")
	}
}

func TestExtractKeyID(t *testing.T) {
	if _, ok := extractKeyID(map[string]any{}); ok {
		t.Fatalf("expected false for missing KeyId")
	}
	if _, ok := extractKeyID(map[string]any{"KeyId": 123}); ok {
		t.Fatalf("expected false for non-string KeyId")
	}
	if _, ok := extractKeyID(map[string]any{"KeyId": "   "}); ok {
		t.Fatalf("expected false for empty KeyId")
	}
	if id, ok := extractKeyID(map[string]any{"KeyId": "arn:good"}); !ok || id != "arn:good" {
		t.Fatalf("expected arn:good, got %q (ok=%v)", id, ok)
	}
}

func TestParseCSVSetLower(t *testing.T) {
	set := parseCSVSetLower("GenerateDataKey, Decrypt ")
	if _, ok := set["generatedatakey"]; !ok {
		t.Fatalf("expected generatedatakey in set")
	}
	if _, ok := set["decrypt"]; !ok {
		t.Fatalf("expected decrypt in set")
	}
}

func TestParseCSVSetExact(t *testing.T) {
	set := parseCSVSetExact(" arn:one ,arn:two ")
	if _, ok := set["arn:one"]; !ok {
		t.Fatalf("expected arn:one in set")
	}
	if _, ok := set["arn:two"]; !ok {
		t.Fatalf("expected arn:two in set")
	}
}

func TestLoadConfigRequiresAllowlist(t *testing.T) {
	setEnv(t, "KMS_PROXY_ALLOWED_ACTIONS", "GenerateDataKey")
	setEnv(t, "KMS_PROXY_ALLOWED_KEY_ARNS", "")
	setEnv(t, "KMS_PROXY_ALLOW_ALL_KEYS", "false")

	if _, err := loadConfig(); err == nil {
		t.Fatalf("expected error when allowlist is empty and allow-all is false")
	}
}

func TestLoadConfigAllowAllKeys(t *testing.T) {
	setEnv(t, "KMS_PROXY_ALLOWED_ACTIONS", "GenerateDataKey")
	setEnv(t, "KMS_PROXY_ALLOWED_KEY_ARNS", "")
	setEnv(t, "KMS_PROXY_ALLOW_ALL_KEYS", "true")

	if _, err := loadConfig(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServeHTTP_OperationNotAllowed(t *testing.T) {
	p := &proxy{
		allowedActions:   map[string]struct{}{},
		allowedKeyARNs:   map[string]struct{}{},
		enforceRecipient: true,
	}

	req := httptest.NewRequest(http.MethodPost, "http://example.com/", strings.NewReader(`{}`))
	req.Header.Set("X-Amz-Target", "KMS.GenerateDataKey")
	rr := httptest.NewRecorder()

	p.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}
}

func TestServeHTTP_MissingRecipient(t *testing.T) {
	p := &proxy{
		allowedActions:   map[string]struct{}{"generatedatakey": {}},
		allowedKeyARNs:   map[string]struct{}{},
		enforceRecipient: true,
		maxBodyBytes:     1024,
	}

	req := httptest.NewRequest(http.MethodPost, "http://example.com/", strings.NewReader(`{"KeyId":"arn:ok"}`))
	req.Header.Set("X-Amz-Target", "KMS.GenerateDataKey")
	rr := httptest.NewRecorder()

	p.ServeHTTP(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}

func TestServeHTTP_KeyNotAllowed(t *testing.T) {
	p := &proxy{
		allowedActions:   map[string]struct{}{"generatedatakey": {}},
		allowedKeyARNs:   map[string]struct{}{"arn:allowed": {}},
		enforceRecipient: true,
		maxBodyBytes:     1024,
	}

	body := `{"KeyId":"arn:denied","Recipient":{"AttestationDocument":"x","KeyEncryptionAlgorithm":"RSAES_OAEP_SHA_256"}}`
	req := httptest.NewRequest(http.MethodPost, "http://example.com/", strings.NewReader(body))
	req.Header.Set("X-Amz-Target", "KMS.GenerateDataKey")
	rr := httptest.NewRecorder()

	p.ServeHTTP(rr, req)

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", rr.Code)
	}
}
