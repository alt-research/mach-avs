package chainio

import (
	"context"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/avsregistry"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	logging "github.com/Layr-Labs/eigensdk-go/logging"

	csservicemanager "github.com/alt-research/avs/contracts/bindings/MachServiceManager"
	"github.com/alt-research/avs/legacy/core/config"
	"github.com/alt-research/avs/legacy/core/message"
)

type AvsWriterer interface {
	SendConfirmAlert(ctx context.Context,
		alertHeader *message.AlertTaskInfo,
		nonSignerStakesAndSignature csservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature,
	) (*types.Receipt, error)
}

type AvsWriter struct {
	*avsregistry.ChainWriter
	AvsContractBindings *AvsManagersBindings
	logger              logging.Logger
	TxMgr               txmgr.TxManager
}

var _ AvsWriterer = (*AvsWriter)(nil)

func BuildAvsWriterFromConfig(c *config.Config) (*AvsWriter, error) {
	return BuildAvsWriter(c.TxMgr, c.RegistryCoordinatorAddr, c.OperatorStateRetrieverAddr, c.ServiceManagerAddress, c.EthHttpClient, c.Logger)
}

func BuildAvsWriter(txMgr txmgr.TxManager, registryCoordinatorAddr, operatorStateRetrieverAddr, serviceManagerAddr gethcommon.Address, ethHttpClient eth.HttpBackend, logger logging.Logger) (*AvsWriter, error) {
	avsServiceBindings, err := NewAvsManagersBindings(registryCoordinatorAddr, operatorStateRetrieverAddr, ethHttpClient, logger)
	if err != nil {
		logger.Error("Failed to create contract bindings", "err", err)
		return nil, err
	}

	logger.Info("build avs", "registryCoordinatorAddr", registryCoordinatorAddr, "operatorStateRetrieverAddr", operatorStateRetrieverAddr)

	avsRegistryWriter, err := avsregistry.NewWriterFromConfig(avsregistry.Config{
		RegistryCoordinatorAddress:    registryCoordinatorAddr,
		OperatorStateRetrieverAddress: operatorStateRetrieverAddr,
		ServiceManagerAddress:         serviceManagerAddr,
	}, ethHttpClient, txMgr, logger)
	if err != nil {
		return nil, err
	}
	return NewAvsWriter(avsRegistryWriter, avsServiceBindings, logger, txMgr), nil
}
func NewAvsWriter(avsRegistryWriter *avsregistry.ChainWriter, avsServiceBindings *AvsManagersBindings, logger logging.Logger, txMgr txmgr.TxManager) *AvsWriter {
	return &AvsWriter{
		ChainWriter:         avsRegistryWriter,
		AvsContractBindings: avsServiceBindings,
		logger:              logger,
		TxMgr:               txMgr,
	}
}

func (w *AvsWriter) SendConfirmAlert(ctx context.Context,
	alertHeader *message.AlertTaskInfo,
	nonSignerStakesAndSignature csservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature,
) (*types.Receipt, error) {
	txOpts, err := w.TxMgr.GetNoSendTxOpts()
	if err != nil {
		w.logger.Errorf("Error getting tx opts")
		return nil, err
	}
	tx, err := w.AvsContractBindings.ServiceManager.ConfirmAlert(txOpts, alertHeader.ToIMachServiceManagerAlertHeader(), nonSignerStakesAndSignature)
	if err != nil {
		w.logger.Error("Error submitting SubmitTaskResponse tx while calling respondToTask", "err", err)
		return nil, err
	}
	receipt, err := w.TxMgr.Send(ctx, tx, true)
	if err != nil {
		w.logger.Errorf("Error submitting respondToTask tx")
		return nil, err
	}
	return receipt, nil
}
