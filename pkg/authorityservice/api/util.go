package api

import (
	"encoding/binary"
	"encoding/hex"
	"errors"

	"github.com/horizen-pes/pkg/common"
)

// BuildMessage builds the raw message to be signed for /getreport.
func BuildMessage(chainID uint64, appID common.ApplicationIdType, reportID common.RequestIdType, nonce []byte) []byte {
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

// ParseRequestID converts the hex report_id string into the fixed-size type.
func ParseRequestID(id string) (common.RequestIdType, error) {
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
