package authorityservice

import (
	"context"
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
	"github.com/horizen-cce-common-go/wallet/subgraph"
	"github.com/horizen-pes/pkg/authorityservice/api"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/logger"
)

// AuthorityService exposes HTTP endpoints for authorities to fetch reports.
type AuthorityService struct {
	secret         []byte
	nonceTTL       time.Duration
	chainID        uint64
	reportPath     string
	subgraphClient subgraph.Client
	clock          func() time.Time
	log            logger.Logger
}

// NewAuthorityService builds a new service instance.
func NewAuthorityService(chainID uint64, nonceTTL time.Duration, reportPath string, sg subgraph.Client, log logger.Logger) (*AuthorityService, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, fmt.Errorf("failed to generate HMAC secret: %w", err)
	}

	if sg == nil {
		return nil, fmt.Errorf("subgraph client is required")
	}

	return &AuthorityService{
		secret:         secret,
		nonceTTL:       nonceTTL,
		chainID:        chainID,
		reportPath:     reportPath,
		subgraphClient: sg,
		clock:          time.Now,
		log:            log,
	}, nil
}

// Handler returns an http.Handler with registered routes.
func (s *AuthorityService) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/nonce", s.handleNonce)
	mux.HandleFunc("/getreport", s.handleGetReport)
	return mux
}

// findCompletionEvent queries the subgraph for the RequestCompleted event.
func (s *AuthorityService) findCompletionEvent(ctx context.Context, reportID common.RequestIdType) (*common.RequestResult, error) {
	event, err := s.subgraphClient.GetRequestCompletedByID(ctx, reportID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, nil
	}
	return &common.RequestResult{
		Status:       event.Status,
		ErrorCode:    event.ErrorCode,
		ErrorMessage: event.ErrorMessage,
	}, nil
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
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if report.Authority != signerAddr {
		s.log.Error("getreport: authority mismatch for report %s: expected %s got %s", reportID.String(), report.Authority.Hex(), signerAddr.Hex())
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	// Ensure the request was completed on-chain to avoid serving stale reports after reorgs.
	event, err := s.findCompletionEvent(r.Context(), reportID)
	if err != nil {
		s.log.Error("getreport: failed to query RequestCompleted for report %s: %v", reportID.String(), err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if event == nil {
		s.log.Error("getreport: no on-chain confirmation for report %s", reportID.String())
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	if event.Status != common.RequestResultOK {
		s.log.Error("getreport: on-chain status not OK for report %s: %v", reportID.String(), event.Status)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	resp := api.GetReportResponse{
		ApplicationID:   report.ApplicationID.String(),
		ReportID:        report.ReportID.String(),
		Authority:       report.Authority.Hex(),
		EncryptedReport: hex.EncodeToString(report.EncryptedReport),
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

	if tsTime.After(now) {
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
	path := filepath.Join(s.reportPath, common.ReportFilename(appID, reportID))

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("report not found for app %s id %s: %w", appID.String(), reportID.String(), err)
	}

	var report common.DeanonymizationReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("failed to decode report %s for app %s: %w", reportID.String(), appID.String(), err)
	}

	return &report, nil
}
