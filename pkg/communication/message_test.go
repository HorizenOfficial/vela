package communication

import (
	"math/big"
	"testing"

	"github.com/HorizenOfficial/vela/pkg/common"
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
		DepositAmount:   common.NewBig(10),
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
					DepositAmount:   common.NewBig(10),
					MaxFeeValue:     common.NewBig(5),
				},
				ApplicationState: validApplicationState,
				WasmModule:       validWasmModule,
			},
			err: "cannot marshal negative Big value",
		},
		{
			name: "Valid Request - zero DepositAmount",
			data: ProcessRequestData{
				Request: &common.Request{
					ProtocolVersion: 0,
					ApplicationID:   common.NewApplicationId(1),
					RequestID:       common.RequestIdType([32]byte{1}),
					RequestType:     common.Process,
					Payload:         []byte("test"),
					Timestamp:       common.NewBig(100),
					Sender:          [20]byte{1},
					DepositAmount:   common.NewBig(0), // Zero deposit amount
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
					DepositAmount:   common.NewBig(10),
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
				require.Zero(t, tt.data.Request.DepositAmount.ToInt().Cmp(extractedData.Request.DepositAmount.ToInt()))
				require.Zero(t, tt.data.Request.MaxFeeValue.ToInt().Cmp(extractedData.Request.MaxFeeValue.ToInt()))
			}
		})
	}
}
