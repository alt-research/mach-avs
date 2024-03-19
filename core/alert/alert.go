package alert

import (
	"math/big"

	"golang.org/x/crypto/sha3"
)

// The Alert submit to avs
type Alert interface {
	MessageHash() [32]byte
}

// AlertBlockMismatch is submit alert for verifier found a op block output mismatch.
//
//	It just a warning without any prove, the prover verifier should
//	submit a prove to ensure the alert is valid.
//	This alert can for the blocks which had not proposal its output
//	root to layer1, this block may not the checkpoint.
type AlertBlockMismatch struct {
	// The invalid output root verifier got from op-devnet.
	InvalidOutputRoot [32]byte
	// The output root calc by verifier.
	ExpectOutputRoot [32]byte
	// The layer2 block 's number.
	L2BlockNumber *big.Int
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

var _ Alert = (*AlertBlockMismatch)(nil)

//	AlertBlockOutputOracleMismatch is Submit alert for verifier found a op block output root mismatch.
//
// It just a warning without any prove, the prover verifier should
// submit a prove to ensure the alert is valid.
// This alert only for the porposaled output root by proposer,
// so we just sumit the index for this output root.
type AlertBlockOutputOracleMismatch struct {
	// The output root calc by verifier.
	ExpectOutputRoot [32]byte
	// The invalid output root index.
	InvalidOutputIndex *big.Int
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

var _ Alert = (*AlertBlockOutputOracleMismatch)(nil)
