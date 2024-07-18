// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

import "forge-std/Script.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";

import "../src/core/AVSCreator.sol";
import {IndexRegistry} from "eigenlayer-middleware/IndexRegistry.sol";
import {StakeRegistry} from "eigenlayer-middleware/StakeRegistry.sol";
import {BLSApkRegistry} from "eigenlayer-middleware/BLSApkRegistry.sol";
import {RegistryCoordinator} from "eigenlayer-middleware/RegistryCoordinator.sol";
import {MachServiceManager} from "../src/core/MachServiceManager.sol";
import {StakeRegistry, IStrategy} from "eigenlayer-middleware/StakeRegistry.sol";

// DELEGATION_MANAGER=$DELEGATION_MANAGER AVS_DIRECTORY=$AVS_DIRECTORY INITIAL_OWNER=$INITIAL_OWNER forge script ./script/AVSCreatorDeployer.s.sol \
// --private-key $PK \
// --rpc-url $URL \
// --etherscan-api-key $API_KEY \
// --broadcast -vvvv --slow --verify

struct TokenAndWeight {
    address token;
    uint96 weight;
}

contract AVSCreatorDeployer is Script {
    function run() external {
        vm.startBroadcast();

        address initialOwner = vm.envAddress("INITIAL_OWNER");
        address delegationManager = vm.envAddress("DELEGATION_MANAGER");
        address avsDirectory = vm.envAddress("AVS_DIRECTORY");

        AVSCreator creator = new AVSCreator(delegationManager, avsDirectory);

        // Set the bytecodes
        creator.setIndexRegistryBytecode(type(IndexRegistry).creationCode);
        creator.setStakeRegistryBytecode(type(StakeRegistry).creationCode);
        creator.setApkRegistryBytecode(type(BLSApkRegistry).creationCode);
        creator.setRegistryCoordinatorBytecode(type(RegistryCoordinator).creationCode);

        uint256 numQuorums = 1;
        uint256 numStrategies = 3;
        uint96 minimumStake = 0;
        uint32 maxOperatorCount = 50;
        {
            // strategies deployed
            TokenAndWeight[] memory deployedStrategyArray = new TokenAndWeight[](numStrategies);

            {
                // need manually step in
                deployedStrategyArray[0].token = 0xbeaC0eeEeeeeEEeEeEEEEeeEEeEeeeEeeEEBEaC0;
                deployedStrategyArray[1].token = 0x7D704507b76571a51d9caE8AdDAbBFd0ba0e63d3;
                deployedStrategyArray[2].token = 0x3A8fBdf9e77DFc25d09741f51d3E181b25d0c4E0;
            }

            {
                // need manually step in
                deployedStrategyArray[0].weight = 1000000000000000000;
                deployedStrategyArray[1].weight = 997992210000000000;
                deployedStrategyArray[2].weight = 1104234999999999999;
            }

            IRegistryCoordinator.OperatorSetParam[] memory operatorSetParams =
                new IRegistryCoordinator.OperatorSetParam[](numQuorums);

            // prepare _operatorSetParams
            for (uint256 i = 0; i < numQuorums; i++) {
                // hard code these for now
                operatorSetParams[i] = IRegistryCoordinator.OperatorSetParam({
                    maxOperatorCount: maxOperatorCount,
                    kickBIPsOfOperatorStake: 11000, // an operator needs to have kickBIPsOfOperatorStake / 10000 times the stake of the operator with the least stake to kick them out
                    kickBIPsOfTotalStake: 1001 // an operator needs to have less than kickBIPsOfTotalStake / 10000 of the total stake to be kicked out
                });
            }

            // prepare _minimumStakes
            uint96[] memory minimumStakeForQuourm = new uint96[](numQuorums);
            for (uint256 i = 0; i < numQuorums; i++) {
                minimumStakeForQuourm[i] = minimumStake;
            }

            // prepare _strategyParams
            IStakeRegistry.StrategyParams[][] memory strategyParams = new IStakeRegistry.StrategyParams[][](numQuorums);
            for (uint256 i = 0; i < numQuorums; i++) {
                IStakeRegistry.StrategyParams[] memory params = new IStakeRegistry.StrategyParams[](numStrategies);
                for (uint256 j = 0; j < numStrategies; j++) {
                    params[j] = IStakeRegistry.StrategyParams({
                        strategy: IStrategy(deployedStrategyArray[j].token),
                        multiplier: deployedStrategyArray[j].weight
                    });
                }
                strategyParams[i] = params;
            }

            creator.createAVS(
                initialOwner, initialOwner, initialOwner, operatorSetParams, minimumStakeForQuourm, strategyParams
            );
        }

        // type(MachServiceManager).creationCode
        // abi.encode(avsDirectory, registryCoordinatorProxy_, stakeRegistryProxy_)

        // Stop broadcasting transactions
        vm.stopBroadcast();
    }
}
