// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;
// solhint-disable

import "forge-std/Script.sol";
import "../test/ServiceManager.test.sol";

// anvil --fork-url https://eth-goerli.g.alchemy.com/v2/<api-key>
// FILE='/Users/x/z/avs/script/config/deploy.goerli.json' forge script script/Deploy.s.sol:Deploy --rpc-url http://127.0.0.1:8545 --broadcast -vvvv --slow
contract Deploy is Script, AVSDeployer {
    using BN254 for BN254.G1Point;

    function _parseStakeRegistryParams(
        string memory data
    )
        internal
        pure
        returns (
            uint96[] memory minimumStakeForQuourm,
            IVoteWeigher.StrategyAndWeightingMultiplier[][]
                memory strategyAndWeightingMultipliers
        )
    {
        bytes memory stakesConfigsRaw = stdJson.parseRaw(
            data,
            ".minimumStakes"
        );
        minimumStakeForQuourm = abi.decode(stakesConfigsRaw, (uint96[]));

        bytes memory strategyConfigsRaw = stdJson.parseRaw(
            data,
            ".strategyWeights"
        );
        strategyAndWeightingMultipliers = abi.decode(
            strategyConfigsRaw,
            (IVoteWeigher.StrategyAndWeightingMultiplier[][])
        );
    }

    function _parseRegistryCoordinatorParams(
        string memory data
    )
        internal
        returns (
            IBLSRegistryCoordinatorWithIndices.OperatorSetParam[]
                memory operatorSetParams,
            address churner,
            address ejector
        )
    {
        bytes memory operatorConfigsRaw = stdJson.parseRaw(
            data,
            ".operatorSetParams"
        );
        operatorSetParams = abi.decode(
            operatorConfigsRaw,
            (IBLSRegistryCoordinatorWithIndices.OperatorSetParam[])
        );

        churner = stdJson.readAddress(data, ".permissions.churner");
        ejector = stdJson.readAddress(data, ".permissions.ejector");
    }

    function run() public {
        string memory data = vm.readFile(vm.envString("FILE"));

        // READ JSON DATA
        bytes32 pk = stdJson.readBytes32(data, ".pk");
        vm.startBroadcast(uint256(pk));

        address admin = stdJson.readAddress(data, ".permissions.owner");

        {
            ProxyAdmin proxyAdmin = new ProxyAdmin();
            EmptyContract emptyContract = new EmptyContract();
            PauserRegistry pauserRegistry;
            {
                address[] memory pausers = new address[](1);
                pausers[0] = admin;
                pauserRegistry = new PauserRegistry(pausers, admin);
            }

            Slasher slasher = Slasher(stdJson.readAddress(data, ".slasher"));
            StrategyManager strategyManager = StrategyManager(
                stdJson.readAddress(data, ".strategyManager")
            );

            (
                uint96[] memory minimumStakeForQuorum,
                IVoteWeigher.StrategyAndWeightingMultiplier[][]
                    memory quorumStrategiesConsideredAndMultipliers
            ) = _parseStakeRegistryParams(data);
            (
                IBLSRegistryCoordinatorWithIndices.OperatorSetParam[]
                    memory operatorSetParams,
                address churner,
                address ejector
            ) = _parseRegistryCoordinatorParams(data);
            _deployAVS(
                admin,
                emptyContract,
                proxyAdmin,
                pauserRegistry,
                slasher,
                strategyManager,
                quorumStrategiesConsideredAndMultipliers,
                minimumStakeForQuorum,
                operatorSetParams
            );
        }

        DelegationManager delegation = DelegationManager(
            stdJson.readAddress(data, ".delegation")
        );
        IDelegationManager.OperatorDetails
            memory operatorDetails = IDelegationManager.OperatorDetails({
                earningsReceiver: admin,
                delegationApprover: address(0),
                stakerOptOutWindowBlocks: 0
            });
        string memory emptyStringForMetadataURI;
        // delegation.registerAsOperator(
        //     operatorDetails,
        //     emptyStringForMetadataURI
        // );

        BN254.G1Point memory messageHash = compendium.getMessageHash(admin);

        uint256 privKey = 69;
        BN254.G1Point memory pubKeyG1;
        BN254.G2Point memory pubKeyG2;
        BN254.G1Point memory signedMessageHash;

        pubKeyG1 = BN254.generatorG1().scalar_mul(privKey);
        //privKey*G2
        pubKeyG2.X[
            1
        ] = 19101821850089705274637533855249918363070101489527618151493230256975900223847;
        pubKeyG2.X[
            0
        ] = 5334410886741819556325359147377682006012228123419628681352847439302316235957;
        pubKeyG2.Y[
            1
        ] = 354176189041917478648604979334478067325821134838555150300539079146482658331;
        pubKeyG2.Y[
            0
        ] = 4185483097059047421902184823581361466320657066600218863748375739772335928910;

        signedMessageHash = BN254.scalar_mul(messageHash, privKey);
        compendium.registerBLSPublicKey(signedMessageHash, pubKeyG1, pubKeyG2);

        bytes memory quorumNumbers = new bytes(1);
        quorumNumbers[0] = bytes1(0);
        string memory defaultSocket = "69.69.69.69:420";
        registryCoordinator.registerOperatorWithCoordinator(
            quorumNumbers,
            pubKeyG1,
            defaultSocket
        );

        vm.stopBroadcast();
    }
}
