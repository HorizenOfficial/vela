// Package kms provides AWS KMS integration for Nitro Enclaves key management.
package kms

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
	"github.com/mdlayher/vsock"
)

const (
	// ParentCID is the vsock CID for the parent EC2 instance.
	// In Nitro Enclaves, CID 3 is reserved for the parent.
	ParentCID = 3
)

// NitroKMSClient implements KMSClient using AWS KMS with Nitro Enclave attestation support.
// It uses the AWS SDK v2 with the Recipient parameter for attestation-based key operations.
// Inside a Nitro Enclave, it connects to KMS through a vsock proxy running on the parent EC2.
type NitroKMSClient struct {
	client *kms.Client
	keyARN string
}

// NewNitroKMSClient creates a new NitroKMSClient configured for the specified region.
// Inside a Nitro Enclave, it uses vsock to communicate with the kms-proxy on the parent EC2.
// Requests are sent unsigned (anonymous credentials); the proxy signs them with IMDS creds.
//
// Parameters:
//   - ctx: Context for the operation
//   - region: AWS region (defaults to "us-east-1")
//   - keyARN: ARN of the KMS key to use
//   - proxyPort: vsock port where the kms-proxy listens on the parent EC2 (typically 8000)
func NewNitroKMSClient(ctx context.Context, region, keyARN string, proxyPort uint32) (*NitroKMSClient, error) {
	if region == "" {
		region = "us-east-1"
	}
	if keyARN == "" {
		return nil, fmt.Errorf("KMS key ARN is required")
	}

	proxyEndpoint := fmt.Sprintf("http://localhost:%d", proxyPort)

	// Create a custom HTTP client that uses vsock to connect to the kms-proxy
	httpClient := &http.Client{
		Transport: &http.Transport{
			DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
				// Connect to the kms-proxy on the parent EC2 via vsock
				return vsock.Dial(ParentCID, proxyPort, nil)
			},
		},
	}

	// Load AWS config with the custom HTTP client
	// Note: requests must be unsigned inside the enclave; the proxy will sign them.
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithHTTPClient(httpClient),
		config.WithCredentialsProvider(aws.AnonymousCredentials{}),
		config.WithEndpointResolverWithOptions(
			aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               proxyEndpoint,
					SigningRegion:     region,
					HostnameImmutable: true,
				}, nil
			}),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	return &NitroKMSClient{
		client: kms.NewFromConfig(cfg),
		keyARN: keyARN,
	}, nil
}

// GenerateDataKeyWithAttestation generates a new data key using KMS with Nitro Enclave attestation.
// The attestation document is a CBOR-encoded Nitro Enclave attestation that includes:
// - PCR values (especially PCR0 for enclave image identity)
// - A public key to encrypt the response
//
// KMS validates the attestation against the key policy's kms:RecipientAttestation conditions
// before returning the encrypted data key.
//
// Returns:
// - CiphertextBlob: The data key encrypted with the KMS key (for storage/recovery)
// - CiphertextForRecipient: The data key encrypted with the enclave's public key (for immediate use)
func (c *NitroKMSClient) GenerateDataKeyWithAttestation(
	ctx context.Context,
	keyARN string,
	attestationDoc []byte,
) (*DataKeyOutput, error) {
	if keyARN == "" {
		keyARN = c.keyARN
	}
	if attestationDoc == nil || len(attestationDoc) == 0 {
		return nil, fmt.Errorf("attestation document is required for Nitro Enclave KMS operations")
	}

	input := &kms.GenerateDataKeyInput{
		KeyId:   aws.String(keyARN),
		KeySpec: types.DataKeySpecAes256,
		Recipient: &types.RecipientInfo{
			AttestationDocument:    attestationDoc,
			KeyEncryptionAlgorithm: types.KeyEncryptionMechanismRsaesOaepSha256,
		},
	}

	output, err := c.client.GenerateDataKey(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("KMS GenerateDataKey failed: %w", err)
	}
	fmt.Printf("kms generate data key response: ciphertext_blob_len=%d ciphertext_for_recipient_len=%d\n", len(output.CiphertextBlob), len(output.CiphertextForRecipient))

	// When using Recipient parameter:
	// - Plaintext is nil (not returned for security)
	// - CiphertextBlob contains the data key encrypted with KMS key
	// - CiphertextForRecipient contains the data key encrypted with the enclave's public key
	if output.CiphertextForRecipient == nil {
		return nil, fmt.Errorf("KMS did not return CiphertextForRecipient - attestation may have failed")
	}

	return &DataKeyOutput{
		CiphertextBlob:         output.CiphertextBlob,
		CiphertextForRecipient: output.CiphertextForRecipient,
	}, nil
}

// DecryptWithAttestation decrypts a ciphertext blob using KMS with Nitro Enclave attestation.
// This is used during key restoration to decrypt the stored CiphertextBlob.
//
// The attestation document must be from an enclave that matches the KMS key policy's
// PCR conditions. This ensures only authorized enclave images can decrypt the data.
//
// Returns the plaintext encrypted for the enclave's public key (CiphertextForRecipient),
// which must be decrypted using the enclave's private RSA key.
func (c *NitroKMSClient) DecryptWithAttestation(
	ctx context.Context,
	ciphertext []byte,
	attestationDoc []byte,
) ([]byte, error) {
	if ciphertext == nil || len(ciphertext) == 0 {
		return nil, fmt.Errorf("ciphertext is required")
	}
	if attestationDoc == nil || len(attestationDoc) == 0 {
		return nil, fmt.Errorf("attestation document is required for Nitro Enclave KMS operations")
	}

	input := &kms.DecryptInput{
		CiphertextBlob: ciphertext,
		Recipient: &types.RecipientInfo{
			AttestationDocument:    attestationDoc,
			KeyEncryptionAlgorithm: types.KeyEncryptionMechanismRsaesOaepSha256,
		},
	}

	output, err := c.client.Decrypt(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("KMS Decrypt failed: %w", err)
	}

	// When using Recipient parameter:
	// - Plaintext is nil (not returned for security)
	// - CiphertextForRecipient contains the plaintext encrypted with the enclave's public key
	if output.CiphertextForRecipient == nil {
		return nil, fmt.Errorf("KMS did not return CiphertextForRecipient - attestation may have failed")
	}

	return output.CiphertextForRecipient, nil
}
