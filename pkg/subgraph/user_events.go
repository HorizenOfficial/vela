package subgraph

import (
	"context"
	"fmt"

	"github.com/horizen-pes/pkg/common"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/crypto"
)

// FetchAndDecryptUserEvents queries the subgraph for user events and decrypts them.
// It applies the optional filter on decrypted payloads and can stop at the first match.
func FetchAndDecryptUserEvents(
	ctx context.Context,
	sg Client,
	teePubKey *cryptotypes.PublicKeyP521,
	privKey cryptotypes.PrivateKeyP521,
	applicationID common.ApplicationIdType,
	eventSubType string,
	limit int,
	filter func([]byte) bool,
	stopAtFirst bool,
) ([][]byte, error) {
	if sg == nil {
		return nil, fmt.Errorf("subgraph client is required")
	}
	if teePubKey == nil {
		return nil, fmt.Errorf("tee public key is required")
	}
	if limit <= 0 {
		limit = 10
	}

	events, err := sg.GetUserEvents(ctx, applicationID, eventSubType, limit)
	if err != nil {
		return nil, err
	}

	var decryptedEvents [][]byte
	for _, ev := range events {
		plain, err := crypto.Decrypt(teePubKey, &privKey, ev.EncryptedData)
		if err != nil {
			continue
		}
		if filter != nil && !filter(plain) {
			continue
		}
		decryptedEvents = append(decryptedEvents, plain)
		if stopAtFirst && len(decryptedEvents) > 0 {
			return decryptedEvents, nil
		}
	}

	return decryptedEvents, nil
}
