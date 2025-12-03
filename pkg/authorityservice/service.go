package authorityservice

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	ethCommon "github.com/ethereum/go-ethereum/common"
	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/horizen-pes/pkg/common"
	"github.com/horizen-pes/pkg/storage"
)

// AuthorityService exposes HTTP endpoints for authorities to fetch reports.
type AuthorityService struct {
	dataLayer storage.DataLayer
	secret    []byte
	nonceTTL  time.Duration
	chainID   uint64
	clock     func() time.Time
}

// NewAuthorityService builds a new service instance.
func NewAuthorityService(dl storage.DataLayer, chainID uint64, nonceTTL time.Duration) (*AuthorityService, error) {
	secret := make([]byte, 32)
	if _, err := rand.Read(secret); err != nil {
		return nil, fmt.Errorf("failed to generate HMAC secret: %w", err)
	}

	return &AuthorityService{
		dataLayer: dl,
		secret:    secret,
		nonceTTL:  nonceTTL,
		chainID:   chainID,
		clock:     time.Now,
	}, nil
}

// Handler returns an http.Handler with registered routes.
func (s *AuthorityService) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/nonce", s.handleNonce)
	mux.HandleFunc("/getreport", s.handleGetReport)
	return mux
}

type nonceResponse struct {
	Salt      string `json:"salt"`
	Nonce     string `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
}

func (s *AuthorityService) handleNonce(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		http.Error(w, "failed to generate nonce", http.StatusInternalServerError)
		return
	}

	ts := s.clock().Unix()
	nonceBytes := s.computeNonce(salt, ts)

	resp := nonceResponse{
		Salt:      hex.EncodeToString(salt),
		Nonce:     hex.EncodeToString(nonceBytes),
		Timestamp: ts,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to write nonce response: %v", err)
	}
}

type getReportRequest struct {
	ChainID   uint64 `json:"chain_id"`
	AppID     uint64 `json:"app_id"`
	ReportID  string `json:"report_id"`
	Salt      string `json:"salt"`
	Nonce     string `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

type getReportResponse struct {
	ApplicationID   string `json:"applicationId"`
	ReportID        string `json:"reportId"`
	Authority       string `json:"authority"`
	EncryptedReport string `json:"encryptedReport"`
	RefundAmount    string `json:"refundAmount,omitempty"`
	ApplicationFee  string `json:"applicationFee,omitempty"`
}

func (s *AuthorityService) handleGetReport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req getReportRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid json body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	if s.chainID != 0 && req.ChainID != s.chainID {
		http.Error(w, "unexpected chain_id", http.StatusBadRequest)
		return
	}

	reportID, err := parseRequestID(req.ReportID)
	if err != nil {
		http.Error(w, "invalid report_id", http.StatusBadRequest)
		return
	}

	appID := common.ApplicationIdType(req.AppID)

	if err := s.validateNonce(req.Salt, req.Nonce, req.Timestamp); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	signerAddr, err := s.recoverSigner(req, reportID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	report, err := s.dataLayer.GetDeanonymizationReport(ctx, reportID)
	if err != nil {
		http.Error(w, "report not found", http.StatusNotFound)
		return
	}

	if report.ApplicationID != appID {
		http.Error(w, "applicationId mismatch", http.StatusForbidden)
		return
	}

	if report.Authority != signerAddr {
		http.Error(w, "authority mismatch", http.StatusForbidden)
		return
	}

	resp := getReportResponse{
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
		log.Printf("Failed to write report response: %v", err)
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

func (s *AuthorityService) recoverSigner(req getReportRequest, reportID common.RequestIdType) (ethCommon.Address, error) {
	nonceBytes, err := hex.DecodeString(req.Nonce)
	if err != nil {
		return ethCommon.Address{}, fmt.Errorf("invalid nonce")
	}

	msg := buildMessage(req.ChainID, common.ApplicationIdType(req.AppID), reportID, nonceBytes)
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

func buildMessage(chainID uint64, appID common.ApplicationIdType, reportID common.RequestIdType, nonce []byte) []byte {
	buf := make([]byte, 0, 8+8+len(reportID)+len(nonce))

	var tmp [8]byte
	binary.BigEndian.PutUint64(tmp[:], chainID)
	buf = append(buf, tmp[:]...)

	binary.BigEndian.PutUint64(tmp[:], uint64(appID))
	buf = append(buf, tmp[:]...)

	buf = append(buf, reportID[:]...)
	buf = append(buf, nonce...)

	return buf
}

func parseRequestID(id string) (common.RequestIdType, error) {
	var reportID common.RequestIdType
	bytes, err := hex.DecodeString(id)
	if err != nil {
		return reportID, err
	}
	if len(bytes) != len(reportID) {
		return reportID, errors.New("invalid length for report_id")
	}
	copy(reportID[:], bytes)
	return reportID, nil
}
