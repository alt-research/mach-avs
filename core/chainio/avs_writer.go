package chainio

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/Layr-Labs/eigensdk-go/chainio/clients/avsregistry"
	"github.com/Layr-Labs/eigensdk-go/chainio/clients/eth"
	"github.com/Layr-Labs/eigensdk-go/chainio/txmgr"
	logging "github.com/Layr-Labs/eigensdk-go/logging"

	csservicemanager "github.com/alt-research/avs/contracts/bindings/MachServiceManager"
	"github.com/alt-research/avs/core/config"
	"github.com/alt-research/avs/core/message"
)

type AvsWriterer interface {
	avsregistry.AvsRegistryWriter

	SendConfirmAlert(ctx context.Context,
		alertHeader *message.AlertTaskInfo,
		nonSignerStakesAndSignature csservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature,
	) (*types.Receipt, error)

	SendGenericConfirm(ctx context.Context,
		task *message.GenericTaskData,
		nonSignerStakesAndSignature csservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature,
	) (*types.Receipt, error)
}

type AvsWriter struct {
	avsregistry.AvsRegistryWriter
	AvsContractBindings *AvsManagersBindings
	logger              logging.Logger
	TxMgr               txmgr.TxManager
	ethHttpClient       eth.Client

	avsCfg *config.GenericAVSConfig
}

var _ AvsWriterer = (*AvsWriter)(nil)

func BuildAvsWriter(
	txMgr txmgr.TxManager,
	registryCoordinatorAddr, operatorStateRetrieverAddr gethcommon.Address,
	ethHttpClient eth.Client,
	logger logging.Logger,
	avsCfg *config.GenericAVSConfig) (*AvsWriter, error) {
	avsServiceBindings, err := NewAvsManagersBindings(registryCoordinatorAddr, operatorStateRetrieverAddr, ethHttpClient, logger)
	if err != nil {
		logger.Error("Failed to create contract bindings", "err", err)
		return nil, err
	}

	logger.Info("build avs", "registryCoordinatorAddr", registryCoordinatorAddr, "operatorStateRetrieverAddr", operatorStateRetrieverAddr)

	avsRegistryWriter, err := avsregistry.BuildAvsRegistryChainWriter(registryCoordinatorAddr, operatorStateRetrieverAddr, logger, ethHttpClient, txMgr)
	if err != nil {
		return nil, err
	}
	return NewAvsWriter(avsRegistryWriter, avsServiceBindings, logger, txMgr, ethHttpClient, avsCfg), nil
}
func NewAvsWriter(
	avsRegistryWriter avsregistry.AvsRegistryWriter,
	avsServiceBindings *AvsManagersBindings,
	logger logging.Logger,
	txMgr txmgr.TxManager,
	ethHttpClient eth.Client,
	avsCfg *config.GenericAVSConfig) *AvsWriter {
	return &AvsWriter{
		AvsRegistryWriter:   avsRegistryWriter,
		AvsContractBindings: avsServiceBindings,
		logger:              logger,
		TxMgr:               txMgr,
		ethHttpClient:       ethHttpClient,

		avsCfg: avsCfg,
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
	receipt, err := w.TxMgr.Send(ctx, tx)
	if err != nil {
		w.logger.Errorf("Error submitting CreateNewTask tx")
		return nil, err
	}
	return receipt, nil
}

func (w *AvsWriter) SendGenericConfirm(ctx context.Context,
	task *message.GenericTaskData,
	nonSignerStakesAndSignature csservicemanager.IBLSSignatureCheckerNonSignerStakesAndSignature,
) (*types.Receipt, error) {
	if w.avsCfg == nil {
		return nil, fmt.Errorf("not use avs config, so cannot send generic confirm")
	}

	txOpts, err := w.TxMgr.GetNoSendTxOpts()
	if err != nil {
		w.logger.Errorf("Error getting tx opts")
		return nil, err
	}

	params := make([]interface{}, 0, len(task.CallParams)+1)
	params = append(params, task.CallParams...)
	params = append(params, nonSignerStakesAndSignature)

	input, err := w.avsCfg.Abi.Pack(task.CallMethod, params...)
	if err != nil {
		return nil, err
	}

	boundContract := bind.NewBoundContract(w.avsCfg.AVSContractAddress, w.avsCfg.Abi, w.ethHttpClient, w.ethHttpClient, w.ethHttpClient)

	tx, err := boundContract.RawTransact(txOpts, input)
	if err != nil {
		w.logger.Error("Error submitting SubmitTaskResponse tx while calling respondToTask", "err", err)
		return nil, err
	}
	receipt, err := w.TxMgr.Send(ctx, tx)
	if err != nil {
		w.logger.Errorf("Error submitting CreateNewTask tx")
		return nil, err
	}
	return receipt, nil
}
