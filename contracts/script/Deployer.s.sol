// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.9;

import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/proxy/beacon/UpgradeableBeacon.sol";
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

// forge script script/Deployer.s.sol --rpc-url $RPC_URL  --private-key $PRIVATE_KEY --broadcast -vvvv
contract Deployer is Script {
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
        IBeaconChainOracle beaconOracle;
        UpgradeableBeacon eigenPodBeacon;
        EigenPod eigenPodImplementation;
        StrategyBase baseStrategyImplementation;
    }

    struct Param {
        uint64 EIGENPOD_MAX_RESTAKED_BALANCE_GWEI_PER_VALIDATOR;
        uint64 GENESIS_TIME;
        address ETHPOSDepositAddress;
        uint256 AVS_DIRECTORY_INIT_PAUSED_STATUS;
        uint256 DELEGATION_MANAGER_INIT_PAUSED_STATUS;
        uint256 DELEGATION_MANAGER_MIN_WITHDRAWAL_DELAY_BLOCKS;
        address STRATEGY_MANAGER_WHITELISTER;
        uint256 STRATEGY_MANAGER_INIT_PAUSED_STATUS;
        uint256 SLASHER_INIT_PAUSED_STATUS;
        uint256 EIGENPOD_MANAGER_INIT_PAUSED_STATUS;
        uint256 DELAYED_WITHDRAWAL_ROUTER_INIT_PAUSED_STATUS;
        // one week in blocks -- 50400
        uint32 DELAYED_WITHDRAWAL_ROUTER_INIT_WITHDRAWAL_DELAY_BLOCKS;
        // Strategy Deployment
        uint256 STRATEGY_MAX_PER_DEPOSIT;
        uint256 STRATEGY_MAX_TOTAL_DEPOSITS;
    }

    function run() external {
        // PreDeployed Contracts
        address executorMultisig = msg.sender;
        address operationsMultisig = msg.sender;
        address pauserMultisig = msg.sender;

        // EigenLayer Contracts
        EigenLayerContracts memory eigenLayerContracts;
        Param memory param;
        param.EIGENPOD_MAX_RESTAKED_BALANCE_GWEI_PER_VALIDATOR = 32;
        param.GENESIS_TIME = 1695902400;
        param.ETHPOSDepositAddress = 0x4242424242424242424242424242424242424242;
        param.AVS_DIRECTORY_INIT_PAUSED_STATUS = 0;
        param.DELEGATION_MANAGER_INIT_PAUSED_STATUS = 0;
        param.DELEGATION_MANAGER_MIN_WITHDRAWAL_DELAY_BLOCKS = 50400;
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
                eigenLayerContracts.beaconOracle,
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
        vm.stopBroadcast();
    }
}
