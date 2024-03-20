// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import "eigenlayer-core/test/mocks/EmptyContract.sol";
import "forge-std/Script.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "eigenlayer-core/contracts/strategies/StrategyBaseTVLLimits.sol";
import {PauserRegistry} from "eigenlayer-core/contracts/permissions/PauserRegistry.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry, IDelegationManager} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {IIndexRegistry} from "eigenlayer-middleware/interfaces/IIndexRegistry.sol";
import {IBLSApkRegistry} from "eigenlayer-middleware/interfaces/IBLSApkRegistry.sol";
import {RegistryCoordinator} from "eigenlayer-middleware/RegistryCoordinator.sol";
import {IndexRegistry} from "eigenlayer-middleware/IndexRegistry.sol";
import {StakeRegistry, IStrategy} from "eigenlayer-middleware/StakeRegistry.sol";
import {BLSApkRegistry} from "eigenlayer-middleware/BLSApkRegistry.sol";
import {MachServiceManager} from "../src/core/MachServiceManager.sol";
import {IMachServiceManager} from "../src/interfaces/IMachServiceManager.sol";

// forge script script/MachServiceManagerDeployer.s.sol --rpc-url $RPC_URL  --private-key $PRIVATE_KEY --broadcast -vvvv
contract MachServiceManagerDeployer is Script {
    struct MachServiceContract {
        MachServiceManager machServiceManager;
        RegistryCoordinator registryCoordinator;
        IRegistryCoordinator registryCoordinatorImplementation;
        IIndexRegistry indexRegistry;
        IIndexRegistry indexRegistryImplementation;
        IStakeRegistry stakeRegistry;
        IStakeRegistry stakeRegistryImplementation;
        BLSApkRegistry apkRegistry;
        BLSApkRegistry apkRegistryImplementation;
    }

    struct AddressConfig {
        // from team
        address machAVSCommunityMultisig;
        address machAVSPauser;
        address churner;
        address ejector;
        address confirmer;
        // from eigenlayer contracts
        address delegation;
        address eigenDAPauserReg;
    }

    function run() external {
        uint8 numStrategies;
        uint256 maxOperatorCount;
        // strategies deployed
        StrategyBaseTVLLimits[] memory deployedStrategyArray;
        vm.startBroadcast();
        // deploy proxy admin for ability to upgrade proxy contracts
        ProxyAdmin machAVSProxyAdmin = new ProxyAdmin();
        EmptyContract emptyContract = new EmptyContract();

        AddressConfig memory addressConfig;
        addressConfig.machAVSCommunityMultisig = msg.sender;
        addressConfig.machAVSPauser = msg.sender;
        addressConfig.churner = msg.sender;
        addressConfig.ejector = msg.sender;
        addressConfig.confirmer = msg.sender;

        PauserRegistry pauserRegistry;

        // deploy pauser registry
        {
            address[] memory pausers = new address[](2);
            pausers[0] = addressConfig.machAVSPauser;
            pausers[1] = addressConfig.machAVSCommunityMultisig;
            pauserRegistry = new PauserRegistry(pausers, addressConfig.machAVSCommunityMultisig);
        }

        MachServiceContract memory machServiceContract;

        /**
         * First, deploy upgradeable proxy contracts that **will point** to the implementations. Since the implementation contracts are
         * not yet deployed, we give these proxies an empty contract as the initial implementation, to act as if they have no code.
         */
        machServiceContract.machServiceManager = MachServiceManager(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(machAVSProxyAdmin), ""))
        );
        machServiceContract.registryCoordinator = RegistryCoordinator(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(machAVSProxyAdmin), ""))
        );
        machServiceContract.indexRegistry = IIndexRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(machAVSProxyAdmin), ""))
        );
        machServiceContract.stakeRegistry = IStakeRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(machAVSProxyAdmin), ""))
        );
        machServiceContract.apkRegistry = BLSApkRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(machAVSProxyAdmin), ""))
        );

        machServiceContract.indexRegistryImplementation = new IndexRegistry(machServiceContract.registryCoordinator);
        machAVSProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(machServiceContract.indexRegistry))),
            address(machServiceContract.indexRegistryImplementation)
        );

        machServiceContract.stakeRegistryImplementation =
            new StakeRegistry(machServiceContract.registryCoordinator, IDelegationManager(addressConfig.delegation));
        machAVSProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(machServiceContract.stakeRegistry))),
            address(machServiceContract.stakeRegistryImplementation)
        );

        machServiceContract.apkRegistryImplementation = new BLSApkRegistry(machServiceContract.registryCoordinator);
        machAVSProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(machServiceContract.apkRegistry))),
            address(machServiceContract.apkRegistryImplementation)
        );
        machServiceContract.registryCoordinatorImplementation = new RegistryCoordinator(
            IMachServiceManager(address(machServiceContract.machServiceManager)),
            machServiceContract.stakeRegistry,
            machServiceContract.apkRegistry,
            machServiceContract.indexRegistry
        );

        {
            IRegistryCoordinator.OperatorSetParam[] memory operatorSetParams =
                new IRegistryCoordinator.OperatorSetParam[](numStrategies);
            for (uint256 i = 0; i < numStrategies; i++) {
                // hard code these for now
                operatorSetParams[i] = IRegistryCoordinator.OperatorSetParam({
                    maxOperatorCount: uint32(maxOperatorCount),
                    kickBIPsOfOperatorStake: 11000, // an operator needs to have kickBIPsOfOperatorStake / 10000 times the stake of the operator with the least stake to kick them out
                    kickBIPsOfTotalStake: 1001 // an operator needs to have less than kickBIPsOfTotalStake / 10000 of the total stake to be kicked out
                });
            }
            uint96[] memory minimumStakeForQuourm = new uint96[](numStrategies);
            IStakeRegistry.StrategyParams[][] memory strategyAndWeightingMultipliers =
                new IStakeRegistry.StrategyParams[][](numStrategies);
            for (uint256 i = 0; i < numStrategies; i++) {
                strategyAndWeightingMultipliers[i] = new IStakeRegistry.StrategyParams[](1);
                strategyAndWeightingMultipliers[i][0] = IStakeRegistry.StrategyParams({
                    strategy: IStrategy(address(deployedStrategyArray[i])),
                    multiplier: 1 ether
                });
            }
            machAVSProxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(machServiceContract.registryCoordinator))),
                address(machServiceContract.registryCoordinatorImplementation),
                abi.encodeWithSelector(
                    RegistryCoordinator.initialize.selector,
                    addressConfig.machAVSCommunityMultisig,
                    addressConfig.churner,
                    addressConfig.ejector,
                    IPauserRegistry(addressConfig.eigenDAPauserReg),
                    0, // initial paused status is nothing paused
                    operatorSetParams,
                    minimumStakeForQuourm,
                    strategyAndWeightingMultipliers
                )
            );
        }
        vm.stopBroadcast();
    }
}
