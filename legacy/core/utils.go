package core

import (
	"math/big"

	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"
	csservicemanager "github.com/alt-research/avs/contracts/bindings/MachServiceManager"
)

// BINDING UTILS - conversion from contract structs to golang structs

// BN254.sol is a library, so bindings for G1 Points and G2 Points are only generated
// in every contract that imports that library. Thus the output here will need to be
// type casted if G1Point is needed to interface with another contract (eg: BLSPublicKeyCompendium.sol)
func ConvertToBN254G1Point(input *bls.G1Point) csservicemanager.BN254G1Point {
	output := csservicemanager.BN254G1Point{
		X: input.X.BigInt(big.NewInt(0)),
		Y: input.Y.BigInt(big.NewInt(0)),
	}
	return output
}

func ConvertToBN254G2Point(input *bls.G2Point) csservicemanager.BN254G2Point {
	output := csservicemanager.BN254G2Point{
		X: [2]*big.Int{input.X.A1.BigInt(big.NewInt(0)), input.X.A0.BigInt(big.NewInt(0))},
		Y: [2]*big.Int{input.Y.A1.BigInt(big.NewInt(0)), input.Y.A0.BigInt(big.NewInt(0))},
	}
	return output
}

func ConvertQuorumNumbersFromBytes(numbers []byte) sdktypes.QuorumNums {
	quorumNumbers := make([]sdktypes.QuorumNum, len(numbers))
	for i, v := range numbers {
		quorumNumbers[i] = sdktypes.QuorumNum(v)
	}
	return quorumNumbers
}

func ConvertQuorumThresholdPercentagesFromBytes(numbers []byte) sdktypes.QuorumThresholdPercentages {
	quorumNumbers := make([]sdktypes.QuorumThresholdPercentage, len(numbers))
	for i, v := range numbers {
		quorumNumbers[i] = sdktypes.QuorumThresholdPercentage(v)
	}
	return quorumNumbers
}
