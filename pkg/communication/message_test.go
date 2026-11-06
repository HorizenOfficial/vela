package communication

import (
	"math/big"
	"testing"

	velacommon "github.com/HorizenOfficial/vela-common-go/common"
	"github.com/HorizenOfficial/vela/pkg/common"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
)

func TestProcessRequestDataValidate(t *testing.T) {
	// A valid request for testing
	validRequest := &common.Request{
		ProtocolVersion: 0,
		ApplicationID:   common.NewApplicationId(1),
		RequestID:       common.RequestIdType([32]byte{1}),
		RequestType:     common.Process,
		Payload:         []byte("test"),
		Timestamp:       common.NewBig(100),
		Sender:          [20]byte{1},
		TokenAddress:    velacommon.ETH_TOKEN,
		AssetAmount:     common.NewBig(10),
		MaxFeeValue:     common.NewBig(5),
	}

	validApplicationState := &common.ApplicationState{
		ApplicationID:  common.NewApplicationId(1),
		StateRoot:      [32]byte{1},
		EncryptedState: []byte("encrypted"),
	}

	validWasmModule := []byte("wasm")

	tests := []struct {
		name string
		data ProcessRequestData
		err  string // Expected error string, if any
	}{
		{
			name: "valid ProcessRequestData",
			data: ProcessRequestData{
				Request:          validRequest,
				ApplicationState: validApplicationState,
				WasmModule:       validWasmModule,
			},
			err: "",
		},
		{
			name: "missing Request",
			data: ProcessRequestData{
				Request:          nil,
				ApplicationState: validApplicationState,
				WasmModule:       validWasmModule,
			},
			err: "Request is required",
		},
		{
			name: "missing ApplicationState",
			data: ProcessRequestData{
				Request:          validRequest,
				ApplicationState: nil,
				WasmModule:       validWasmModule,
			},
			err: "",
		},
		{
			name: "empty WasmModule",
			data: ProcessRequestData{
				Request:          validRequest,
				ApplicationState: validApplicationState,
				WasmModule:       []byte{},
			},
			err: "",
		},
		{
			name: "invalid Request - negative Timestamp",
			data: ProcessRequestData{
				Request: &common.Request{
					ProtocolVersion: 0,
					ApplicationID:   common.NewApplicationId(1),
					RequestID:       common.RequestIdType([32]byte{1}),
					RequestType:     common.Process,
					Payload:         []byte("test"),
					Timestamp:       common.ToBig(big.NewInt(-1)), // Invalid timestamp
					Sender:          [20]byte{1},
					TokenAddress:    velacommon.ETH_TOKEN,
					AssetAmount:     common.NewBig(10),
					MaxFeeValue:     common.NewBig(5),
				},
				ApplicationState: validApplicationState,
				WasmModule:       validWasmModule,
			},
			err: "cannot marshal negative Big value",
		},
		{
			name: "Valid Request - zero AssetAmount",
			data: ProcessRequestData{
				Request: &common.Request{
					ProtocolVersion: 0,
					ApplicationID:   common.NewApplicationId(1),
					RequestID:       common.RequestIdType([32]byte{1}),
					RequestType:     common.Process,
					Payload:         []byte("test"),
					Timestamp:       common.NewBig(100),
					Sender:          [20]byte{1},
					TokenAddress:    velacommon.ETH_TOKEN,
					AssetAmount:     common.NewBig(0), // Zero asset amount
					MaxFeeValue:     common.NewBig(5),
				},
				ApplicationState: validApplicationState,
				WasmModule:       validWasmModule,
			},
			err: "",
		},
		{
			name: "Valid Request - zero MaxFeeValue",
			data: ProcessRequestData{
				Request: &common.Request{
					ProtocolVersion: 0,
					ApplicationID:   common.NewApplicationId(1),
					RequestID:       common.RequestIdType([32]byte{1}),
					RequestType:     common.Process,
					Payload:         []byte("test"),
					Timestamp:       common.NewBig(100),
					Sender:          [20]byte{1},
					TokenAddress:    velacommon.ETH_TOKEN,
					AssetAmount:     common.NewBig(10),
					MaxFeeValue:     common.NewBig(0), // zero max fee value
				},
				ApplicationState: validApplicationState,
				WasmModule:       validWasmModule,
			},
			err: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal and Unmarshal to simulate message transfer and trigger validation via extractData
			msg := Message{
				ID:   generateID(),
				Type: ProcessRequestMessage,
				Data: tt.data,
			}

			extractedData, err := extractData[ProcessRequestData](msg.Data)

			if tt.err != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
				require.Nil(t, extractedData)
			} else {
				require.NoError(t, err)
				require.NotNil(t, extractedData)
				// Further checks to ensure data integrity after extraction if needed
				require.Zero(t, tt.data.Request.Timestamp.ToInt().Cmp(extractedData.Request.Timestamp.ToInt()))
				require.Zero(t, tt.data.Request.AssetAmount.ToInt().Cmp(extractedData.Request.AssetAmount.ToInt()))
				require.Zero(t, tt.data.Request.MaxFeeValue.ToInt().Cmp(extractedData.Request.MaxFeeValue.ToInt()))
			}
		})
	}
}

func TestBatchProcessRequestDataValidate(t *testing.T) {
	validRequest := func() *common.Request {
		return &common.Request{
			ProtocolVersion: 0,
			ApplicationID:   common.NewApplicationId(1),
			RequestID:       common.RequestIdType([32]byte{1}),
			RequestType:     common.Process,
			Payload:         []byte("test"),
			Timestamp:       common.NewBig(100),
			Sender:          [20]byte{1},
			TokenAddress:    velacommon.ETH_TOKEN,
			AssetAmount:     common.NewBig(10),
			MaxFeeValue:     common.NewBig(5),
		}
	}

	// Marshals fine but fails Request.Validate(): zero assetAmount requires the native-token sentinel.
	invalidRequest := validRequest()
	invalidRequest.AssetAmount = common.NewBig(0)
	invalidRequest.TokenAddress = ethCommon.Address{2}

	requestForOtherApp := func() *common.Request {
		return &common.Request{
			ProtocolVersion: 0,
			ApplicationID:   common.NewApplicationId(12),
			RequestID:       common.RequestIdType([32]byte{1}),
			RequestType:     common.Process,
			Payload:         []byte("test"),
			Timestamp:       common.NewBig(100),
			Sender:          [20]byte{1},
			TokenAddress:    velacommon.ETH_TOKEN,
			AssetAmount:     common.NewBig(10),
			MaxFeeValue:     common.NewBig(5),
		}
	}

	validApplicationState := &common.ApplicationState{
		ApplicationID:  common.NewApplicationId(1),
		StateRoot:      [32]byte{1},
		EncryptedState: []byte("encrypted"),
	}

	tests := []struct {
		name string
		data BatchProcessRequestData
		err  string
	}{
		{
			name: "valid batch of multiple requests",
			data: BatchProcessRequestData{
				Requests:         []*common.Request{validRequest(), validRequest()},
				ApplicationState: validApplicationState,
				WasmModule:       []byte("wasm"),
			},
			err: "",
		},
		{
			name: "empty Requests",
			data: BatchProcessRequestData{
				Requests:         []*common.Request{},
				ApplicationState: validApplicationState,
				WasmModule:       []byte("wasm"),
			},
			err: "Requests is required",
		},
		{
			// The manager passes a nil state when it has none for the application;
			// this must be rejected, not dereferenced.
			name: "nil ApplicationState",
			data: BatchProcessRequestData{
				Requests:         []*common.Request{validRequest()},
				ApplicationState: nil,
				WasmModule:       []byte("wasm"),
			},
			err: "ApplicationState is required",
		},
		{
			name: "nil Requests",
			data: BatchProcessRequestData{
				Requests:         nil,
				ApplicationState: validApplicationState,
				WasmModule:       []byte("wasm"),
			},
			err: "Requests is required",
		},
		{
			name: "nil request in slice",
			data: BatchProcessRequestData{
				Requests:         []*common.Request{validRequest(), nil},
				ApplicationState: validApplicationState,
				WasmModule:       []byte("wasm"),
			},
			err: "Requests[1] is required",
		},
		{
			name: "invalid request in slice - zero AssetAmount with non-native token",
			data: BatchProcessRequestData{
				Requests:         []*common.Request{invalidRequest},
				ApplicationState: validApplicationState,
				WasmModule:       []byte("wasm"),
			},
			err: "invalid Requests[0]",
		},
		{
			name: "requests with different ApplicationIDs",
			data: BatchProcessRequestData{
				Requests:         []*common.Request{validRequest(), requestForOtherApp()},
				ApplicationState: validApplicationState,
				WasmModule:       []byte("wasm"),
			},
			err: "invalid Requests[1]",
		},
		{
			name: "request has different ApplicationIDs than state",
			data: BatchProcessRequestData{
				Requests:         []*common.Request{requestForOtherApp()},
				ApplicationState: validApplicationState,
				WasmModule:       []byte("wasm"),
			},
			err: "invalid Requests[0]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := Message{
				ID:   generateID(),
				Type: BatchProcessRequestMessage,
				Data: tt.data,
			}

			extractedData, err := extractData[BatchProcessRequestData](msg.Data)

			if tt.err != "" {
				require.Error(t, err)
				require.Contains(t, err.Error(), tt.err)
				require.Nil(t, extractedData)
			} else {
				require.NoError(t, err)
				require.NotNil(t, extractedData)
				require.Len(t, extractedData.Requests, len(tt.data.Requests))
			}
		})
	}
}
