package chainio

import (
	"context"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"

	sdkavsregistry "github.com/Layr-Labs/eigensdk-go/chainio/clients/avsregistry"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	logging "github.com/Layr-Labs/eigensdk-go/logging"

	csservicemanager "github.com/alt-research/avs/contracts/bindings/MachServiceManager"
	"github.com/alt-research/avs/core/config"
)

type AvsReaderer interface {
	sdkavsregistry.AvsRegistryReader

	CheckSignatures(
		ctx context.Context, msgHash [32]byte, quorumNumbers []byte, referenceBlockNumber uint32, nonSignerStakesAndSignature csservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature,
	) (csservicemanager.IBLSSignatureCheckerQuorumStakeTotals, error)
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
