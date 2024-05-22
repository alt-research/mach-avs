package genericproxy

import (
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	sighHashAbiParams, _ = abi.NewType("tuple", "Hash32SighHashParam", []abi.ArgumentMarshaling{
		{Name: "messageHash", Type: "bytes32"},
		{Name: "referenceBlockNumber", Type: "uint32"},
		{Name: "rollupChainID", Type: "uint256"},
	})

	sighHashAbiArgs = abi.Arguments{
		{Type: sighHashAbiParams, Name: "one"},
	}
)

func CalcSighHash(messageHash [32]byte, referenceBlockNumber uint32, chainId *big.Int) ([32]byte, error) {
	record := struct {
		MessageHash          [32]byte `abi:"messageHash"`
		ReferenceBlockNumber uint32   `abi:"referenceBlockNumber"`
		ChainId              *big.Int `abi:"rollupChainID"`
	}{
		MessageHash:          messageHash,
		ReferenceBlockNumber: referenceBlockNumber,
		ChainId:              chainId,
	}

	packed, err := sighHashAbiArgs.Pack(&record)
	if err != nil {
		return [32]byte{}, err
	}

	return crypto.Keccak256Hash(packed), nil

}
