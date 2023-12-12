// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

import {IRiscZeroVerifier} from "../src/interfaces/IMachOptimism.sol";
import {IMachOptimismL2OutputOracle} from "../src/interfaces/IMachOptimismL2OutputOracle.sol";
import {RiscZeroGroth16Verifier, ControlID} from "../src/groth16/RiscZeroGroth16Verifier.sol";

import {IBLSRegistryCoordinatorWithIndices, ServiceManager, ServiceManagerBase} from "../src/ServiceManager.sol";
import {TransparentUpgradeableProxy} from "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "eigenlayer-contracts/src/contracts/strategies/StrategyBase.sol";

import "@openzeppelin/contracts/token/ERC20/presets/ERC20PresetFixedSupply.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/proxy/beacon/IBeacon.sol";
import "@openzeppelin/contracts/proxy/beacon/UpgradeableBeacon.sol";

import "eigenlayer-contracts/src/contracts/interfaces/IDelegationManager.sol";
import "eigenlayer-contracts/src/contracts/core/DelegationManager.sol";

import "eigenlayer-contracts/src/contracts/interfaces/IETHPOSDeposit.sol";
import "eigenlayer-contracts/src/contracts/interfaces/IBeaconChainOracle.sol";

import "eigenlayer-contracts/src/contracts/core/StrategyManager.sol";
import "eigenlayer-contracts/src/contracts/strategies/StrategyBase.sol";
import "eigenlayer-contracts/src/contracts/core/Slasher.sol";

import "eigenlayer-contracts/src/contracts/pods/EigenPod.sol";
import "eigenlayer-contracts/src/contracts/pods/EigenPodManager.sol";
import "eigenlayer-contracts/src/contracts/pods/DelayedWithdrawalRouter.sol";

import "eigenlayer-contracts/src/contracts/permissions/PauserRegistry.sol";

import "eigenlayer-contracts/src/test/utils/Operators.sol";

import "eigenlayer-contracts/src/test/mocks/LiquidStakingToken.sol";
import "eigenlayer-contracts/src/test/mocks/EmptyContract.sol";
import "eigenlayer-contracts/src/test/mocks/ETHDepositMock.sol";
import "eigenlayer-contracts/src/test/mocks/BeaconChainOracleMock.sol";

import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";

import {Slasher} from "eigenlayer-contracts/src/contracts/core/Slasher.sol";
import {ISlasher} from "eigenlayer-contracts/src/contracts/interfaces/ISlasher.sol";
import {PauserRegistry} from "eigenlayer-contracts/src/contracts/permissions/PauserRegistry.sol";
import {IStrategyManager} from "eigenlayer-contracts/src/contracts/interfaces/IStrategyManager.sol";
import {IStrategy} from "eigenlayer-contracts/src/contracts/interfaces/IStrategy.sol";
import {ISignatureUtils} from "eigenlayer-contracts/src/contracts/interfaces/ISignatureUtils.sol";

import {BitmapUtils} from "eigenlayer-middleware/src/libraries/BitmapUtils.sol";
import {BN254, BLSPublicKeyCompendium} from "eigenlayer-middleware/src/BLSPublicKeyCompendium.sol";
import {BLSOperatorStateRetriever} from "eigenlayer-middleware/src/BLSOperatorStateRetriever.sol";
import {BLSRegistryCoordinatorWithIndices} from "eigenlayer-middleware/src/BLSRegistryCoordinatorWithIndices.sol";
import {BLSPubkeyRegistry} from "eigenlayer-middleware/src/BLSPubkeyRegistry.sol";
import {StakeRegistry} from "eigenlayer-middleware/src/StakeRegistry.sol";
import {IndexRegistry} from "eigenlayer-middleware/src/IndexRegistry.sol";
import {IServiceManager} from "eigenlayer-middleware/src/interfaces/IServiceManager.sol";
import {IBLSPubkeyRegistry} from "eigenlayer-middleware/src/interfaces/IBLSPubkeyRegistry.sol";
import {IRegistryCoordinator} from "eigenlayer-middleware/src/interfaces/IRegistryCoordinator.sol";
import {IVoteWeigher} from "eigenlayer-middleware/src/interfaces/IVoteWeigher.sol";
import {IStakeRegistry} from "eigenlayer-middleware/src/interfaces/IStakeRegistry.sol";
import {IIndexRegistry} from "eigenlayer-middleware/src/interfaces/IIndexRegistry.sol";

import "forge-std/Test.sol";

contract EigenLayerDeployer is Test {
    Vm cheats = Vm(HEVM_ADDRESS);

    Slasher public slasher;
    DelegationManager public delegation;
    StrategyManager public strategyManager;
    EigenPodManager public eigenPodManager;
    IEigenPod public pod;
    IDelayedWithdrawalRouter public delayedWithdrawalRouter;
    IETHPOSDeposit public ethPOSDeposit;
    IBeacon public eigenPodBeacon;

    // testing/mock contracts
    IERC20 public eigenToken;
    IERC20 public weth;
    StrategyBase public wethStrat;
    StrategyBase public eigenStrat;
    StrategyBase public baseStrategyImplementation;

    mapping(uint256 => IStrategy) public strategies;

    uint256 wethInitialSupply = 10e50;
    uint256 public constant eigenTotalSupply = 1000e18;
    uint256 nonce = 69;
    uint256 public gasLimit = 750000;
    uint32 PARTIAL_WITHDRAWAL_FRAUD_PROOF_PERIOD_BLOCKS = 7 days / 12 seconds;
    uint256 REQUIRED_BALANCE_WEI = 32 ether;
    uint64 MAX_PARTIAL_WTIHDRAWAL_AMOUNT_GWEI = 1 ether / 1e9;
    uint64 MAX_RESTAKED_BALANCE_GWEI_PER_VALIDATOR = 32e9;
    uint64 GOERLI_GENESIS_TIME = 1616508000;

    address theMultiSig = address(420);
    address operator = address(0x4206904396bF2f8b173350ADdEc5007A52664293); //sk: e88d9d864d5d731226020c5d2f02b62a4ce2a4534a39c225d32d3db795f83319
    address _challenger = address(0x6966904396bF2f8b173350bCcec5007A52669873);

    address public constant CONTRACT_OWNER = address(101);

    function _deployEigenLayerContractsLocal(
        EmptyContract emptyContract,
        ProxyAdmin proxyAdmin,
        PauserRegistry pauserRegistry
    ) internal returns (Slasher, StrategyManager) {
        /**
         * First, deploy upgradeable proxy contracts that **will point** to the implementations. Since the implementation contracts are
         * not yet deployed, we give these proxies an empty contract as the initial implementation, to act as if they have no code.
         */
        delegation = DelegationManager(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(proxyAdmin),
                    ""
                )
            )
        );
        strategyManager = StrategyManager(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(proxyAdmin),
                    ""
                )
            )
        );
        slasher = Slasher(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(proxyAdmin),
                    ""
                )
            )
        );
        eigenPodManager = EigenPodManager(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(proxyAdmin),
                    ""
                )
            )
        );
        delayedWithdrawalRouter = DelayedWithdrawalRouter(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(proxyAdmin),
                    ""
                )
            )
        );

        ethPOSDeposit = new ETHPOSDepositMock();
        pod = new EigenPod(
            ethPOSDeposit,
            delayedWithdrawalRouter,
            eigenPodManager,
            MAX_RESTAKED_BALANCE_GWEI_PER_VALIDATOR,
            GOERLI_GENESIS_TIME
        );

        eigenPodBeacon = new UpgradeableBeacon(address(pod));

        // Second, deploy the *implementation* contracts, using the *proxy contracts* as inputs
        DelegationManager delegationImplementation = new DelegationManager(
            strategyManager,
            slasher,
            eigenPodManager
        );
        StrategyManager strategyManagerImplementation = new StrategyManager(
            delegation,
            eigenPodManager,
            slasher
        );
        Slasher slasherImplementation = new Slasher(
            strategyManager,
            delegation
        );
        EigenPodManager eigenPodManagerImplementation = new EigenPodManager(
            ethPOSDeposit,
            eigenPodBeacon,
            strategyManager,
            slasher,
            delegation
        );
        DelayedWithdrawalRouter delayedWithdrawalRouterImplementation = new DelayedWithdrawalRouter(
                eigenPodManager
            );

        // Third, upgrade the proxy contracts to use the correct implementation contracts and initialize them.
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(delegation))),
            address(delegationImplementation),
            abi.encodeWithSelector(
                DelegationManager.initialize.selector,
                CONTRACT_OWNER,
                pauserRegistry,
                0 /*initialPausedStatus*/
            )
        );
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(strategyManager))),
            address(strategyManagerImplementation),
            abi.encodeWithSelector(
                StrategyManager.initialize.selector,
                CONTRACT_OWNER,
                CONTRACT_OWNER,
                pauserRegistry,
                0 /*initialPausedStatus*/
            )
        );
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(slasher))),
            address(slasherImplementation),
            abi.encodeWithSelector(
                Slasher.initialize.selector,
                CONTRACT_OWNER,
                pauserRegistry,
                0 /*initialPausedStatus*/
            )
        );
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(eigenPodManager))),
            address(eigenPodManagerImplementation),
            abi.encodeWithSelector(
                EigenPodManager.initialize.selector,
                type(uint256).max, // maxPods
                address(0),
                CONTRACT_OWNER,
                pauserRegistry,
                0 /*initialPausedStatus*/
            )
        );
        uint256 initPausedStatus = 0;
        uint256 withdrawalDelayBlocks = PARTIAL_WITHDRAWAL_FRAUD_PROOF_PERIOD_BLOCKS;
        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(
                payable(address(delayedWithdrawalRouter))
            ),
            address(delayedWithdrawalRouterImplementation),
            abi.encodeWithSelector(
                DelayedWithdrawalRouter.initialize.selector,
                CONTRACT_OWNER,
                pauserRegistry,
                initPausedStatus,
                withdrawalDelayBlocks
            )
        );

        //simple ERC20 (**NOT** WETH-like!), used in a test strategy
        weth = new ERC20PresetFixedSupply(
            "weth",
            "WETH",
            wethInitialSupply,
            CONTRACT_OWNER
        );

        // deploy StrategyBase contract implementation, then create upgradeable proxy that points to implementation and initialize it
        baseStrategyImplementation = new StrategyBase(strategyManager);
        wethStrat = StrategyBase(
            address(
                new TransparentUpgradeableProxy(
                    address(baseStrategyImplementation),
                    address(proxyAdmin),
                    abi.encodeWithSelector(
                        StrategyBase.initialize.selector,
                        weth,
                        pauserRegistry
                    )
                )
            )
        );

        eigenToken = new ERC20PresetFixedSupply(
            "eigen",
            "EIGEN",
            wethInitialSupply,
            CONTRACT_OWNER
        );

        // deploy upgradeable proxy that points to StrategyBase implementation and initialize it
        eigenStrat = StrategyBase(
            address(
                new TransparentUpgradeableProxy(
                    address(baseStrategyImplementation),
                    address(proxyAdmin),
                    abi.encodeWithSelector(
                        StrategyBase.initialize.selector,
                        eigenToken,
                        pauserRegistry
                    )
                )
            )
        );

        return (slasher, strategyManager);
    }
}

contract AVSDeployer {
    BLSPublicKeyCompendium compendium;
    BLSOperatorStateRetriever public operatorStateRetriever;

    ServiceManager public serviceManager;
    ServiceManager public serviceManagerImplementation;

    BLSRegistryCoordinatorWithIndices public registryCoordinator;
    BLSRegistryCoordinatorWithIndices public registryCoordinatorImplementation;

    StakeRegistry public stakeRegistry;
    StakeRegistry public stakeRegistryImplementation;

    BLSPubkeyRegistry public blsPubkeyRegistry;
    BLSPubkeyRegistry public blsPubkeyRegistryImplementation;

    IndexRegistry public indexRegistry;
    IndexRegistry public indexRegistryImplementation;

    function _deployAVS(
        address contractOwner,
        EmptyContract emptyContract,
        ProxyAdmin proxyAdmin,
        PauserRegistry pauserRegistry,
        Slasher slasher,
        StrategyManager strategyManager,
        IVoteWeigher.StrategyAndWeightingMultiplier[][]
            memory quorumStrategiesConsideredAndMultipliers,
        uint96[] memory minimumStakeForQuorum,
        IBLSRegistryCoordinatorWithIndices.OperatorSetParam[]
            memory operatorSetParams
    ) internal {
        // Compendium
        //
        compendium = new BLSPublicKeyCompendium();
        // make the CONTRACT_OWNER the owner of the serviceManager contract
        serviceManager = ServiceManager(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(proxyAdmin),
                    ""
                )
            )
        );

        registryCoordinator = BLSRegistryCoordinatorWithIndices(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(proxyAdmin),
                    ""
                )
            )
        );

        stakeRegistry = StakeRegistry(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(proxyAdmin),
                    ""
                )
            )
        );

        indexRegistry = IndexRegistry(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(proxyAdmin),
                    ""
                )
            )
        );

        blsPubkeyRegistry = BLSPubkeyRegistry(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(proxyAdmin),
                    ""
                )
            )
        );

        stakeRegistryImplementation = new StakeRegistry(
            registryCoordinator,
            strategyManager,
            serviceManager
        );

        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(stakeRegistry))),
            address(stakeRegistryImplementation),
            abi.encodeWithSelector(
                StakeRegistry.initialize.selector,
                minimumStakeForQuorum,
                quorumStrategiesConsideredAndMultipliers
            )
        );

        registryCoordinatorImplementation = new BLSRegistryCoordinatorWithIndices(
            slasher,
            serviceManager,
            stakeRegistry,
            blsPubkeyRegistry,
            indexRegistry
        );

        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(registryCoordinator))),
            address(registryCoordinatorImplementation),
            abi.encodeWithSelector(
                BLSRegistryCoordinatorWithIndices.initialize.selector,
                address(0),
                address(0),
                operatorSetParams,
                pauserRegistry,
                0 /*initialPausedStatus*/
            )
        );

        blsPubkeyRegistryImplementation = new BLSPubkeyRegistry(
            registryCoordinator,
            BLSPublicKeyCompendium(address(compendium))
        );

        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(blsPubkeyRegistry))),
            address(blsPubkeyRegistryImplementation)
        );

        indexRegistryImplementation = new IndexRegistry(registryCoordinator);

        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(indexRegistry))),
            address(indexRegistryImplementation)
        );

        operatorStateRetriever = new BLSOperatorStateRetriever();

        serviceManagerImplementation = new ServiceManager(
            IBLSRegistryCoordinatorWithIndices(address(registryCoordinator)),
            ISlasher(address(slasher))
        );

        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(serviceManager))),
            address(serviceManagerImplementation),
            abi.encodeWithSelector(
                ServiceManager.initialize.selector,
                pauserRegistry,
                contractOwner,
                0xf8ea046803fdc1fe07d3b56188ff60e3e3c23b771395f9c1244e800b788a4dcc,
                IMachOptimismL2OutputOracle(
                    address(0xAaE1866Bc68c49ede8b779d6c5Ad61b0C3FeAB86)
                ),
                new RiscZeroGroth16Verifier(
                    ControlID.CONTROL_ID_0,
                    ControlID.CONTROL_ID_1
                )
            )
        );
    }
}

contract ServiceManagerTest is AVSDeployer, EigenLayerDeployer {
    using BN254 for BN254.G1Point;

    address public constant ALICE = address(102);

    function setUp() public {
        vm.startPrank(CONTRACT_OWNER);

        address pauser = address(69);
        address unpauser = address(489);
        address[] memory pausers = new address[](1);
        pausers[0] = pauser;
        PauserRegistry pauserRegistry = new PauserRegistry(pausers, unpauser);
        ProxyAdmin proxyAdmin = new ProxyAdmin();

        EmptyContract emptyContract = new EmptyContract();
        (
            Slasher slasher,
            StrategyManager strategyManager
        ) = _deployEigenLayerContractsLocal(
                emptyContract,
                proxyAdmin,
                pauserRegistry
            );

        {
            uint8 numQuorums = 1;
            // setup the dummy quorum strategies
            IVoteWeigher.StrategyAndWeightingMultiplier[][]
                memory quorumStrategiesConsideredAndMultipliers = new IVoteWeigher.StrategyAndWeightingMultiplier[][](
                    numQuorums
                );
            for (
                uint256 i = 0;
                i < quorumStrategiesConsideredAndMultipliers.length;
                i++
            ) {
                quorumStrategiesConsideredAndMultipliers[
                    i
                ] = new IVoteWeigher.StrategyAndWeightingMultiplier[](1);
                quorumStrategiesConsideredAndMultipliers[i][0] = IVoteWeigher
                    .StrategyAndWeightingMultiplier(wethStrat, 1);
            }

            // setup the dummy minimum stake for quorum
            uint96[] memory minimumStakeForQuorum = new uint96[](numQuorums);
            for (uint256 i; i < minimumStakeForQuorum.length; i++) {
                minimumStakeForQuorum[i] = 1000;
            }

            IBLSRegistryCoordinatorWithIndices.OperatorSetParam[]
                memory operatorSetParams = new IBLSRegistryCoordinatorWithIndices.OperatorSetParam[](
                    numQuorums
                );

            for (uint i = 0; i < numQuorums; i++) {
                // hard code these for now
                operatorSetParams[i] = IBLSRegistryCoordinatorWithIndices
                    .OperatorSetParam({
                        maxOperatorCount: 10,
                        kickBIPsOfOperatorStake: 15000,
                        kickBIPsOfTotalStake: 150
                    });
            }

            _deployAVS(
                CONTRACT_OWNER,
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

        IStrategy[] memory strategiesToWhitelist = new IStrategy[](1);
        strategiesToWhitelist[0] = wethStrat;
        strategyManager.addStrategiesToDepositWhitelist(strategiesToWhitelist);

        vm.stopPrank();
    }

    function testRegister() public {
        uint256 amountToDeposit = 1000 ether;
        uint256 privKey = 69;

        vm.startPrank(CONTRACT_OWNER);
        weth.transfer(ALICE, amountToDeposit);
        vm.stopPrank();

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

        vm.startPrank(ALICE);

        IDelegationManager.OperatorDetails
            memory operatorDetails = IDelegationManager.OperatorDetails({
                earningsReceiver: ALICE,
                delegationApprover: address(0),
                stakerOptOutWindowBlocks: 0
            });
        string memory emptyStringForMetadataURI;
        delegation.registerAsOperator(
            operatorDetails,
            emptyStringForMetadataURI
        );

        BN254.G1Point memory messageHash = compendium.getMessageHash(ALICE);
        signedMessageHash = BN254.scalar_mul(messageHash, privKey);
        compendium.registerBLSPublicKey(signedMessageHash, pubKeyG1, pubKeyG2);

        weth.approve(address(strategyManager), amountToDeposit);

        strategyManager.depositIntoStrategy(wethStrat, weth, amountToDeposit);
        require(
            strategyManager.stakerStrategyShares(ALICE, wethStrat) ==
                amountToDeposit,
            "amountToDeposit mismatch"
        );

        require(
            stakeRegistry.weightOfOperatorForQuorum(0, ALICE) ==
                amountToDeposit / 1 ether,
            "weight mismatch"
        );

        bytes memory quorumNumbers = new bytes(1);
        quorumNumbers[0] = bytes1(0);
        string memory defaultSocket = "69.69.69.69:420";
        registryCoordinator.registerOperatorWithCoordinator(
            quorumNumbers,
            pubKeyG1,
            defaultSocket
        );
        assertEq(
            registryCoordinator.getOperatorId(ALICE),
            BN254.hashG1Point(pubKeyG1),
            "operator ID mismatch"
        );
        vm.stopPrank();
    }
}
