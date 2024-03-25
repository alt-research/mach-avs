package operator

// OUTDATED
// This file contains cli functions for registering an operator with the AVS and printing status
// However, all of this functionality has been moved to the plugin/ package
// we are just waiting for eigenlayer-cli to be open sourced so we can completely get rid of this registration functionality in the operator

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"

	"github.com/Layr-Labs/eigensdk-go/crypto/bls"
	eigenSdkTypes "github.com/Layr-Labs/eigensdk-go/types"
	sdktypes "github.com/Layr-Labs/eigensdk-go/types"

	regcoord "github.com/Layr-Labs/eigensdk-go/contracts/bindings/RegistryCoordinator"
)

func (o *Operator) RegisterOperatorWithEigenlayer() error {
	op := eigenSdkTypes.Operator{
		Address:                 o.operatorAddr.String(),
		EarningsReceiverAddress: o.operatorAddr.String(),
		MetadataUrl:             o.metadataURI,
	}
	_, err := o.eigenlayerWriter.RegisterAsOperator(context.Background(), op)
	if err != nil {
		o.logger.Error("Error registering operator with eigenlayer", err)
		return err
	}
	return nil
}

func (o *Operator) DepositIntoStrategy(strategyAddr common.Address, amount *big.Int) error {
	_, err := o.eigenlayerWriter.DepositERC20IntoStrategy(context.Background(), strategyAddr, amount)
	if err != nil {
		o.logger.Errorf("Error depositing into strategy", "err", err)
		return err
	}
	return nil
}

// Registration specific functions
func (o *Operator) RegisterOperatorWithAvs(
	operatorEcdsaKeyPair *ecdsa.PrivateKey,
) error {
	// hardcode these things for now
	quorumNumbers := []byte{0}
	socket := "Not Needed"
	operatorToAvsRegistrationSigSalt := [32]byte{123}
	curBlockNum, err := o.ethClient.BlockNumber(context.Background())
	if err != nil {
		o.logger.Errorf("Unable to get current block number")
		return err
	}
	curBlock, err := o.ethClient.BlockByNumber(context.Background(), big.NewInt(int64(curBlockNum)))
	if err != nil {
		o.logger.Errorf("Unable to get current block")
		return err
	}
	sigValidForSeconds := int64(1_000_000)
	operatorToAvsRegistrationSigExpiry := big.NewInt(int64(curBlock.Time()) + sigValidForSeconds)

	o.logger.Info(
		"RegisterOperatorInQuorumWithAVSRegistryCoordinator",
		"quorumNumbers", quorumNumbers[0],
		"socket", socket,
		"operatorToAvsRegistrationSigExpiry", operatorToAvsRegistrationSigExpiry,
	)

	quorumNumbersToSDK := make([]sdktypes.QuorumNum, len(quorumNumbers))
	for i, _ := range quorumNumbers {
		quorumNumbersToSDK[i] = sdktypes.QuorumNum(uint8(quorumNumbers[i]))
	}

	_, err = o.avsWriter.RegisterOperatorInQuorumWithAVSRegistryCoordinator(
		context.Background(),
		operatorEcdsaKeyPair, operatorToAvsRegistrationSigSalt, operatorToAvsRegistrationSigExpiry,
		o.blsKeypair, quorumNumbersToSDK, socket,
	)
	if err != nil {
		o.logger.Error("Unable to register operator with avs registry coordinator", err)
		return err
	}
	o.logger.Infof("Registered operator with avs registry coordinator.")

	return nil
}

// PRINTING STATUS OF OPERATOR: 1
// operator address: 0xa0ee7a142d267c1f36714e4a8f75612f20a79720
// dummy token balance: 0
// delegated shares in dummyTokenStrat: 200
// operator pubkey hash in AVS pubkey compendium (0 if not registered): 0x4b7b8243d970ff1c90a7c775c008baad825893ec6e806dfa5d3663dc093ed17f
// operator is opted in to eigenlayer: true
// operator is opted in to playgroundAVS (aka can be slashed): true
// operator status in AVS registry: REGISTERED
//
//	operatorId: 0x4b7b8243d970ff1c90a7c775c008baad825893ec6e806dfa5d3663dc093ed17f
//	middlewareTimesLen (# of stake updates): 0
//
// operator is frozen: false
type OperatorStatus struct {
	EcdsaAddress string
	// pubkey compendium related
	PubkeysRegistered bool
	G1Pubkey          string
	G2Pubkey          string
	// avs related
	RegisteredWithAvs bool
	OperatorId        string
}

func (o *Operator) PrintOperatorStatus() error {
	fmt.Println("Printing operator status")
	operatorId, err := o.avsReader.GetOperatorId(&bind.CallOpts{}, o.operatorAddr)
	if err != nil {
		return err
	}
	pubkeysRegistered := operatorId != [32]byte{}
	registeredWithAvs := o.operatorId != [32]byte{}
	operatorStatus := OperatorStatus{
		EcdsaAddress:      o.operatorAddr.String(),
		PubkeysRegistered: pubkeysRegistered,
		G1Pubkey:          o.blsKeypair.GetPubKeyG1().String(),
		G2Pubkey:          o.blsKeypair.GetPubKeyG2().String(),
		RegisteredWithAvs: registeredWithAvs,
		OperatorId:        hex.EncodeToString(o.operatorId[:]),
	}
	operatorStatusJson, err := json.MarshalIndent(operatorStatus, "", " ")
	if err != nil {
		return err
	}
	fmt.Println(string(operatorStatusJson))
	return nil
}

func pubKeyG1ToBN254G1Point(p *bls.G1Point) regcoord.BN254G1Point {
	return regcoord.BN254G1Point{
		X: p.X.BigInt(new(big.Int)),
		Y: p.Y.BigInt(new(big.Int)),
	}
}
