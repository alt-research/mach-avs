package chainio

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"

	sdkavsregistry "github.com/Layr-Labs/eigensdk-go/chainio/clients/avsregistry"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	logging "github.com/Layr-Labs/eigensdk-go/logging"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"

	csservicemanager "github.com/alt-research/avs/contracts/bindings/MachServiceManager"
	"github.com/alt-research/avs/core/config"
)

type AvsReaderer interface {
	sdkavsregistry.AvsRegistryReader

	CheckSignatures(
		ctx context.Context, msgHash [32]byte, quorumNumbers []byte, referenceBlockNumber uint32, nonSignerStakesAndSignature csservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature,
	) (csservicemanager.IBLSSignatureCheckerQuorumStakeTotals, error)

	IsAlertContains(ctx context.Context, messageHash [32]byte) (bool, error)

	// GetQuorumsByBlockNumber
	GetQuorumsByBlockNumber(ctx context.Context, blockNumber uint32) ([]sdktypes.QuorumNum, error)

	// GetQuorumThresholdPercentages
	GetQuorumThresholdPercentages(ctx context.Context, blockNumber uint32, quorums []sdktypes.QuorumNum) ([]sdktypes.QuorumThresholdPercentage, error)
}

type AvsReader struct {
	sdkavsregistry.AvsRegistryReader
	AvsServiceBindings *AvsManagersBindings
	logger             logging.Logger
}

var _ AvsReaderer = (*AvsReader)(nil)

func BuildAvsReaderFromConfig(c *config.Config) (*AvsReader, error) {
	return BuildAvsReader(c.RegistryCoordinatorAddr, c.OperatorStateRetrieverAddr, c.EthHttpClient, c.Logger)
}
func BuildAvsReader(registryCoordinatorAddr, operatorStateRetrieverAddr gethcommon.Address, ethHttpClient eth.Client, logger logging.Logger) (*AvsReader, error) {
	avsManagersBindings, err := NewAvsManagersBindings(registryCoordinatorAddr, operatorStateRetrieverAddr, ethHttpClient, logger)
	if err != nil {
		return nil, err
	}
	avsRegistryReader, err := sdkavsregistry.BuildAvsRegistryChainReader(registryCoordinatorAddr, operatorStateRetrieverAddr, ethHttpClient, logger)
	if err != nil {
		return nil, err
	}
	return NewAvsReader(avsRegistryReader, avsManagersBindings, logger)
}
func NewAvsReader(avsRegistryReader sdkavsregistry.AvsRegistryReader, avsServiceBindings *AvsManagersBindings, logger logging.Logger) (*AvsReader, error) {
	return &AvsReader{
		AvsRegistryReader:  avsRegistryReader,
		AvsServiceBindings: avsServiceBindings,
		logger:             logger,
	}, nil
}

func (r *AvsReader) CheckSignatures(
	ctx context.Context, msgHash [32]byte, quorumNumbers []byte, referenceBlockNumber uint32, nonSignerStakesAndSignature csservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature,
) (csservicemanager.IBLSSignatureCheckerQuorumStakeTotals, error) {
	stakeTotalsPerQuorum, _, err := r.AvsServiceBindings.ServiceManager.CheckSignatures(
		&bind.CallOpts{}, msgHash, quorumNumbers, referenceBlockNumber, nonSignerStakesAndSignature,
	)
	if err != nil {
		return csservicemanager.IBLSSignatureCheckerQuorumStakeTotals{}, err
	}
	return stakeTotalsPerQuorum, nil
}

func (r *AvsReader) IsAlertContains(ctx context.Context, messageHash [32]byte) (bool, error) {
	isContain, err := r.AvsServiceBindings.ServiceManager.Contains(&bind.CallOpts{
		Context: ctx,
	}, messageHash)
	if err != nil {
		return false, err
	}

	return isContain, nil
}

func (r *AvsReader) GetQuorumsByBlockNumber(ctx context.Context, blockNumber uint32) ([]sdktypes.QuorumNum, error) {
	quorumCount, err := r.AvsRegistryReader.GetQuorumCount(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})
	if err != nil {
		return nil, err
	}

	res := make([]sdktypes.QuorumNum, 0, quorumCount)
	for i := uint8(0); i < quorumCount; i++ {
		res = append(res, sdktypes.QuorumNum(i))
	}

	return res, nil
}

func (r *AvsReader) GetQuorumThresholdPercentages(ctx context.Context, blockNumber uint32, quorums []sdktypes.QuorumNum) ([]sdktypes.QuorumThresholdPercentage, error) {
	quorumThresholdPercentage, err := r.AvsServiceBindings.ServiceManager.QuorumThresholdPercentage(&bind.CallOpts{
		Context:     ctx,
		BlockNumber: big.NewInt(int64(blockNumber)),
	})

	if err != nil {
		return nil, err
	}

	res := make([]sdktypes.QuorumThresholdPercentage, 0, len(quorums))
	for i := 0; i < len(quorums); i++ {
		res = append(res, sdktypes.QuorumThresholdPercentage(uint8(quorumThresholdPercentage)))
	}

	return res, nil

}
