package nsmutil

import (
	"fmt"

	"github.com/hf/nsm/request"
	"github.com/hf/nsm/response"
)

// Session abstracts nsm.Session for testability.
type Session interface {
	Send(req request.Request) (response.Response, error)
	Close() error
}

// SessionOpener creates a new NSM session.
type SessionOpener func() (Session, error)

// AttestWithSession sends an attestation request using the provided session opener.
// It returns the attestation document on success.
func AttestWithSession(opener SessionOpener, publicKey, userData []byte) ([]byte, error) {
	session, err := opener()
	if err != nil {
		return nil, fmt.Errorf("failed to open NSM session (not running in enclave?): %w", err)
	}
	defer session.Close()

	res, err := session.Send(&request.Attestation{
		PublicKey: publicKey,
		UserData:  userData,
	})
	if err != nil {
		return nil, fmt.Errorf("NSM attestation request failed: %w", err)
	}
	if res.Error != "" {
		return nil, fmt.Errorf("NSM device returned an error: %s", res.Error)
	}
	if res.Attestation == nil || res.Attestation.Document == nil {
		return nil, fmt.Errorf("NSM returned empty attestation document")
	}

	return res.Attestation.Document, nil
}
