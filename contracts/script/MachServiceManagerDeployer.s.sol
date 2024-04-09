// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

import "eigenlayer-core/test/mocks/EmptyContract.sol";
import "forge-std/Script.sol";
import "forge-std/console2.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "eigenlayer-core/contracts/strategies/StrategyBaseTVLLimits.sol";
import "eigenlayer-core/contracts/core/AVSDirectory.sol";
import "eigenlayer-core/contracts/core/DelegationManager.sol";
import "eigenlayer-core/contracts/core/StrategyManager.sol";
import "eigenlayer-core/contracts/strategies/StrategyBaseTVLLimits.sol";
import {IAVSDirectory} from "eigenlayer-core/contracts/interfaces/IAVSDirectory.sol";
import {PauserRegistry} from "eigenlayer-core/contracts/permissions/PauserRegistry.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {IStakeRegistry, IDelegationManager} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {IIndexRegistry} from "eigenlayer-middleware/interfaces/IIndexRegistry.sol";
import {IBLSApkRegistry} from "eigenlayer-middleware/interfaces/IBLSApkRegistry.sol";
import {RegistryCoordinator} from "eigenlayer-middleware/RegistryCoordinator.sol";
import {IndexRegistry} from "eigenlayer-middleware/IndexRegistry.sol";
import {StakeRegistry, IStrategy} from "eigenlayer-middleware/StakeRegistry.sol";
import {BLSApkRegistry} from "eigenlayer-middleware/BLSApkRegistry.sol";
import {OperatorStateRetriever} from "eigenlayer-middleware/OperatorStateRetriever.sol";
import {MachServiceManager} from "../src/core/MachServiceManager.sol";
import {IMachServiceManager} from "../src/interfaces/IMachServiceManager.sol";

// forge script script/MachServiceManagerDeployer.s.sol --rpc-url $RPC_URL  --private-key $PRIVATE_KEY --broadcast -vvvv
contract MachServiceManagerDeployer is Script {
    struct MachServiceContract {
        MachServiceManager machServiceManager;
        MachServiceManager machServiceManagerImplementation;
        RegistryCoordinator registryCoordinator;
        IRegistryCoordinator registryCoordinatorImplementation;
        IIndexRegistry indexRegistry;
        IIndexRegistry indexRegistryImplementation;
        IStakeRegistry stakeRegistry;
        IStakeRegistry stakeRegistryImplementation;
        BLSApkRegistry apkRegistry;
        BLSApkRegistry apkRegistryImplementation;
        OperatorStateRetriever operatorStateRetriever;
    }

    struct EigenLayerContracts {
        AVSDirectory avsDirectory;
        DelegationManager delegationManager;
        StrategyManager strategyManager;
        address beaconETH;
        address oETH;
        address ETHx;
        address mETH;
        address sfrxETH;
        address lsETH;
        address cbETH;
        address ankrETH;
        address stETH;
        address osETH;
        address wBETH;
        address swETH;
        address rETH;
        uint96 beaconETH_Multiplier;
        uint96 oETH_Multiplier;
        uint96 ETHx_Multiplier;
        uint96 mETH_Multiplier;
        uint96 sfrxETH_Multiplier;
        uint96 lsETH_Multiplier;
        uint96 cbETH_Multiplier;
        uint96 ankrETH_Multiplier;
        uint96 stETH_Multiplier;
        uint96 osETH_Multiplier;
        uint96 wBETH_Multiplier;
        uint96 swETH_Multiplier;
        uint96 rETH_Multiplier;
    }

    struct TokenAndWeight {
        address token;
        uint96 weight;
    }

    struct DeploymentConfig {
        // from team
        address machAVSCommunityMultisig;
        address machAVSPauser;
        address churner;
        address ejector;
        address whitelister;
        address confirmer;
        uint256 chainId;
        uint256 numStrategies;
        uint256 numQuorum;
        uint256 maxOperatorCount;
        uint96 minimumStake;
        // from eigenlayer contracts
        address avsDirectory;
        address delegationManager;
    }

    function run() external {
        EigenLayerContracts memory eigenLayerContracts;
        DeploymentConfig memory deploymentConfig;

        {
            string memory EIGENLAYER = "EIGENLAYER_ADDRESSES_OUTPUT_PATH";
            string memory defaultPath = "./script/output/deployment_input.json";
            string memory deployedPath = vm.envOr(EIGENLAYER, defaultPath);
            string memory deployedEigenLayerAddresses = vm.readFile(deployedPath);

            bytes memory deployedStrategyManagerData = vm.parseJson(deployedEigenLayerAddresses, ".strategyManager");
            address deployedStrategyManager = abi.decode(deployedStrategyManagerData, (address));
            bytes memory deployedAvsDirectoryData = vm.parseJson(deployedEigenLayerAddresses, ".avsDirectory");
            address deployedAvsDirectory = abi.decode(deployedAvsDirectoryData, (address));
            bytes memory deployedDelegationManagerData = vm.parseJson(deployedEigenLayerAddresses, ".delegationManager");
            address deployedDelegationManager = abi.decode(deployedDelegationManagerData, (address));

            eigenLayerContracts.avsDirectory = AVSDirectory(deployedAvsDirectory);
            eigenLayerContracts.strategyManager = StrategyManager(deployedStrategyManager);
            eigenLayerContracts.delegationManager = DelegationManager(deployedDelegationManager);
            eigenLayerContracts.beaconETH =
                abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".beaconETH"), (address));
            eigenLayerContracts.oETH = abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".oETH"), (address));
            eigenLayerContracts.ETHx = abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".ETHx"), (address));
            eigenLayerContracts.mETH = abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".mETH"), (address));
            eigenLayerContracts.sfrxETH = abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".sfrxETH"), (address));
            eigenLayerContracts.lsETH = abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".lsETH"), (address));
            eigenLayerContracts.cbETH = abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".cbETH"), (address));
            eigenLayerContracts.ankrETH = abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".ankrETH"), (address));
            eigenLayerContracts.stETH = abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".stETH"), (address));
            eigenLayerContracts.osETH = abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".osETH"), (address));
            eigenLayerContracts.wBETH = abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".wBETH"), (address));
            eigenLayerContracts.swETH = abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".swETH"), (address));
            eigenLayerContracts.rETH = abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".rETH"), (address));
            eigenLayerContracts.beaconETH_Multiplier =
                abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".beaconETH_Multiplier"), (uint96));
            eigenLayerContracts.oETH_Multiplier =
                abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".oETH_Multiplier"), (uint96));
            eigenLayerContracts.ETHx_Multiplier =
                abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".ETHx_Multiplier"), (uint96));
            eigenLayerContracts.mETH_Multiplier =
                abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".mETH_Multiplier"), (uint96));
            eigenLayerContracts.sfrxETH_Multiplier =
                abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".sfrxETH_Multiplier"), (uint96));
            eigenLayerContracts.lsETH_Multiplier =
                abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".lsETH_Multiplier"), (uint96));
            eigenLayerContracts.cbETH_Multiplier =
                abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".cbETH_Multiplier"), (uint96));
            eigenLayerContracts.ankrETH_Multiplier =
                abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".ankrETH_Multiplier"), (uint96));
            eigenLayerContracts.stETH_Multiplier =
                abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".stETH_Multiplier"), (uint96));
            eigenLayerContracts.osETH_Multiplier =
                abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".osETH_Multiplier"), (uint96));
            eigenLayerContracts.wBETH_Multiplier =
                abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".wBETH_Multiplier"), (uint96));
            eigenLayerContracts.swETH_Multiplier =
                abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".swETH_Multiplier"), (uint96));
            eigenLayerContracts.rETH_Multiplier =
                abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".rETH_Multiplier"), (uint96));

            {
                deploymentConfig.machAVSCommunityMultisig =
                    abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".owner"), (address));

                deploymentConfig.churner = abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".churner"), (address));

                deploymentConfig.ejector = abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".ejector"), (address));

                deploymentConfig.confirmer =
                    abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".confirmer"), (address));

                deploymentConfig.whitelister =
                    abi.decode(vm.parseJson(deployedEigenLayerAddresses, ".whitelister"), (address));
            }
        }

        deploymentConfig.chainId = 10;
        deploymentConfig.numQuorum = 1;
        deploymentConfig.maxOperatorCount = 50;
        deploymentConfig.minimumStake = 0;
        deploymentConfig.numStrategies = 13;

        deploymentConfig.avsDirectory = address(eigenLayerContracts.avsDirectory);
        deploymentConfig.delegationManager = address(eigenLayerContracts.delegationManager);

        // strategies deployed
        TokenAndWeight[] memory deployedStrategyArray = new TokenAndWeight[](deploymentConfig.numStrategies);

        {
            // need manually step in
            deployedStrategyArray[0].token = eigenLayerContracts.beaconETH;
            deployedStrategyArray[1].token = eigenLayerContracts.oETH;
            deployedStrategyArray[2].token = eigenLayerContracts.ETHx;
            deployedStrategyArray[3].token = eigenLayerContracts.mETH;
            deployedStrategyArray[4].token = eigenLayerContracts.sfrxETH;
            deployedStrategyArray[5].token = eigenLayerContracts.lsETH;
            deployedStrategyArray[6].token = eigenLayerContracts.cbETH;
            deployedStrategyArray[7].token = eigenLayerContracts.ankrETH;
            deployedStrategyArray[8].token = eigenLayerContracts.stETH;
            deployedStrategyArray[9].token = eigenLayerContracts.osETH;
            deployedStrategyArray[10].token = eigenLayerContracts.wBETH;
            deployedStrategyArray[11].token = eigenLayerContracts.swETH;
            deployedStrategyArray[12].token = eigenLayerContracts.rETH;
        }

        {
            // need manually step in
            deployedStrategyArray[0].weight = eigenLayerContracts.beaconETH_Multiplier;
            deployedStrategyArray[1].weight = eigenLayerContracts.oETH_Multiplier;
            deployedStrategyArray[2].weight = eigenLayerContracts.ETHx_Multiplier;
            deployedStrategyArray[3].weight = eigenLayerContracts.mETH_Multiplier;
            deployedStrategyArray[4].weight = eigenLayerContracts.sfrxETH_Multiplier;
            deployedStrategyArray[5].weight = eigenLayerContracts.lsETH_Multiplier;
            deployedStrategyArray[6].weight = eigenLayerContracts.cbETH_Multiplier;
            deployedStrategyArray[7].weight = eigenLayerContracts.ankrETH_Multiplier;
            deployedStrategyArray[8].weight = eigenLayerContracts.stETH_Multiplier;
            deployedStrategyArray[9].weight = eigenLayerContracts.osETH_Multiplier;
            deployedStrategyArray[10].weight = eigenLayerContracts.wBETH_Multiplier;
            deployedStrategyArray[11].weight = eigenLayerContracts.swETH_Multiplier;
            deployedStrategyArray[12].weight = eigenLayerContracts.rETH_Multiplier;
        }

        vm.startBroadcast();
        // deploy proxy admin for ability to upgrade proxy contracts
        ProxyAdmin machAVSProxyAdmin = new ProxyAdmin();
        EmptyContract emptyContract = new EmptyContract();

        PauserRegistry pauserRegistry;

        // deploy pauser registry
        {
            address[] memory pausers = new address[](1);
            pausers[0] = deploymentConfig.machAVSCommunityMultisig;
            pauserRegistry = new PauserRegistry(pausers, deploymentConfig.machAVSCommunityMultisig);
        }

        MachServiceContract memory machServiceContract;

        /**
         * First, deploy upgradeable proxy contracts that **will point** to the implementations. Since the implementation contracts are
         * not yet deployed, we give these proxies an empty contract as the initial implementation, to act as if they have no code.
         */
        machServiceContract.indexRegistry = IIndexRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(machAVSProxyAdmin), ""))
        );
        machServiceContract.stakeRegistry = IStakeRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(machAVSProxyAdmin), ""))
        );
        machServiceContract.apkRegistry = BLSApkRegistry(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(machAVSProxyAdmin), ""))
        );
        machServiceContract.registryCoordinator = RegistryCoordinator(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(machAVSProxyAdmin), ""))
        );
        machServiceContract.machServiceManager = MachServiceManager(
            address(new TransparentUpgradeableProxy(address(emptyContract), address(machAVSProxyAdmin), ""))
        );

        // Second, deploy the *implementation* contracts, using the *proxy contracts* as inputs
        machServiceContract.indexRegistryImplementation = new IndexRegistry(machServiceContract.registryCoordinator);
        machAVSProxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(machServiceContract.indexRegistry))),
            address(machServiceContract.indexRegistryImplementation)
        );

        machServiceContract.stakeRegistryImplementation = new StakeRegistry(
            machServiceContract.registryCoordinator, IDelegationManager(deploymentConfig.delegationManager)
        );
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
        machServiceContract.operatorStateRetriever = new OperatorStateRetriever();

        {
            IRegistryCoordinator.OperatorSetParam[] memory operatorSetParams =
                new IRegistryCoordinator.OperatorSetParam[](deploymentConfig.numQuorum);

            // prepare _operatorSetParams
            for (uint256 i = 0; i < deploymentConfig.numQuorum; i++) {
                // hard code these for now
                operatorSetParams[i] = IRegistryCoordinator.OperatorSetParam({
                    maxOperatorCount: uint32(deploymentConfig.maxOperatorCount),
                    kickBIPsOfOperatorStake: 11000, // an operator needs to have kickBIPsOfOperatorStake / 10000 times the stake of the operator with the least stake to kick them out
                    kickBIPsOfTotalStake: 1001 // an operator needs to have less than kickBIPsOfTotalStake / 10000 of the total stake to be kicked out
                });
            }

            // prepare _minimumStakes
            uint96[] memory minimumStakeForQuourm = new uint96[](deploymentConfig.numQuorum);
            for (uint256 i = 0; i < deploymentConfig.numQuorum; i++) {
                minimumStakeForQuourm[i] = deploymentConfig.minimumStake;
            }

            // prepare _strategyParams
            IStakeRegistry.StrategyParams[][] memory strategyParams =
                new IStakeRegistry.StrategyParams[][](deploymentConfig.numQuorum);
            for (uint256 i = 0; i < deploymentConfig.numQuorum; i++) {
                IStakeRegistry.StrategyParams[] memory params =
                    new IStakeRegistry.StrategyParams[](deploymentConfig.numStrategies);
                for (uint256 j = 0; j < deploymentConfig.numStrategies; j++) {
                    params[j] = IStakeRegistry.StrategyParams({
                        strategy: IStrategy(deployedStrategyArray[j].token),
                        multiplier: deployedStrategyArray[j].weight
                    });
                }
                strategyParams[i] = params;
            }

            // initialize
            machAVSProxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(machServiceContract.registryCoordinator))),
                address(machServiceContract.registryCoordinatorImplementation),
                abi.encodeWithSelector(
                    RegistryCoordinator.initialize.selector,
                    deploymentConfig.machAVSCommunityMultisig,
                    deploymentConfig.churner,
                    deploymentConfig.ejector,
                    IPauserRegistry(pauserRegistry),
                    0, // initial paused status is nothing paused
                    operatorSetParams,
                    minimumStakeForQuourm,
                    strategyParams
                )
            );
        }
        machServiceContract.machServiceManagerImplementation = new MachServiceManager(
            IAVSDirectory(deploymentConfig.avsDirectory),
            machServiceContract.registryCoordinator,
            machServiceContract.stakeRegistry,
            deploymentConfig.chainId
        );
        // Third, upgrade the proxy contracts to use the correct implementation contracts and initialize them.
        machAVSProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(machServiceContract.machServiceManager))),
            address(machServiceContract.machServiceManagerImplementation),
            abi.encodeWithSelector(
                MachServiceManager.initialize.selector,
                IPauserRegistry(pauserRegistry),
                0,
                deploymentConfig.machAVSCommunityMultisig,
                deploymentConfig.confirmer,
                deploymentConfig.whitelister
            )
        );
        vm.stopBroadcast();

        string memory MACH = "MACHAVS_ADDRESSES_OUTPUT_PATH";
        string memory defaultMachPath = "./script/output/machavs_deploy_output.json";
        string memory deployedMachPath = vm.envOr(MACH, defaultMachPath);

        string memory output = "machAVS deployment output";
        vm.serializeAddress(output, "machServiceManager", address(machServiceContract.machServiceManager));
        vm.serializeAddress(
            output, "machServiceManagerImpl", address(machServiceContract.machServiceManagerImplementation)
        );
        vm.serializeAddress(output, "registryCoordinator", address(machServiceContract.registryCoordinator));
        vm.serializeAddress(
            output, "registryCoordinatorImpl", address(machServiceContract.registryCoordinatorImplementation)
        );
        vm.serializeAddress(output, "indexRegistry", address(machServiceContract.indexRegistry));
        vm.serializeAddress(output, "indexRegistryImpl", address(machServiceContract.indexRegistryImplementation));
        vm.serializeAddress(output, "stakeRegistry", address(machServiceContract.stakeRegistry));
        vm.serializeAddress(output, "stakeRegistryImpl", address(machServiceContract.stakeRegistryImplementation));
        vm.serializeAddress(output, "apkRegistry", address(machServiceContract.apkRegistry));
        vm.serializeAddress(output, "apkRegistryImpl", address(machServiceContract.apkRegistryImplementation));
        vm.serializeAddress(output, "pauserRegistry", address(pauserRegistry));
        vm.serializeAddress(output, "machAVSProxyAdmin", address(machAVSProxyAdmin));
        vm.serializeAddress(output, "emptyContract", address(emptyContract));
        vm.serializeAddress(output, "operatorStateRetriever", address(machServiceContract.operatorStateRetriever));
        string memory finalJson = vm.serializeString(output, "object", output);
        vm.createDir("./script/output", true);
        vm.writeJson(finalJson, deployedMachPath);
    }
}
