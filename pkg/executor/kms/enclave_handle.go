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
	"os"
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

var kmsDebug = strings.EqualFold(os.Getenv("EXECUTOR_KMS_DEBUG"), "true") || os.Getenv("EXECUTOR_KMS_DEBUG") == "1"

func debugf(format string, args ...interface{}) {
	if !kmsDebug {
		return
	}
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	fmt.Printf(format, args...)
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
	debugf("kms attestation rsa key size=%d", rsaKey.Size())

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
	debugf("kms attestation public key: der_len=%d sha256=%x", len(publicKeyDER), publicKeyHash)

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
			debugf("kms enveloped key unwrap: source=%s encrypted_len=%d", source, len(unwrapped))
			envelopedKey = unwrapped
		} else {
			prefixLen := 8
			if len(envelopedKey) < prefixLen {
				prefixLen = len(envelopedKey)
			}
			debugf("kms enveloped key unwrap: cms=false ciphertext_len=%d prefix=%x", len(envelopedKey), envelopedKey[:prefixLen])
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

const (
	cmsEnvelopedDataOID = "1.2.840.113549.1.7.3"
	rsaesOAEPoid        = "1.2.840.113549.1.1.7"
)

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
	if err != nil || packet == nil || packet.Tag != ber.TagSequence || len(packet.Children) < 2 {
		return nil, false
	}

	if packet.Children[0].Tag != ber.TagObjectIdentifier {
		return nil, false
	}
	oid, err := decodeOID(packet.Children[0].ByteValue)
	if err != nil || oid != cmsEnvelopedDataOID {
		return nil, false
	}

	content := findCMSContent(packet.Children[1:])
	enveloped := extractSequence(content)
	recipientInfos := findRecipientInfos(enveloped)
	encryptedKey := findEncryptedKeyFromRecipientInfos(recipientInfos, expectedLen)
	if len(encryptedKey) != expectedLen {
		return nil, false
	}

	return encryptedKey, true
}

func findCMSContent(children []*ber.Packet) *ber.Packet {
	for _, child := range children {
		if child.ClassType == ber.ClassContext && child.Tag == 0 {
			if len(child.Children) == 1 {
				return child.Children[0]
			}
			return child
		}
	}
	return nil
}

func extractSequence(packet *ber.Packet) *ber.Packet {
	if packet == nil {
		return nil
	}
	if packet.Tag == ber.TagSequence {
		return packet
	}
	if len(packet.Children) == 1 && packet.Children[0].Tag == ber.TagSequence {
		return packet.Children[0]
	}
	return nil
}

func findRecipientInfos(enveloped *ber.Packet) *ber.Packet {
	if enveloped == nil || enveloped.Tag != ber.TagSequence || len(enveloped.Children) < 2 {
		return nil
	}

	idx := 0
	if enveloped.Children[idx].Tag != ber.TagInteger {
		return nil
	}
	idx++

	if idx < len(enveloped.Children) && enveloped.Children[idx].ClassType == ber.ClassContext && enveloped.Children[idx].Tag == 0 {
		idx++
	}

	if idx < len(enveloped.Children) {
		if enveloped.Children[idx].Tag == ber.TagSet || enveloped.Children[idx].Tag == ber.TagSequence {
			return enveloped.Children[idx]
		}
	}

	for _, child := range enveloped.Children {
		if child.Tag == ber.TagSet {
			return child
		}
	}

	return nil
}

func findEncryptedKeyFromRecipientInfos(recipientInfos *ber.Packet, expectedLen int) []byte {
	if recipientInfos == nil {
		return nil
	}

	for _, recipient := range recipientInfos.Children {
		if recipient.Tag != ber.TagSequence {
			continue
		}

		algOID := extractKeyEncryptionAlgorithmOID(recipient)
		if algOID != "" && algOID != rsaesOAEPoid {
			continue
		}

		if encryptedKey := extractEncryptedKey(recipient, expectedLen); len(encryptedKey) == expectedLen {
			return encryptedKey
		}
	}

	return nil
}

func extractKeyEncryptionAlgorithmOID(recipient *ber.Packet) string {
	if recipient == nil || recipient.Tag != ber.TagSequence || len(recipient.Children) < 3 {
		return ""
	}

	algSeq := extractSequence(recipient.Children[2])
	if algSeq == nil || len(algSeq.Children) == 0 || algSeq.Children[0].Tag != ber.TagObjectIdentifier {
		return ""
	}

	oid, err := decodeOID(algSeq.Children[0].ByteValue)
	if err != nil {
		return ""
	}

	return oid
}

func extractEncryptedKey(recipient *ber.Packet, expectedLen int) []byte {
	if recipient == nil || recipient.Tag != ber.TagSequence || len(recipient.Children) == 0 {
		return nil
	}

	last := recipient.Children[len(recipient.Children)-1]
	if last.Tag == ber.TagOctetString && len(last.ByteValue) == expectedLen {
		return last.ByteValue
	}

	return findBEROctetStringByLen(recipient, expectedLen)
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
			debugf("kms enveloped key unwrap: key_alg=%s", alg.Algorithm.String())
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
	if !kmsDebug {
		return
	}
	var cose []interface{}
	if err := cbor.Unmarshal(attestationDoc, &cose); err != nil {
		debugf("kms attestation doc decode error: %v", err)
		return
	}
	if len(cose) < 4 {
		debugf("kms attestation doc decode error: unexpected cose length=%d", len(cose))
		return
	}

	payload, ok := cose[2].([]byte)
	if !ok {
		debugf("kms attestation doc decode error: payload has type %T", cose[2])
		return
	}

	var att map[string]interface{}
	if err := cbor.Unmarshal(payload, &att); err != nil {
		debugf("kms attestation doc payload decode error: %v", err)
		return
	}

	pkRaw, ok := att["public_key"].([]byte)
	if !ok {
		debugf("kms attestation doc decode error: public_key has type %T", att["public_key"])
		return
	}

	pkHash := sha256.Sum256(pkRaw)
	match := bytes.Equal(pkRaw, expectedPublicKey)
	debugf("kms attestation doc public key: len=%d sha256=%x match=%t", len(pkRaw), pkHash, match)

	parsedKey, err := x509.ParsePKIXPublicKey(pkRaw)
	if err != nil {
		debugf("kms attestation doc public key parse error: %v", err)
		return
	}
	if rsaKey, ok := parsedKey.(*rsa.PublicKey); ok {
		debugf("kms attestation doc public key parsed: rsa_size=%d", rsaKey.Size())
		return
	}
	debugf("kms attestation doc public key parsed: type=%T", parsedKey)
}
