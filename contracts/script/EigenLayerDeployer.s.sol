// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/proxy/beacon/UpgradeableBeacon.sol";
import "@openzeppelin/contracts/token/ERC20/presets/ERC20PresetFixedSupply.sol";
import "eigenlayer-core/contracts/permissions/PauserRegistry.sol";
import "eigenlayer-core/test/mocks/EmptyContract.sol";
import "eigenlayer-core/contracts/core/Slasher.sol";
import "eigenlayer-core/contracts/core/AVSDirectory.sol";
import "eigenlayer-core/contracts/core/DelegationManager.sol";
import "eigenlayer-core/contracts/core/StrategyManager.sol";
import "eigenlayer-core/contracts/pods/EigenPodManager.sol";
import "eigenlayer-core/contracts/pods/DelayedWithdrawalRouter.sol";
import "eigenlayer-core/contracts/pods/EigenPod.sol";
import "eigenlayer-core/contracts/strategies/StrategyBase.sol";
import "eigenlayer-core/contracts/strategies/StrategyBaseTVLLimits.sol";

import "forge-std/Script.sol";

struct StrategyUnderlyingTokenConfig {
    address tokenAddress;
    string tokenName;
    string tokenSymbol;
}

// struct used to encode token info in config file
struct StrategyConfig {
    uint256 maxDeposits;
    uint256 maxPerDeposit;
    address tokenAddress;
    string tokenSymbol;
}

// forge script script/EigenLayerDeployer.s.sol --rpc-url $RPC_URL  --private-key $PRIVATE_KEY --broadcast -vvvv
contract EigenLayerDeployer is Script {
    struct EigenLayerContracts {
        ProxyAdmin eigenLayerProxyAdmin;
        PauserRegistry eigenLayerPauserReg;
        EmptyContract emptyContract;
        Slasher slasher;
        Slasher slasherImplementation;
        AVSDirectory avsDirectory;
        AVSDirectory avsDirectoryImplementation;
        DelegationManager delegationManager;
        DelegationManager delegationManagerImplementation;
        StrategyManager strategyManager;
        StrategyManager strategyManagerImplementation;
        EigenPodManager eigenPodManager;
        EigenPodManager eigenPodManagerImplementation;
        DelayedWithdrawalRouter delayedWithdrawalRouter;
        DelayedWithdrawalRouter delayedWithdrawalRouterImplementation;
        UpgradeableBeacon eigenPodBeacon;
        EigenPod eigenPodImplementation;
        StrategyBase baseStrategyImplementation;
    }

    struct Param {
        // StrategyManager
        uint256 STRATEGY_MANAGER_INIT_PAUSED_STATUS;
        address STRATEGY_MANAGER_WHITELISTER;
        // Slasher
        uint256 SLASHER_INIT_PAUSED_STATUS;
        // DelegationManager
        uint256 DELEGATION_MANAGER_INIT_PAUSED_STATUS;
        uint256 DELEGATION_MANAGER_MIN_WITHDRAWAL_DELAY_BLOCKS;
        // AVSDirectory
        uint256 AVS_DIRECTORY_INIT_PAUSED_STATUS;
        // EigenPodManager
        uint256 EIGENPOD_MANAGER_INIT_PAUSED_STATUS;
        uint64 EIGENPOD_MANAGER_DENEB_FORK_TIMESTAMP;
        // EigenPod
        uint64 EIGENPOD_GENESIS_TIME;
        uint64 EIGENPOD_MAX_RESTAKED_BALANCE_GWEI_PER_VALIDATOR;
        uint64 EIGENPOD_MAX_PODS;
        address ETHPOSDepositAddress;
        uint256 DELAYED_WITHDRAWAL_ROUTER_INIT_PAUSED_STATUS;
        // Strategy Deployment
        uint256 STRATEGY_MAX_PER_DEPOSIT;
        uint256 STRATEGY_MAX_TOTAL_DEPOSITS;
        // one week in blocks -- 50400
        uint32 DELAYED_WITHDRAWAL_ROUTER_INIT_WITHDRAWAL_DELAY_BLOCKS;
        uint64 GENESIS_TIME;
    }

    function run() external {
        // PreDeployed Contracts
        address executorMultisig = msg.sender;
        address operationsMultisig = msg.sender;
        address pauserMultisig = msg.sender;
        IBeaconChainOracle beaconOracle = IBeaconChainOracle(0x4C116BB629bff7A8373c2378bBd919f8349B8f25);

        // EigenLayer Contracts
        EigenLayerContracts memory eigenLayerContracts;
        Param memory param;
        // StrategyManager
        param.STRATEGY_MANAGER_INIT_PAUSED_STATUS = 0;
        param.STRATEGY_MANAGER_WHITELISTER = 0x0000000000000000000000000000000000000000;
        // Slasher
        param.SLASHER_INIT_PAUSED_STATUS = 0;
        // DelegationManager
        param.DELEGATION_MANAGER_INIT_PAUSED_STATUS = 0;
        param.DELEGATION_MANAGER_MIN_WITHDRAWAL_DELAY_BLOCKS = 50400;
        // AVSDirectory
        param.AVS_DIRECTORY_INIT_PAUSED_STATUS = 0;
        // EigenPodManager
        param.EIGENPOD_MANAGER_INIT_PAUSED_STATUS = 0;
        param.EIGENPOD_MANAGER_DENEB_FORK_TIMESTAMP = 1707305664;
        // EigenPod
        param.EIGENPOD_GENESIS_TIME = 1695902400;
        param.EIGENPOD_MAX_RESTAKED_BALANCE_GWEI_PER_VALIDATOR = 32;
        param.EIGENPOD_MAX_PODS = 100;
        param.ETHPOSDepositAddress = 0x4242424242424242424242424242424242424242;

        param.DELAYED_WITHDRAWAL_ROUTER_INIT_PAUSED_STATUS = 0;
        vm.startBroadcast();

        // Deploy ProxyAdmin, later set admins for all proxies to be executorMultisig
        eigenLayerContracts.eigenLayerProxyAdmin = new ProxyAdmin();

        address[] memory pausers = new address[](3);
        pausers[0] = executorMultisig;
        pausers[1] = operationsMultisig;
        pausers[2] = pauserMultisig;
        address unpauser = executorMultisig;
        eigenLayerContracts.eigenLayerPauserReg = new PauserRegistry(pausers, unpauser);

        eigenLayerContracts.emptyContract = new EmptyContract();

        eigenLayerContracts.slasher = Slasher(
            address(
                new TransparentUpgradeableProxy(
                    address(eigenLayerContracts.emptyContract), address(eigenLayerContracts.eigenLayerProxyAdmin), ""
                )
            )
        );
        eigenLayerContracts.avsDirectory = AVSDirectory(
            address(
                new TransparentUpgradeableProxy(
                    address(eigenLayerContracts.emptyContract), address(eigenLayerContracts.eigenLayerProxyAdmin), ""
                )
            )
        );
        eigenLayerContracts.delegationManager = DelegationManager(
            address(
                new TransparentUpgradeableProxy(
                    address(eigenLayerContracts.emptyContract), address(eigenLayerContracts.eigenLayerProxyAdmin), ""
                )
            )
        );
        eigenLayerContracts.strategyManager = StrategyManager(
            address(
                new TransparentUpgradeableProxy(
                    address(eigenLayerContracts.emptyContract), address(eigenLayerContracts.eigenLayerProxyAdmin), ""
                )
            )
        );
        eigenLayerContracts.eigenPodManager = EigenPodManager(
            address(
                new TransparentUpgradeableProxy(
                    address(eigenLayerContracts.emptyContract), address(eigenLayerContracts.eigenLayerProxyAdmin), ""
                )
            )
        );
        eigenLayerContracts.delayedWithdrawalRouter = DelayedWithdrawalRouter(
            address(
                new TransparentUpgradeableProxy(
                    address(eigenLayerContracts.emptyContract), address(eigenLayerContracts.eigenLayerProxyAdmin), ""
                )
            )
        );

        eigenLayerContracts.eigenPodImplementation = new EigenPod(
            IETHPOSDeposit(param.ETHPOSDepositAddress),
            eigenLayerContracts.delayedWithdrawalRouter,
            eigenLayerContracts.eigenPodManager,
            param.EIGENPOD_MAX_RESTAKED_BALANCE_GWEI_PER_VALIDATOR,
            param.GENESIS_TIME
        );
        eigenLayerContracts.eigenPodBeacon = new UpgradeableBeacon(address(eigenLayerContracts.eigenPodImplementation));
        eigenLayerContracts.avsDirectoryImplementation = new AVSDirectory(eigenLayerContracts.delegationManager);
        eigenLayerContracts.delegationManagerImplementation = new DelegationManager(
            eigenLayerContracts.strategyManager, eigenLayerContracts.slasher, eigenLayerContracts.eigenPodManager
        );
        eigenLayerContracts.strategyManagerImplementation = new StrategyManager(
            eigenLayerContracts.delegationManager, eigenLayerContracts.eigenPodManager, eigenLayerContracts.slasher
        );
        eigenLayerContracts.slasherImplementation =
            new Slasher(eigenLayerContracts.strategyManager, eigenLayerContracts.delegationManager);
        eigenLayerContracts.eigenPodManagerImplementation = new EigenPodManager(
            IETHPOSDeposit(param.ETHPOSDepositAddress),
            eigenLayerContracts.eigenPodBeacon,
            eigenLayerContracts.strategyManager,
            eigenLayerContracts.slasher,
            eigenLayerContracts.delegationManager
        );
        eigenLayerContracts.delayedWithdrawalRouterImplementation =
            new DelayedWithdrawalRouter(eigenLayerContracts.eigenPodManager);

        // Third, upgrade the proxy contracts to point to the implementations
        IStrategy[] memory initializeStrategiesToSetDelayBlocks = new IStrategy[](0);
        uint256[] memory initializeWithdrawalDelayBlocks = new uint256[](0);
        // AVSDirectory
        eigenLayerContracts.eigenLayerProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenLayerContracts.avsDirectory))),
            address(eigenLayerContracts.avsDirectoryImplementation),
            abi.encodeWithSelector(
                AVSDirectory.initialize.selector,
                executorMultisig, // initialOwner
                eigenLayerContracts.eigenLayerPauserReg,
                param.AVS_DIRECTORY_INIT_PAUSED_STATUS
            )
        );
        // DelegationManager
        eigenLayerContracts.eigenLayerProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenLayerContracts.delegationManager))),
            address(eigenLayerContracts.delegationManagerImplementation),
            abi.encodeWithSelector(
                DelegationManager.initialize.selector,
                executorMultisig, // initialOwner
                eigenLayerContracts.eigenLayerPauserReg,
                param.DELEGATION_MANAGER_INIT_PAUSED_STATUS,
                param.DELEGATION_MANAGER_MIN_WITHDRAWAL_DELAY_BLOCKS,
                initializeStrategiesToSetDelayBlocks,
                initializeWithdrawalDelayBlocks
            )
        );
        // StrategyManager
        eigenLayerContracts.eigenLayerProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenLayerContracts.strategyManager))),
            address(eigenLayerContracts.strategyManagerImplementation),
            abi.encodeWithSelector(
                StrategyManager.initialize.selector,
                executorMultisig, //initialOwner
                param.STRATEGY_MANAGER_WHITELISTER, //initial whitelister
                eigenLayerContracts.eigenLayerPauserReg,
                param.STRATEGY_MANAGER_INIT_PAUSED_STATUS
            )
        );
        // Slasher
        eigenLayerContracts.eigenLayerProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenLayerContracts.slasher))),
            address(eigenLayerContracts.slasherImplementation),
            abi.encodeWithSelector(
                Slasher.initialize.selector,
                executorMultisig,
                eigenLayerContracts.eigenLayerPauserReg,
                param.SLASHER_INIT_PAUSED_STATUS
            )
        );
        // EigenPodManager
        eigenLayerContracts.eigenLayerProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenLayerContracts.eigenPodManager))),
            address(eigenLayerContracts.eigenPodManagerImplementation),
            abi.encodeWithSelector(
                EigenPodManager.initialize.selector,
                param.EIGENPOD_MAX_PODS,
                beaconOracle,
                msg.sender, // initialOwner is msg.sender for now to set forktimestamp later
                eigenLayerContracts.eigenLayerPauserReg,
                param.EIGENPOD_MANAGER_INIT_PAUSED_STATUS
            )
        );
        // Delayed Withdrawal Router
        eigenLayerContracts.eigenLayerProxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenLayerContracts.delayedWithdrawalRouter))),
            address(eigenLayerContracts.delayedWithdrawalRouterImplementation),
            abi.encodeWithSelector(
                DelayedWithdrawalRouter.initialize.selector,
                executorMultisig, // initialOwner
                eigenLayerContracts.eigenLayerPauserReg,
                param.DELAYED_WITHDRAWAL_ROUTER_INIT_PAUSED_STATUS,
                param.DELAYED_WITHDRAWAL_ROUTER_INIT_WITHDRAWAL_DELAY_BLOCKS
            )
        );
        // deploy a token and create a strategy config for each token
        uint256 i = 1;
        address tokenAddress = address(
            new ERC20PresetFixedSupply(
                string(abi.encodePacked("Token", i)), string(abi.encodePacked("TOK", i)), 1000 ether, msg.sender
            )
        );
        StrategyConfig memory strategyConfig = StrategyConfig({
            maxDeposits: type(uint256).max,
            maxPerDeposit: type(uint256).max,
            tokenAddress: tokenAddress,
            tokenSymbol: string(abi.encodePacked("TOK", i))
        });
        eigenLayerContracts.baseStrategyImplementation = new StrategyBaseTVLLimits(eigenLayerContracts.strategyManager);

        StrategyBaseTVLLimits strategyBaseTVLLimits = StrategyBaseTVLLimits(
            address(
                new TransparentUpgradeableProxy(
                    address(eigenLayerContracts.baseStrategyImplementation),
                    address(eigenLayerContracts.eigenLayerProxyAdmin),
                    abi.encodeWithSelector(
                        StrategyBaseTVLLimits.initialize.selector,
                        strategyConfig.maxPerDeposit,
                        strategyConfig.maxDeposits,
                        IERC20(strategyConfig.tokenAddress),
                        eigenLayerContracts.eigenLayerPauserReg
                    )
                )
            )
        );

        vm.stopBroadcast();
        string memory output = "eigenlayer contracts deployment output";
        vm.serializeAddress(output, "eigenLayerProxyAdmin", address(eigenLayerContracts.eigenLayerProxyAdmin));
        vm.serializeAddress(output, "eigenLayerPauserRegistry", address(eigenLayerContracts.eigenLayerPauserReg));
        vm.serializeAddress(output, "slasher", address(eigenLayerContracts.slasher));
        vm.serializeAddress(output, "slasherImplementation", address(eigenLayerContracts.slasherImplementation));
        vm.serializeAddress(output, "avsDirectory", address(eigenLayerContracts.avsDirectory));
        vm.serializeAddress(
            output, "avsDirectoryImplementation", address(eigenLayerContracts.avsDirectoryImplementation)
        );
        vm.serializeAddress(output, "delegationManager", address(eigenLayerContracts.delegationManager));
        vm.serializeAddress(
            output, "delegationManagerImplementation", address(eigenLayerContracts.delegationManagerImplementation)
        );
        vm.serializeAddress(output, "strategyManager", address(eigenLayerContracts.strategyManager));
        vm.serializeAddress(
            output, "strategyManagerImplementation", address(eigenLayerContracts.strategyManagerImplementation)
        );
        vm.serializeAddress(output, "eigenPodManager", address(eigenLayerContracts.eigenPodManager));
        vm.serializeAddress(
            output, "eigenPodManagerImplementation", address(eigenLayerContracts.eigenPodManagerImplementation)
        );
        vm.serializeAddress(output, "delayedWithdrawalRouter", address(eigenLayerContracts.delayedWithdrawalRouter));
        vm.serializeAddress(
            output,
            "delayedWithdrawalRouterImplementation",
            address(eigenLayerContracts.delayedWithdrawalRouterImplementation)
        );
        vm.serializeAddress(output, "eigenPodBeacon", address(eigenLayerContracts.eigenPodBeacon));
        vm.serializeAddress(output, "eigenPodImplementation", address(eigenLayerContracts.eigenPodImplementation));
        vm.serializeAddress(
            output, "baseStrategyImplementation", address(eigenLayerContracts.baseStrategyImplementation)
        );
        vm.serializeAddress(output, "underlayingToken", tokenAddress);
        vm.serializeAddress(output, "strategyBaseTVLLimits", address(strategyBaseTVLLimits));

        vm.createDir("./script/output", true);
        string memory finalJson = vm.serializeString(output, "object", output);
        vm.writeJson(finalJson, "./script/output/eigenlayer_deploy_output.json");
    }
}
