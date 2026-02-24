package kms

import (
	"encoding/asn1"
	"strconv"
	"strings"
	"testing"

	ber "github.com/go-asn1-ber/asn1-ber"
)

func parseOIDString(value string) (asn1.ObjectIdentifier, error) {
	parts := strings.Split(value, ".")
	out := make(asn1.ObjectIdentifier, 0, len(parts))
	for _, part := range parts {
		n, err := strconv.Atoi(part)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, nil
}

func parseOID(t *testing.T, value string) asn1.ObjectIdentifier {
	t.Helper()
	out, err := parseOIDString(value)
	if err != nil {
		t.Fatalf("invalid OID %q: %v", value, err)
	}
	return out
}

func mustOID(value string) asn1.ObjectIdentifier {
	oid, err := parseOIDString(value)
	if err != nil {
		panic(err)
	}
	return oid
}

func oidValueBytes(t *testing.T, value string) []byte {
	t.Helper()
	der, err := asn1.Marshal(parseOID(t, value))
	if err != nil {
		t.Fatalf("marshal OID: %v", err)
	}
	var raw asn1.RawValue
	if _, err := asn1.Unmarshal(der, &raw); err != nil {
		t.Fatalf("unmarshal OID: %v", err)
	}
	return raw.Bytes
}

func oidPacket(t *testing.T, value string) *ber.Packet {
	t.Helper()
	p := ber.Encode(ber.ClassUniversal, ber.TypePrimitive, ber.TagObjectIdentifier, nil, "")
	p.ByteValue = oidValueBytes(t, value)
	p.Data.Write(p.ByteValue)
	return p
}

func mustOIDPacket(value string) *ber.Packet {
	oid, err := parseOIDString(value)
	if err != nil {
		panic(err)
	}
	der, err := asn1.Marshal(oid)
	if err != nil {
		panic(err)
	}
	var raw asn1.RawValue
	if _, err := asn1.Unmarshal(der, &raw); err != nil {
		panic(err)
	}
	p := ber.Encode(ber.ClassUniversal, ber.TypePrimitive, ber.TagObjectIdentifier, nil, "")
	p.ByteValue = raw.Bytes
	p.Data.Write(p.ByteValue)
	return p
}

func octetStringPacket(value []byte) *ber.Packet {
	return ber.Encode(ber.ClassUniversal, ber.TypePrimitive, ber.TagOctetString, string(value), "")
}

func buildCMS(encryptedKey []byte, oid string) []byte {
	top := ber.NewSequence("cms")
	top.AppendChild(mustOIDPacket(oid))
	content := ber.Encode(ber.ClassContext, ber.TypeConstructed, ber.Tag(0), nil, "")

	enveloped := ber.NewSequence("enveloped")
	enveloped.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagInteger, 0, ""))

	recipientInfos := ber.Encode(ber.ClassUniversal, ber.TypeConstructed, ber.TagSet, nil, "")
	recipient := ber.NewSequence("recipient")
	recipient.AppendChild(ber.NewInteger(ber.ClassUniversal, ber.TypePrimitive, ber.TagInteger, 0, ""))
	recipient.AppendChild(ber.NewSequence("rid"))

	algSeq := ber.NewSequence("alg")
	algSeq.AppendChild(mustOIDPacket(rsaesOAEPoid))
	recipient.AppendChild(algSeq)
	recipient.AppendChild(octetStringPacket(encryptedKey))

	recipientInfos.AppendChild(recipient)
	enveloped.AppendChild(recipientInfos)
	content.AppendChild(enveloped)
	top.AppendChild(content)

	return top.Bytes()
}

func buildKeyTrans(encryptedKey []byte) []byte {
	info := struct {
		Version                int
		RecipientID            asn1.RawValue
		KeyEncryptionAlgorithm algorithmIdentifier
		EncryptedKey           []byte
	}{
		Version:     0,
		RecipientID: asn1.RawValue{FullBytes: []byte{0x30, 0x00}},
		KeyEncryptionAlgorithm: algorithmIdentifier{
			Algorithm: mustOID(rsaesOAEPoid),
		},
		EncryptedKey: encryptedKey,
	}
	der, err := asn1.Marshal(info)
	if err != nil {
		panic(err)
	}
	return der
}

func TestDecodeOID(t *testing.T) {
	raw := oidValueBytes(t, rsaesOAEPoid)
	got, err := decodeOID(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != rsaesOAEPoid {
		t.Fatalf("expected %s, got %s", rsaesOAEPoid, got)
	}
}

func TestDecodeOIDInvalid(t *testing.T) {
	if _, err := decodeOID(nil); err == nil {
		t.Fatalf("expected error for empty OID")
	}
	if _, err := decodeOID([]byte{0x2a, 0x86}); err == nil {
		t.Fatalf("expected error for incomplete OID")
	}
}

func TestUnwrapCMSRecipientEncryptedKey(t *testing.T) {
	encrypted := []byte("0123456789abcdef")
	cms := buildCMS(encrypted, cmsEnvelopedDataOID)
	got, ok := unwrapCMSRecipientEncryptedKey(cms, len(encrypted))
	if !ok {
		t.Fatalf("expected unwrap ok")
	}
	if string(got) != string(encrypted) {
		t.Fatalf("unexpected encrypted key")
	}
}

func TestUnwrapCMSRecipientEncryptedKeyWrongOID(t *testing.T) {
	encrypted := []byte("0123456789abcdef")
	cms := buildCMS(encrypted, "1.2.3.4.5")
	if _, ok := unwrapCMSRecipientEncryptedKey(cms, len(encrypted)); ok {
		t.Fatalf("expected unwrap to fail for wrong OID")
	}
}

func TestUnwrapRecipientEncryptedKeySources(t *testing.T) {
	encrypted := []byte("0123456789abcdef")
	cms := buildCMS(encrypted, cmsEnvelopedDataOID)
	if got, source, ok := unwrapRecipientEncryptedKey(cms, len(encrypted)); !ok || source != "cms" || string(got) != string(encrypted) {
		t.Fatalf("expected cms source")
	}

	keyTrans := buildKeyTrans(encrypted)
	if got, source, ok := unwrapRecipientEncryptedKey(keyTrans, len(encrypted)); !ok || source != "keytrans" || string(got) != string(encrypted) {
		t.Fatalf("expected keytrans source")
	}
}

func TestUnwrapKeyTransRecipientInfo(t *testing.T) {
	encrypted := []byte("0123456789abcdef")
	keyTrans := buildKeyTrans(encrypted)
	got, ok := unwrapKeyTransRecipientInfo(keyTrans, len(encrypted))
	if !ok {
		t.Fatalf("expected unwrap ok")
	}
	if string(got) != string(encrypted) {
		t.Fatalf("unexpected encrypted key")
	}
}
