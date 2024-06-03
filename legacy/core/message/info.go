package message

import "github.com/ethereum/go-ethereum/common/hexutil"

type BlockWorkProof struct {
	BlockHash   hexutil.Bytes `json:"block_hash"`
	BlockNumber hexutil.Big   `json:"block_number"`
}

type HealthCheckMsg struct {
	AvsName string         `json:"avs_name"`
	Method  string         `json:"method"`
	Proof   BlockWorkProof `json:"proof"`
}
