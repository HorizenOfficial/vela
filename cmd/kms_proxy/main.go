package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/mdlayher/vsock"
)

const (
	defaultPort       = uint32(8000)
	defaultRegion     = "us-east-1"
	defaultMaxBody    = 1024 * 1024
	defaultReadHeader = 5 * time.Second
	defaultRead       = 15 * time.Second
	defaultWrite      = 30 * time.Second
)

var hopByHopHeaders = map[string]struct{}{
	"Connection":          {},
	"Keep-Alive":          {},
	"Proxy-Authenticate":  {},
	"Proxy-Authorization": {},
	"TE":                  {},
	"Trailers":            {},
	"Transfer-Encoding":   {},
	"Upgrade":             {},
}

type proxy struct {
	region           string
	allowedActions   map[string]struct{}
	allowedKeyARNs   map[string]struct{}
	creds            *aws.CredentialsCache
	signer           *v4.Signer
	client           *http.Client
	enforceRecipient bool
	maxBodyBytes     int64
	debug            bool
}

func main() {
	ctx := context.Background()

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalf("kms proxy config error: %v", err)
	}

	awsCfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(cfg.region))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}

	p := &proxy{
		region:           cfg.region,
		allowedActions:   cfg.allowedActions,
		allowedKeyARNs:   cfg.allowedKeyARNs,
		creds:            aws.NewCredentialsCache(awsCfg.Credentials),
		signer:           v4.NewSigner(),
		client:           &http.Client{Timeout: cfg.upstreamTimeout},
		enforceRecipient: true,
		maxBodyBytes:     cfg.maxBodyBytes,
		debug:            cfg.debug,
	}

	listener, err := vsock.Listen(cfg.port, nil)
	if err != nil {
		log.Fatalf("failed to listen on vsock port %d: %v", cfg.port, err)
	}

	srv := &http.Server{
		Handler:           p,
		ReadHeaderTimeout: defaultReadHeader,
		ReadTimeout:       defaultRead,
		WriteTimeout:      defaultWrite,
	}

	log.Printf("kms proxy listening on vsock port %d (region=%s)", cfg.port, cfg.region)

	go func() {
		if err := srv.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("kms proxy server error: %v", err)
		}
	}()

	waitForSignal()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

func (p *proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	xTarget := r.Header.Get("X-Amz-Target")
	if xTarget == "" {
		http.Error(w, "missing X-Amz-Target header", http.StatusBadRequest)
		return
	}

	operation, err := parseOperation(xTarget)
	if err != nil {
		http.Error(w, "invalid X-Amz-Target header", http.StatusBadRequest)
		return
	}
	if !p.isActionAllowed(operation) {
		http.Error(w, "operation not allowed", http.StatusForbidden)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, p.maxBodyBytes))
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	payload, err := decodePayload(body)
	if err != nil {
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	if p.debug {
		keyID, _ := extractKeyID(payload)
		log.Printf("kms proxy request: action=%s key_id=%s", operation, keyID)
		logRecipientInfo(payload)
	}

	if p.enforceRecipient && !hasRecipient(payload) {
		http.Error(w, "missing Recipient attestation", http.StatusBadRequest)
		return
	}

	if keyID, ok := extractKeyID(payload); ok && len(p.allowedKeyARNs) > 0 {
		if !p.isKeyAllowed(keyID) {
			http.Error(w, "key not allowed", http.StatusForbidden)
			return
		}
	}

	kmsHost := fmt.Sprintf("kms.%s.amazonaws.com", p.region)
	kmsURL := fmt.Sprintf("https://%s%s", kmsHost, r.URL.RequestURI())
	outReq, err := http.NewRequestWithContext(r.Context(), r.Method, kmsURL, bytes.NewReader(body))
	if err != nil {
		http.Error(w, "failed to create upstream request", http.StatusBadRequest)
		return
	}

	outReq.Header = cloneHeaders(r.Header)
	stripHopHeaders(outReq.Header)
	outReq.Header.Del("Authorization")
	outReq.Header.Del("X-Amz-Date")
	outReq.Header.Del("X-Amz-Security-Token")
	outReq.Header.Del("X-Amz-Content-Sha256")
	outReq.Host = kmsHost
	outReq.ContentLength = int64(len(body))

	payloadHash := hashSHA256Hex(body)
	outReq.Header.Set("X-Amz-Content-Sha256", payloadHash)

	creds, err := p.creds.Retrieve(r.Context())
	if err != nil {
		http.Error(w, "failed to retrieve AWS credentials", http.StatusBadGateway)
		return
	}

	if err := p.signer.SignHTTP(r.Context(), creds, outReq, payloadHash, "kms", p.region, time.Now()); err != nil {
		http.Error(w, "failed to sign request", http.StatusBadGateway)
		return
	}

	resp, err := p.client.Do(outReq)
	if err != nil {
		http.Error(w, "failed to reach KMS", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read KMS response", http.StatusBadGateway)
		return
	}

	if p.debug {
		logRecipientResponseInfo(respBody, resp.StatusCode)
	}

	copyHeaders(w.Header(), resp.Header)
	w.WriteHeader(resp.StatusCode)
	_, _ = w.Write(respBody)
}

type proxyConfig struct {
	port           uint32
	region         string
	allowedActions map[string]struct{}
	allowedKeyARNs map[string]struct{}
	upstreamTimeout time.Duration
	maxBodyBytes   int64
	debug          bool
}

func loadConfig() (*proxyConfig, error) {
	port, err := getEnvUint32("KMS_PROXY_PORT", getEnvUint32Fallback("EXECUTOR_KMS_PROXY_PORT", defaultPort))
	if err != nil {
		return nil, fmt.Errorf("invalid KMS_PROXY_PORT: %w", err)
	}

	region := getEnvString("KMS_PROXY_REGION", "")
	if region == "" {
		region = getEnvString("AWS_REGION", "")
	}
	if region == "" {
		region = getEnvString("AWS_DEFAULT_REGION", "")
	}
	if region == "" {
		region = defaultRegion
	}

	allowedActions := parseCSVSetLower(getEnvString("KMS_PROXY_ALLOWED_ACTIONS", "GenerateDataKey,Decrypt"))
	if len(allowedActions) == 0 {
		return nil, fmt.Errorf("KMS_PROXY_ALLOWED_ACTIONS must include at least one action")
	}

	allowedKeyARNs := parseCSVSetExact(getEnvString("KMS_PROXY_ALLOWED_KEY_ARNS", ""))

	timeoutSec := getEnvInt("KMS_PROXY_UPSTREAM_TIMEOUT_SEC", 15)
	maxBody := getEnvInt("KMS_PROXY_MAX_BODY_BYTES", defaultMaxBody)
	debug := getEnvBool("KMS_PROXY_DEBUG", false)

	return &proxyConfig{
		port:            port,
		region:          region,
		allowedActions:  allowedActions,
		allowedKeyARNs:  allowedKeyARNs,
		upstreamTimeout: time.Duration(timeoutSec) * time.Second,
		maxBodyBytes:    int64(maxBody),
		debug:           debug,
	}, nil
}

func parseOperation(target string) (string, error) {
	parts := strings.Split(target, ".")
	if len(parts) != 2 || parts[1] == "" {
		return "", fmt.Errorf("invalid X-Amz-Target")
	}
	return strings.ToLower(parts[1]), nil
}

func decodePayload(body []byte) (map[string]any, error) {
	payload := make(map[string]any)
	if len(body) == 0 {
		return payload, nil
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, err
	}
	return payload, nil
}

func hasRecipient(payload map[string]any) bool {
	raw, ok := payload["Recipient"]
	if !ok || raw == nil {
		return false
	}
	recipient, ok := raw.(map[string]any)
	if !ok {
		return false
	}
	if recipient["AttestationDocument"] == nil {
		return false
	}
	if recipient["KeyEncryptionAlgorithm"] == nil {
		return false
	}
	return true
}

func extractKeyID(payload map[string]any) (string, bool) {
	raw, ok := payload["KeyId"]
	if !ok {
		return "", false
	}
	keyID, ok := raw.(string)
	if !ok || strings.TrimSpace(keyID) == "" {
		return "", false
	}
	return keyID, true
}

func (p *proxy) isActionAllowed(action string) bool {
	_, ok := p.allowedActions[action]
	return ok
}

func (p *proxy) isKeyAllowed(keyID string) bool {
	_, ok := p.allowedKeyARNs[keyID]
	return ok
}

func parseCSVSetLower(value string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, part := range strings.Split(value, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		set[strings.ToLower(trimmed)] = struct{}{}
	}
	return set
}

func parseCSVSetExact(value string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, part := range strings.Split(value, ",") {
		trimmed := strings.TrimSpace(part)
		if trimmed == "" {
			continue
		}
		set[trimmed] = struct{}{}
	}
	return set
}

func cloneHeaders(h http.Header) http.Header {
	clone := make(http.Header, len(h))
	for k, v := range h {
		values := make([]string, len(v))
		copy(values, v)
		clone[k] = values
	}
	return clone
}

func stripHopHeaders(h http.Header) {
	for key := range hopByHopHeaders {
		h.Del(key)
	}
}

func copyHeaders(dst, src http.Header) {
	for k, v := range src {
		values := make([]string, len(v))
		copy(values, v)
		dst[k] = values
	}
}

func logRecipientInfo(payload map[string]any) {
	raw, ok := payload["Recipient"]
	if !ok || raw == nil {
		log.Printf("kms proxy recipient: missing")
		return
	}
	recipient, ok := raw.(map[string]any)
	if !ok {
		log.Printf("kms proxy recipient: invalid type")
		return
	}

	alg, _ := recipient["KeyEncryptionAlgorithm"].(string)
	attestationLen := 0
	switch v := recipient["AttestationDocument"].(type) {
	case string:
		attestationLen = len(v)
	case []byte:
		attestationLen = len(v)
	}

	log.Printf("kms proxy recipient: alg=%s attestation_len=%d", alg, attestationLen)
}

func logRecipientResponseInfo(body []byte, status int) {
	if len(body) == 0 {
		log.Printf("kms proxy response: status=%d empty_body", status)
		return
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		log.Printf("kms proxy response: status=%d invalid_json", status)
		return
	}

	ciphertextLen := 0
	switch v := payload["CiphertextForRecipient"].(type) {
	case string:
		ciphertextLen = len(v)
	case []byte:
		ciphertextLen = len(v)
	}

	if ciphertextLen == 0 {
		log.Printf("kms proxy response: status=%d no_ciphertext_for_recipient", status)
		return
	}

	log.Printf("kms proxy response: status=%d ciphertext_for_recipient_len=%d", status, ciphertextLen)
}

func hashSHA256Hex(body []byte) string {
	sum := sha256.Sum256(body)
	return hex.EncodeToString(sum[:])
}

func getEnvString(key, fallback string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			return parsed
		}
	}
	return fallback
}

func getEnvUint32Fallback(key string, fallback uint32) uint32 {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		if parsed, err := strconv.ParseUint(val, 10, 32); err == nil {
			return uint32(parsed)
		}
	}
	return fallback
}

func getEnvUint32(key string, fallback uint32) (uint32, error) {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return fallback, nil
	}
	parsed, err := strconv.ParseUint(val, 10, 32)
	if err != nil {
		return 0, err
	}
	return uint32(parsed), nil
}

func getEnvBool(key string, fallback bool) bool {
	val := strings.TrimSpace(os.Getenv(key))
	if val == "" {
		return fallback
	}
	parsed, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}
	return parsed
}

func waitForSignal() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
}
