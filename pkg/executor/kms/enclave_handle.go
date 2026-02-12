// Package kms provides AWS KMS integration for Nitro Enclaves key management.
package kms

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/asn1"
	"fmt"
	"strconv"
	"strings"

	ber "github.com/go-asn1-ber/asn1-ber"
	"github.com/fxamacker/cbor/v2"
	"github.com/hf/nsm"
	"github.com/hf/nsm/request"
)

// NitroEnclaveHandle implements EnclaveHandle using the Nitro Secure Module (NSM).
// It provides attestation generation and KMS envelope decryption capabilities
// that only work inside an AWS Nitro Enclave.
type NitroEnclaveHandle struct {
	// rsaKey is the ephemeral RSA key used for KMS envelope encryption.
	// The public key is included in the attestation document.
	rsaKey *rsa.PrivateKey
}

// NewNitroEnclaveHandle creates a new NitroEnclaveHandle.
// It generates an ephemeral RSA-2048 key pair for KMS envelope encryption.
// The public key will be included in attestation documents sent to KMS.
func NewNitroEnclaveHandle() (*NitroEnclaveHandle, error) {
	// Generate ephemeral RSA key pair for KMS communication
	// KMS will encrypt the data key using this public key
	rsaKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA key pair: %w", err)
	}
	fmt.Printf("kms attestation rsa key size=%d\n", rsaKey.Size())

	return &NitroEnclaveHandle{
		rsaKey: rsaKey,
	}, nil
}

// Attest generates a Nitro Enclave attestation document.
// The attestation document is a CBOR-encoded structure that includes:
// - PCR values (PCR0 = enclave image hash, PCR1 = kernel hash, PCR2 = application hash)
// - The enclave's RSA public key (for KMS to encrypt the response)
// - Optional user data (can be used for nonce/challenge-response)
//
// This function only works inside a Nitro Enclave. Outside an enclave,
// it will fail with "NSM not available".
func (h *NitroEnclaveHandle) Attest(userData []byte) ([]byte, error) {
	// Open connection to NSM (Nitro Secure Module)
	session, err := nsm.OpenDefaultSession()
	if err != nil {
		return nil, fmt.Errorf("failed to open NSM session (not running in enclave?): %w", err)
	}
	defer session.Close()

	// Marshal the RSA public key in DER (PKIX/SPKI) format for inclusion in attestation
	publicKeyDER, err := x509.MarshalPKIXPublicKey(&h.rsaKey.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal RSA public key: %w", err)
	}
	publicKeyHash := sha256.Sum256(publicKeyDER)
	fmt.Printf("kms attestation public key: der_len=%d sha256=%x\n", len(publicKeyDER), publicKeyHash)

	// Request attestation from NSM
	// - UserData: optional application-specific data (e.g., nonce)
	// - PublicKey: RSA public key for KMS to encrypt the response
	// - Nonce: optional additional randomness
	res, err := session.Send(&request.Attestation{
		UserData:  userData,
		PublicKey: publicKeyDER,
		Nonce:     nil,
	})
	if err != nil {
		return nil, fmt.Errorf("NSM attestation request failed: %w", err)
	}

	if res.Attestation == nil || res.Attestation.Document == nil {
		return nil, fmt.Errorf("NSM returned empty attestation document")
	}
	logAttestationDocPublicKey(res.Attestation.Document, publicKeyDER)

	return res.Attestation.Document, nil
}

// DecryptKMSEnvelopedKey decrypts a KMS CiphertextForRecipient blob.
// KMS encrypts the data key using RSAES-OAEP with SHA-256 using the
// public key from the attestation document.
//
// This function uses the enclave's private RSA key to decrypt the envelope.
func (h *NitroEnclaveHandle) DecryptKMSEnvelopedKey(envelopedKey []byte) ([]byte, error) {
	if envelopedKey == nil || len(envelopedKey) == 0 {
		return nil, fmt.Errorf("enveloped key is empty")
	}

	if len(envelopedKey) != h.rsaKey.Size() {
		if unwrapped, source, ok := unwrapRecipientEncryptedKey(envelopedKey, h.rsaKey.Size()); ok {
			fmt.Printf("kms enveloped key unwrap: source=%s encrypted_len=%d\n", source, len(unwrapped))
			envelopedKey = unwrapped
		} else {
			prefixLen := 8
			if len(envelopedKey) < prefixLen {
				prefixLen = len(envelopedKey)
			}
			fmt.Printf("kms enveloped key unwrap: cms=false ciphertext_len=%d prefix=%x\n", len(envelopedKey), envelopedKey[:prefixLen])
		}
	}

	// KMS uses RSAES-OAEP with SHA-256 for envelope encryption
	plaintext, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		h.rsaKey,
		envelopedKey,
		nil, // No label
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to decrypt KMS envelope (ciphertext_len=%d key_size=%d): %w",
			len(envelopedKey),
			h.rsaKey.Size(),
			err,
		)
	}

	return plaintext, nil
}

// GetPublicKey returns the RSA public key in DER (PKIX/SPKI) format.
// This can be used for debugging or verification purposes.
func (h *NitroEnclaveHandle) GetPublicKey() ([]byte, error) {
	return x509.MarshalPKIXPublicKey(&h.rsaKey.PublicKey)
}

const cmsEnvelopedDataOID = "1.2.840.113549.1.7.3"

type keyTransRecipientInfo struct {
	Version                int
	RecipientID            asn1.RawValue
	KeyEncryptionAlgorithm asn1.RawValue
	EncryptedKey           []byte
}

type algorithmIdentifier struct {
	Algorithm  asn1.ObjectIdentifier
	Parameters asn1.RawValue `asn1:"optional"`
}

func unwrapRecipientEncryptedKey(envelopedKey []byte, expectedLen int) ([]byte, string, bool) {
	if encryptedKey, ok := unwrapCMSRecipientEncryptedKey(envelopedKey, expectedLen); ok {
		return encryptedKey, "cms", true
	}

	if encryptedKey, ok := unwrapKeyTransRecipientInfo(envelopedKey, expectedLen); ok {
		return encryptedKey, "keytrans", true
	}

	return nil, "", false
}

func unwrapCMSRecipientEncryptedKey(envelopedKey []byte, expectedLen int) ([]byte, bool) {
	packet, err := ber.DecodePacketErr(envelopedKey)
	if err != nil || packet == nil || packet.Tag != ber.TagSequence || len(packet.Children) == 0 {
		return nil, false
	}

	if packet.Children[0].Tag != ber.TagObjectIdentifier {
		return nil, false
	}
	oid, err := decodeOID(packet.Children[0].ByteValue)
	if err != nil || oid != cmsEnvelopedDataOID {
		return nil, false
	}

	encryptedKey := findBEROctetStringByLen(packet, expectedLen)
	if len(encryptedKey) != expectedLen {
		return nil, false
	}

	return encryptedKey, true
}

func findBEROctetStringByLen(packet *ber.Packet, expectedLen int) []byte {
	if packet == nil {
		return nil
	}
	if packet.Tag == ber.TagOctetString {
		if len(packet.ByteValue) == expectedLen {
			return packet.ByteValue
		}
	}
	for _, child := range packet.Children {
		if value := findBEROctetStringByLen(child, expectedLen); value != nil {
			return value
		}
	}
	return nil
}

func unwrapKeyTransRecipientInfo(envelopedKey []byte, expectedLen int) ([]byte, bool) {
	var info keyTransRecipientInfo
	if _, err := asn1.Unmarshal(envelopedKey, &info); err != nil {
		return nil, false
	}

	if len(info.EncryptedKey) != expectedLen {
		return nil, false
	}

	if len(info.KeyEncryptionAlgorithm.FullBytes) > 0 {
		var alg algorithmIdentifier
		if _, err := asn1.Unmarshal(info.KeyEncryptionAlgorithm.FullBytes, &alg); err == nil {
			fmt.Printf("kms enveloped key unwrap: key_alg=%s\n", alg.Algorithm.String())
		}
	}

	return info.EncryptedKey, true
}

func decodeOID(der []byte) (string, error) {
	if len(der) == 0 {
		return "", fmt.Errorf("empty OID")
	}

	first := int(der[0])
	arcs := []int{first / 40, first % 40}

	value := 0
	for i := 1; i < len(der); i++ {
		b := der[i]
		value = (value << 7) | int(b&0x7f)
		if b&0x80 == 0 {
			arcs = append(arcs, value)
			value = 0
		}
	}
	if value != 0 {
		return "", fmt.Errorf("incomplete OID")
	}

	var builder strings.Builder
	for i, arc := range arcs {
		if i > 0 {
			builder.WriteByte('.')
		}
		builder.WriteString(strconv.Itoa(arc))
	}
	return builder.String(), nil
}

func logAttestationDocPublicKey(attestationDoc, expectedPublicKey []byte) {
	var cose []interface{}
	if err := cbor.Unmarshal(attestationDoc, &cose); err != nil {
		fmt.Printf("kms attestation doc decode error: %v\n", err)
		return
	}
	if len(cose) < 4 {
		fmt.Printf("kms attestation doc decode error: unexpected cose length=%d\n", len(cose))
		return
	}

	payload, ok := cose[2].([]byte)
	if !ok {
		fmt.Printf("kms attestation doc decode error: payload has type %T\n", cose[2])
		return
	}

	var att map[string]interface{}
	if err := cbor.Unmarshal(payload, &att); err != nil {
		fmt.Printf("kms attestation doc payload decode error: %v\n", err)
		return
	}

	pkRaw, ok := att["public_key"].([]byte)
	if !ok {
		fmt.Printf("kms attestation doc decode error: public_key has type %T\n", att["public_key"])
		return
	}

	pkHash := sha256.Sum256(pkRaw)
	match := bytes.Equal(pkRaw, expectedPublicKey)
	fmt.Printf("kms attestation doc public key: len=%d sha256=%x match=%t\n", len(pkRaw), pkHash, match)

	parsedKey, err := x509.ParsePKIXPublicKey(pkRaw)
	if err != nil {
		fmt.Printf("kms attestation doc public key parse error: %v\n", err)
		return
	}
	if rsaKey, ok := parsedKey.(*rsa.PublicKey); ok {
		fmt.Printf("kms attestation doc public key parsed: rsa_size=%d\n", rsaKey.Size())
		return
	}
	fmt.Printf("kms attestation doc public key parsed: type=%T\n", parsedKey)
}
