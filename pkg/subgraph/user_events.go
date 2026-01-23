package subgraph

import (
	"context"
	"fmt"
	"math/big"

	"github.com/horizen-pes/pkg/common"
	cryptotypes "github.com/horizen-pes/pkg/common/crypto"
	"github.com/horizen-pes/pkg/crypto"
)

var userEventsPageSize = 1000

// FetchAndDecryptUserEvents queries the subgraph for user events and decrypts them.
// The limit caps the number of decrypted events returned; limit <= 0 means no cap.
// The page size is internal and capped to avoid query errors.
// It applies the optional filter on decrypted payloads.
func FetchAndDecryptUserEvents(
	ctx context.Context,
	sg Client,
	teePubKey *cryptotypes.PublicKeyP521,
	privKey cryptotypes.PrivateKeyP521,
	applicationID common.ApplicationIdType,
	eventSubType string,
	limit int,
	filter func([]byte) bool,
) ([][]byte, error) {
	if sg == nil {
		return nil, fmt.Errorf("subgraph client is required")
	}
	if teePubKey == nil {
		return nil, fmt.Errorf("tee public key is required")
	}

	maxResults := limit
	if maxResults < 0 {
		maxResults = 0
	}
	pageSize := userEventsPageSize
	if pageSize <= 0 {
		pageSize = 1000
	}
	if pageSize > 1000 {
		pageSize = 1000
	}

	var decryptedEvents [][]byte
	var before *big.Int
	for {
		events, err := sg.GetUserEvents(ctx, applicationID, eventSubType, pageSize, before)
		if err != nil {
			return nil, err
		}
		if len(events) == 0 {
			break
		}

		for _, ev := range events {
			plain, err := crypto.Decrypt(teePubKey, &privKey, ev.EncryptedData)
			if err != nil {
				continue
			}
			if filter != nil && !filter(plain) {
				continue
			}
			decryptedEvents = append(decryptedEvents, plain)
			if maxResults > 0 && len(decryptedEvents) >= maxResults {
				return decryptedEvents, nil
			}
		}

		if len(events) < pageSize {
			break
		}
		before = userEventSortKey(events[len(events)-1])
	}

	return decryptedEvents, nil
}

// Must match SORT_BASE in the subgraph mapping.
const userEventSortKeyBase = uint64(1000000000)

func userEventSortKey(ev UserEvent) *big.Int {
	if ev.SortKey != nil {
		return new(big.Int).Set(ev.SortKey)
	}

	block := new(big.Int).SetUint64(ev.BlockNumber)
	base := new(big.Int).SetUint64(userEventSortKeyBase)
	block.Mul(block, base)
	block.Add(block, new(big.Int).SetUint64(ev.LogIndex))
	return block
}
