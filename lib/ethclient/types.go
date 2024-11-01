package ethclient

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
)

// ExecutionPayloadV2 includes only necessary fields of `engine.ExecutableData`.
// See: https://github.com/ethereum/execution-apis/blob/main/src/engine/shanghai.md#executionpayloadv2
type ExecutionPayloadV2 struct {
	ParentHash    common.Hash         `json:"parentHash"    gencodec:"required"`
	FeeRecipient  common.Address      `json:"feeRecipient"  gencodec:"required"`
	StateRoot     common.Hash         `json:"stateRoot"     gencodec:"required"`
	ReceiptsRoot  common.Hash         `json:"receiptsRoot"  gencodec:"required"`
	LogsBloom     []byte              `json:"logsBloom"     gencodec:"required"`
	Random        common.Hash         `json:"prevRandao"    gencodec:"required"`
	Number        uint64              `json:"blockNumber"   gencodec:"required"`
	GasLimit      uint64              `json:"gasLimit"      gencodec:"required"`
	GasUsed       uint64              `json:"gasUsed"       gencodec:"required"`
	Timestamp     uint64              `json:"timestamp"     gencodec:"required"`
	ExtraData     []byte              `json:"extraData"     gencodec:"required"`
	BaseFeePerGas *big.Int            `json:"baseFeePerGas" gencodec:"required"`
	BlockHash     common.Hash         `json:"blockHash"     gencodec:"required"`
	Transactions  [][]byte            `json:"transactions"  gencodec:"required"`
	Withdrawals   []*types.Withdrawal `json:"withdrawals"`
}

// MarshalJSON marshals as JSON.
func (e ExecutionPayloadV2) MarshalJSON() ([]byte, error) {
	type ExecutableData struct {
		ParentHash    common.Hash         `json:"parentHash"    gencodec:"required"`
		FeeRecipient  common.Address      `json:"feeRecipient"  gencodec:"required"`
		StateRoot     common.Hash         `json:"stateRoot"     gencodec:"required"`
		ReceiptsRoot  common.Hash         `json:"receiptsRoot"  gencodec:"required"`
		LogsBloom     hexutil.Bytes       `json:"logsBloom"     gencodec:"required"`
		Random        common.Hash         `json:"prevRandao"    gencodec:"required"`
		Number        hexutil.Uint64      `json:"blockNumber"   gencodec:"required"`
		GasLimit      hexutil.Uint64      `json:"gasLimit"      gencodec:"required"`
		GasUsed       hexutil.Uint64      `json:"gasUsed"       gencodec:"required"`
		Timestamp     hexutil.Uint64      `json:"timestamp"     gencodec:"required"`
		ExtraData     hexutil.Bytes       `json:"extraData"     gencodec:"required"`
		BaseFeePerGas *hexutil.Big        `json:"baseFeePerGas" gencodec:"required"`
		BlockHash     common.Hash         `json:"blockHash"     gencodec:"required"`
		Transactions  []hexutil.Bytes     `json:"transactions"  gencodec:"required"`
		Withdrawals   []*types.Withdrawal `json:"withdrawals"`
	}
	var enc ExecutableData
	enc.ParentHash = e.ParentHash
	enc.FeeRecipient = e.FeeRecipient
	enc.StateRoot = e.StateRoot
	enc.ReceiptsRoot = e.ReceiptsRoot
	enc.LogsBloom = e.LogsBloom
	enc.Random = e.Random
	enc.Number = hexutil.Uint64(e.Number)
	enc.GasLimit = hexutil.Uint64(e.GasLimit)
	enc.GasUsed = hexutil.Uint64(e.GasUsed)
	enc.Timestamp = hexutil.Uint64(e.Timestamp)
	enc.ExtraData = e.ExtraData
	enc.BaseFeePerGas = (*hexutil.Big)(e.BaseFeePerGas)
	enc.BlockHash = e.BlockHash
	if e.Transactions != nil {
		enc.Transactions = make([]hexutil.Bytes, len(e.Transactions))
		for k, v := range e.Transactions {
			enc.Transactions[k] = v
		}
	}
	enc.Withdrawals = e.Withdrawals
	return json.Marshal(&enc)
}

// ExecutionPayloadV3 includes only necessary fields of `engine.ExecutableData`.
// See: https://github.com/ethereum/execution-apis/blob/main/src/engine/cancun.md#executionpayloadv3
type ExecutionPayloadV3 struct {
	ParentHash    common.Hash         `json:"parentHash"    gencodec:"required"`
	FeeRecipient  common.Address      `json:"feeRecipient"  gencodec:"required"`
	StateRoot     common.Hash         `json:"stateRoot"     gencodec:"required"`
	ReceiptsRoot  common.Hash         `json:"receiptsRoot"  gencodec:"required"`
	LogsBloom     []byte              `json:"logsBloom"     gencodec:"required"`
	Random        common.Hash         `json:"prevRandao"    gencodec:"required"`
	Number        uint64              `json:"blockNumber"   gencodec:"required"`
	GasLimit      uint64              `json:"gasLimit"      gencodec:"required"`
	GasUsed       uint64              `json:"gasUsed"       gencodec:"required"`
	Timestamp     uint64              `json:"timestamp"     gencodec:"required"`
	ExtraData     []byte              `json:"extraData"     gencodec:"required"`
	BaseFeePerGas *big.Int            `json:"baseFeePerGas" gencodec:"required"`
	BlockHash     common.Hash         `json:"blockHash"     gencodec:"required"`
	Transactions  [][]byte            `json:"transactions"  gencodec:"required"`
	Withdrawals   []*types.Withdrawal `json:"withdrawals"`
	BlobGasUsed   *uint64             `json:"blobGasUsed"`
	ExcessBlobGas *uint64             `json:"excessBlobGas"`
}

// MarshalJSON marshals as JSON.
func (e ExecutionPayloadV3) MarshalJSON() ([]byte, error) {
	type ExecutableData struct {
		ParentHash    common.Hash         `json:"parentHash"    gencodec:"required"`
		FeeRecipient  common.Address      `json:"feeRecipient"  gencodec:"required"`
		StateRoot     common.Hash         `json:"stateRoot"     gencodec:"required"`
		ReceiptsRoot  common.Hash         `json:"receiptsRoot"  gencodec:"required"`
		LogsBloom     hexutil.Bytes       `json:"logsBloom"     gencodec:"required"`
		Random        common.Hash         `json:"prevRandao"    gencodec:"required"`
		Number        hexutil.Uint64      `json:"blockNumber"   gencodec:"required"`
		GasLimit      hexutil.Uint64      `json:"gasLimit"      gencodec:"required"`
		GasUsed       hexutil.Uint64      `json:"gasUsed"       gencodec:"required"`
		Timestamp     hexutil.Uint64      `json:"timestamp"     gencodec:"required"`
		ExtraData     hexutil.Bytes       `json:"extraData"     gencodec:"required"`
		BaseFeePerGas *hexutil.Big        `json:"baseFeePerGas" gencodec:"required"`
		BlockHash     common.Hash         `json:"blockHash"     gencodec:"required"`
		Transactions  []hexutil.Bytes     `json:"transactions"  gencodec:"required"`
		Withdrawals   []*types.Withdrawal `json:"withdrawals"`
		BlobGasUsed   *hexutil.Uint64     `json:"blobGasUsed"`
		ExcessBlobGas *hexutil.Uint64     `json:"excessBlobGas"`
	}
	var enc ExecutableData
	enc.ParentHash = e.ParentHash
	enc.FeeRecipient = e.FeeRecipient
	enc.StateRoot = e.StateRoot
	enc.ReceiptsRoot = e.ReceiptsRoot
	enc.LogsBloom = e.LogsBloom
	enc.Random = e.Random
	enc.Number = hexutil.Uint64(e.Number)
	enc.GasLimit = hexutil.Uint64(e.GasLimit)
	enc.GasUsed = hexutil.Uint64(e.GasUsed)
	enc.Timestamp = hexutil.Uint64(e.Timestamp)
	enc.ExtraData = e.ExtraData
	enc.BaseFeePerGas = (*hexutil.Big)(e.BaseFeePerGas)
	enc.BlockHash = e.BlockHash
	if e.Transactions != nil {
		enc.Transactions = make([]hexutil.Bytes, len(e.Transactions))
		for k, v := range e.Transactions {
			enc.Transactions[k] = v
		}
	}
	enc.Withdrawals = e.Withdrawals
	enc.BlobGasUsed = (*hexutil.Uint64)(e.BlobGasUsed)
	enc.ExcessBlobGas = (*hexutil.Uint64)(e.ExcessBlobGas)
	return json.Marshal(&enc)
}
