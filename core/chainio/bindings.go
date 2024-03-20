package chainio

import (
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/logging"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"

	regcoord "github.com/Layr-Labs/eigensdk-go/contracts/bindings/RegistryCoordinator"
	csservicemanager "github.com/alt-research/avs/contracts/bindings/MachServiceManager"
)

type AvsManagersBindings struct {
	ServiceManager *csservicemanager.ContractMachServiceManager
	ethClient      eth.Client
	logger         logging.Logger
}

func NewAvsManagersBindings(registryCoordinatorAddr, operatorStateRetrieverAddr gethcommon.Address, ethclient eth.Client, logger logging.Logger) (*AvsManagersBindings, error) {
	contractRegistryCoordinator, err := regcoord.NewContractRegistryCoordinator(registryCoordinatorAddr, ethclient)
	if err != nil {
		return nil, err
	}
	serviceManagerAddr, err := contractRegistryCoordinator.ServiceManager(&bind.CallOpts{})
	if err != nil {
		return nil, err
	}
	contractServiceManager, err := csservicemanager.NewContractMachServiceManager(serviceManagerAddr, ethclient)
	if err != nil {
		logger.Error("Failed to fetch IServiceManager contract", "err", err)
		return nil, err
	}

	return &AvsManagersBindings{
		ServiceManager: contractServiceManager,
		ethClient:      ethclient,
		logger:         logger,
	}, nil
}
