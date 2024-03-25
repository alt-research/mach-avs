package alert

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"golang.org/x/crypto/sha3"
)

type HexEncodedBytes32 [32]byte

func (b HexEncodedBytes32) MarshalJSON() ([]byte, error) {
	hexString := hex.EncodeToString(b[:])

	return json.Marshal(fmt.Sprintf("0x%s", hexString))
}

func (b *HexEncodedBytes32) UnmarshalJSON(data []byte) (err error) {
	var hexString string

	if err = json.Unmarshal(data, &hexString); err != nil {
		return err
	}

	var hexBytes []byte

	hexString = strings.TrimPrefix(hexString, "0x")
	hexString = strings.TrimPrefix(hexString, "0X")

	if hexBytes, err = hex.DecodeString(hexString); err != nil {
		return err
	}

	if len(hexBytes) != 32 {
		return fmt.Errorf("the bytes length not eq to 32")
	}

	copy(b[:], hexBytes[:32])

	return
}

type BigIntJSON struct {
	v *big.Int
}

func (b BigIntJSON) MarshalJSON() ([]byte, error) {
	v := b.v.Uint64()

	return json.Marshal(v)
}

func (b *BigIntJSON) UnmarshalJSON(data []byte) (err error) {
	var v uint64

	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}

	b.v = big.NewInt(0)
	b.v.SetUint64(v)

	return
}

// The Alert submit to avs
type Alert interface {
	MessageHash() [32]byte
}

// The Alert Request Message
type AlertRequest struct {
	Alert   Alert
	ResChan chan AlertResponse
}

// The Alert Response Message
type AlertResponse struct {
	Code      uint32
	TxHash    [32]byte
	TaskIndex uint32
	Err       error
	Msg       string
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
	L2BlockNumber BigIntJSON `json:"l2_block_number"`
}

// Return the message hash for signature in avs
func (a AlertBlockMismatch) MessageHash() [32]byte {
	var res [32]byte

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(a.InvalidOutputRoot[:])
	hasher.Write(a.ExpectOutputRoot[:])
	hasher.Write(a.L2BlockNumber.v.Bytes())
	copy(res[:], hasher.Sum(nil)[:32])

	return res
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
	InvalidOutputIndex BigIntJSON `json:"invalid_output_index"`
}

// Return the message hash for signature in avs
func (a AlertBlockOutputOracleMismatch) MessageHash() [32]byte {
	var res [32]byte

	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(a.ExpectOutputRoot[:])
	hasher.Write(a.InvalidOutputIndex.v.Bytes())
	copy(res[:], hasher.Sum(nil)[:32])

	return res
}

var _ Alert = (*AlertBlockOutputOracleMismatch)(nil)

//	AlertBlockOutputOracleMismatch is Submit alert for verifier found a op block output root mismatch.
//
// It just a warning without any prove, the prover verifier should
// submit a prove to ensure the alert is valid.
// This alert only for the porposaled output root by proposer,
// so we just sumit the index for this output root.
type AlertBlockHashMismatch struct {
	// The block hash which to alert
	Hash HexEncodedBytes32 `json:"hash"`
}

// Return the message hash for signature in avs
func (a AlertBlockHashMismatch) MessageHash() [32]byte {
	return a.Hash
}

var _ Alert = (*AlertBlockHashMismatch)(nil)
