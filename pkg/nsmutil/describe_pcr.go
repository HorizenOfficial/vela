package nsmutil

import (
	"fmt"

	"github.com/hf/nsm/request"
)

// DescribePCRWithSession queries the NSM for the value of the PCR at the given
// index using the provided session opener. It returns the raw PCR bytes on
// success (48 bytes for a SHA-384 measurement such as PCR0).
func DescribePCRWithSession(opener SessionOpener, index uint16) ([]byte, error) {
	session, err := opener()
	if err != nil {
		return nil, fmt.Errorf("failed to open NSM session (not running in enclave?): %w", err)
	}
	defer session.Close()

	res, err := session.Send(&request.DescribePCR{Index: index})
	if err != nil {
		return nil, fmt.Errorf("NSM DescribePCR request failed: %w", err)
	}
	if res.Error != "" {
		return nil, fmt.Errorf("NSM device returned an error: %s", res.Error)
	}
	if res.DescribePCR == nil || res.DescribePCR.Data == nil {
		return nil, fmt.Errorf("NSM returned empty DescribePCR response")
	}

	return res.DescribePCR.Data, nil
}
