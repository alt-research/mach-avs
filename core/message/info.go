package message

import "github.com/ethereum/go-ethereum/common/hexutil"

type BlockWorkProof struct {
	BlockHash   hexutil.Bytes `json:"block_hash"`
	BlockNumber uint64        `json:"block_number"`
}
