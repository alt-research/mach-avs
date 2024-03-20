package alert

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
	"golang.org/x/crypto/sha3"
)

type HexEncodedBytes32 [32]byte

func (b HexEncodedBytes32) MarshalJSON() ([]byte, error) {
	hexString := hex.EncodeToString(b[:])

	return json.Marshal(hexString)
}

func (b *HexEncodedBytes32) UnmarshalJSON(data []byte) (err error) {
	var hexString string

	if err = json.Unmarshal(data, &hexString); err != nil {
		return
	}

	var hexBytes []byte

	if hexBytes, err = hex.DecodeString(hexString); err != nil {
		return err
	}

	if len(hexBytes) != 32 {
		return fmt.Errorf("the bytes length not eq to 32")
	}

	copy(b[:], hexBytes[:32])

	return
}

// The Alert submit to avs
type Alert interface {
	MessageHash() [32]byte
	// Return a uint32 as task index in aggregator, it use a value for aggregator to index.
	TaskIndex() uint32
}

// The Alert Information
type AlertInfo struct {
	AlertHash [32]byte
	TaskIndex uint32
}

// AlertBlockMismatch is submit alert for verifier found a op block output mismatch.
//
//	It just a warning without any prove, the prover verifier should
//	submit a prove to ensure the alert is valid.
//	This alert can for the blocks which had not proposal its output
//	root to layer1, this block may not the checkpoint.
type AlertBlockMismatch struct {
	// The invalid output root verifier got from op-devnet.
	InvalidOutputRoot HexEncodedBytes32 `json:"invalid_output_root"`
	// The output root calc by verifier.
	ExpectOutputRoot HexEncodedBytes32 `json:"expect_output_root"`
	// The layer2 block 's number.
	L2BlockNumber *big.Int `json:"l2_block_number"`
}

// Return the message hash for signature in avs
func (a AlertBlockMismatch) MessageHash() [32]byte {
	var res [32]byte

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(a.InvalidOutputRoot[:])
	hasher.Write(a.ExpectOutputRoot[:])
	hasher.Write(a.L2BlockNumber.Bytes())
	copy(res[:], hasher.Sum(nil)[:32])

	return res
}

func (a AlertBlockMismatch) TaskIndex() uint32 {
	max := big.NewInt(math.MaxInt32)

	if a.L2BlockNumber.Cmp(max) != -1 {
		return uint32(a.L2BlockNumber.Uint64())
	} else {
		// TODO: support task index not use this func
		return uint32(a.L2BlockNumber.Mod(a.L2BlockNumber, max).Uint64())
	}
}

var _ Alert = (*AlertBlockMismatch)(nil)

//	AlertBlockOutputOracleMismatch is Submit alert for verifier found a op block output root mismatch.
//
// It just a warning without any prove, the prover verifier should
// submit a prove to ensure the alert is valid.
// This alert only for the porposaled output root by proposer,
// so we just sumit the index for this output root.
type AlertBlockOutputOracleMismatch struct {
	// The output root calc by verifier.
	ExpectOutputRoot HexEncodedBytes32 `json:"expect_output_root"`
	// The invalid output root index.
	InvalidOutputIndex *big.Int `json:"invalid_output_index"`
}

// Return the message hash for signature in avs
func (a AlertBlockOutputOracleMismatch) MessageHash() [32]byte {
	var res [32]byte

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(a.ExpectOutputRoot[:])
	hasher.Write(a.InvalidOutputIndex.Bytes())
	copy(res[:], hasher.Sum(nil)[:32])

	return res
}

func (a AlertBlockOutputOracleMismatch) TaskIndex() uint32 {
	max := big.NewInt(math.MaxInt32)

	if a.InvalidOutputIndex.Cmp(max) != -1 {
		return uint32(a.InvalidOutputIndex.Uint64())
	} else {
		// TODO: support task index not use this func
		return uint32(a.InvalidOutputIndex.Mod(a.InvalidOutputIndex, max).Uint64())
	}
}

var _ Alert = (*AlertBlockOutputOracleMismatch)(nil)
