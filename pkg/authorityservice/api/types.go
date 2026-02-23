package api

// NonceResponse is returned by GET /nonce.
type NonceResponse struct {
	Salt      string `json:"salt"`
	Nonce     string `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
}

// GetReportRequest is the payload for POST /getreport.
type GetReportRequest struct {
	ChainID   uint64 `json:"chain_id"`
	AppID     uint64 `json:"app_id"`
	ReportID  string `json:"report_id"`
	Salt      string `json:"salt"`
	Nonce     string `json:"nonce"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

// GetReportResponse is returned by POST /getreport.
type GetReportResponse struct {
	ApplicationID   string `json:"applicationId"`
	ReportID        string `json:"reportId"`
	Authority       string `json:"authority"`
	EncryptedReport string `json:"encryptedReport"`
}
