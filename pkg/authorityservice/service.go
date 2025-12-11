package authorityservice

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/horizen-pes/pkg/authorityservice/api"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/logger"
)

// AuthorityService exposes HTTP endpoints for authorities to fetch reports.
type AuthorityService struct {
	secret     []byte
	nonceTTL   time.Duration
	chainID    uint64
	reportPath string
	clock      func() time.Time
	log    logger.Logger
}

// NewAuthorityService builds a new service instance.
func NewAuthorityService(chainID uint64, nonceTTL time.Duration, reportPath string, log logger.Logger) (*AuthorityService, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, fmt.Errorf("failed to generate HMAC secret: %w", err)
	}

	return &AuthorityService{
		secret:     secret,
		nonceTTL:   nonceTTL,
		chainID:    chainID,
		reportPath: reportPath,
		clock:      time.Now,
		log:              log,
	}, nil
}

// Handler returns an http.Handler with registered routes.
func (s *AuthorityService) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/nonce", s.handleNonce)
	mux.HandleFunc("/getreport", s.handleGetReport)
	return mux
}

func (s *AuthorityService) handleNonce(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		http.Error(w, "failed to generate salt", http.StatusInternalServerError)
		return
	}

	ts := s.clock().Unix()
	nonceBytes := s.computeNonce(salt, ts)

	resp := api.NonceResponse{
		Salt:      hex.EncodeToString(salt),
		Nonce:     hex.EncodeToString(nonceBytes),
		Timestamp: ts,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		s.log.Error("Failed to write nonce response: %v", err)
		http.Error(w, "failed to write nonce response", http.StatusInternalServerError)
	}
}

func (s *AuthorityService) handleGetReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req api.GetReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.log.Error("getreport: invalid json body: %v", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if s.chainID == 0 {
		s.log.Error("getreport: authority service chain_id not configured")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if req.ChainID != s.chainID {
		s.log.Error("getreport: unexpected chain_id: got %d expected %d", req.ChainID, s.chainID)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	reportID, err := api.ParseRequestID(req.ReportID)
	if err != nil {
		s.log.Error("getreport: invalid report_id %q: %v", req.ReportID, err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	appID := common.ApplicationIdType(req.AppID)

	if err := s.validateNonce(req.Salt, req.Nonce, req.Timestamp); err != nil {
		s.log.Error("getreport: nonce validation failed for report %s app %d: %v", reportID.String(), appID, err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	report, err := s.loadReport(appID, reportID)
	if err != nil {
		s.log.Error("getreport: report %s not found: %v", reportID.String(), err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	if report.ApplicationID != appID {
		s.log.Error("getreport: applicationId mismatch for report %s: expected %s got %s", reportID.String(), report.ApplicationID.String(), appID.String())
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	signerAddr, err := s.recoverSigner(req, reportID)
	if err != nil {
		s.log.Error("getreport: failed to recover signer for report %s app %d: %v", reportID.String(), appID, err)
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if report.Authority != signerAddr {
		s.log.Error("getreport: authority mismatch for report %s: expected %s got %s", reportID.String(), report.Authority.Hex(), signerAddr.Hex())
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	resp := api.GetReportResponse{
		ApplicationID:   report.ApplicationID.String(),
		ReportID:        report.ReportID.String(),
		Authority:       report.Authority.Hex(),
		EncryptedReport: hex.EncodeToString(report.EncryptedReport),
	}

	if report.RefundAmount != nil {
		resp.RefundAmount = report.RefundAmount.String()
	}
	if report.ApplicationFee != nil {
		resp.ApplicationFee = report.ApplicationFee.String()
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		s.log.Error("Failed to write report response: %v", err)
		http.Error(w, "failed to write report response", http.StatusInternalServerError)
	}
}

func (s *AuthorityService) validateNonce(saltHex, nonceHex string, ts int64) error {
	salt, err := hex.DecodeString(saltHex)
	if err != nil {
		return fmt.Errorf("invalid salt")
	}

	providedNonce, err := hex.DecodeString(nonceHex)
	if err != nil {
		return fmt.Errorf("invalid nonce")
	}

	now := s.clock()
	tsTime := time.Unix(ts, 0)

	if tsTime.After(now.Add(1 * time.Minute)) {
		return fmt.Errorf("nonce timestamp in the future")
	}
	if now.Sub(tsTime) > s.nonceTTL {
		return fmt.Errorf("nonce expired")
	}

	expected := s.computeNonce(salt, ts)
	if !hmac.Equal(expected, providedNonce) {
		return fmt.Errorf("invalid nonce")
	}

	return nil
}

func (s *AuthorityService) computeNonce(salt []byte, ts int64) []byte {
	mac := hmac.New(sha256.New, s.secret)
	mac.Write(salt)

	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], uint64(ts))
	mac.Write(buf[:])

	return mac.Sum(nil)
}

func (s *AuthorityService) recoverSigner(req api.GetReportRequest, reportID common.RequestIdType) (ethCommon.Address, error) {
	nonceBytes, err := hex.DecodeString(req.Nonce)
	if err != nil {
		return ethCommon.Address{}, fmt.Errorf("invalid nonce")
	}

	msg := api.BuildMessage(req.ChainID, common.ApplicationIdType(req.AppID), reportID, nonceBytes)
	hash := ethCrypto.Keccak256Hash(msg)

	sigBytes, err := hex.DecodeString(req.Signature)
	if err != nil {
		return ethCommon.Address{}, fmt.Errorf("invalid signature encoding")
	}

	if len(sigBytes) != 65 {
		return ethCommon.Address{}, fmt.Errorf("signature must be 65 bytes")
	}

	// Normalize V to {0,1}
	if sigBytes[64] >= 27 {
		sigBytes[64] -= 27
	}

	pubKey, err := ethCrypto.SigToPub(hash.Bytes(), sigBytes)
	if err != nil {
		return ethCommon.Address{}, fmt.Errorf("failed to recover signer")
	}

	return ethCrypto.PubkeyToAddress(*pubKey), nil
}

func (s *AuthorityService) loadReport(appID common.ApplicationIdType, reportID common.RequestIdType) (*common.DeanonymizationReport, error) {
	// Try with the provided appID first
	paths := []string{filepath.Join(s.reportPath, common.ReportFilename(appID, reportID))}
	// Fallback: find by reportID regardless of appID to detect mismatches
	glob, _ := filepath.Glob(filepath.Join(s.reportPath, "*_"+reportID.String()))
	paths = append(paths, glob...)

	var lastErr error
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			lastErr = err
			continue
		}
		var report common.DeanonymizationReport
		if err := json.Unmarshal(data, &report); err != nil {
			lastErr = err
			continue
		}
		return &report, nil
	}
	if lastErr != nil {
		return nil, lastErr
	}
	return nil, fmt.Errorf("report not found")
}
