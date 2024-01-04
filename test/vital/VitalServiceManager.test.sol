// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;
import "forge-std/Test.sol";

import "@openzeppelin/contracts/utils/cryptography/ECDSA.sol";

import {IRegistryCoordinator, VitalServiceManager, ServiceManagerBase} from "../../src/vital/VitalServiceManager.sol";
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
import {IStrategy} from "eigenlayer-contracts/src/contracts/interfaces/IStrategy.sol";

import {BitmapUtils} from "eigenlayer-middleware/src/libraries/BitmapUtils.sol";
import {OperatorStateRetriever} from "eigenlayer-middleware/src/OperatorStateRetriever.sol";
import {RegistryCoordinator, ISignatureUtils, IServiceManager} from "eigenlayer-middleware/src/RegistryCoordinator.sol";
import {BN254, BLSApkRegistry} from "eigenlayer-middleware/src/BLSApkRegistry.sol";
import {IStakeRegistry, StakeRegistry} from "eigenlayer-middleware/src/StakeRegistry.sol";
import {IndexRegistry} from "eigenlayer-middleware/src/IndexRegistry.sol";
import {IBLSApkRegistry} from "eigenlayer-middleware/src/BLSApkRegistryStorage.sol";

contract EigenLayerDeployer is Test {
    Vm cheats = Vm(HEVM_ADDRESS);

    Slasher public slasher;
    DelegationManager public delegationManager;
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
    ) internal returns (Slasher, StrategyManager, IDelegationManager) {
        /**
         * First, deploy upgradeable proxy contracts that **will point** to the implementations. Since the implementation contracts are
         * not yet deployed, we give these proxies an empty contract as the initial implementation, to act as if they have no code.
         */
        delegationManager = DelegationManager(
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
            delegationManager,
            eigenPodManager,
            slasher
        );
        Slasher slasherImplementation = new Slasher(
            strategyManager,
            delegationManager
        );

        {
            // Third, upgrade the proxy contracts to use the correct implementation contracts and initialize them.
            proxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(
                    payable(address(delegationManager))
                ),
                address(delegationImplementation),
                abi.encodeWithSelector(
                    DelegationManager.initialize.selector,
                    CONTRACT_OWNER,
                    pauserRegistry,
                    0 /*initialPausedStatus*/,
                    10
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
        }
        {
            EigenPodManager eigenPodManagerImplementation = new EigenPodManager(
                ethPOSDeposit,
                eigenPodBeacon,
                strategyManager,
                slasher,
                delegationManager
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
        }
        {
            DelayedWithdrawalRouter delayedWithdrawalRouterImplementation = new DelayedWithdrawalRouter(
                    eigenPodManager
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
        }

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

        return (slasher, strategyManager, delegationManager);
    }
}

contract VitalAVSDeployer {
    OperatorStateRetriever public operatorStateRetriever;

    IServiceManager public serviceManager;
    IServiceManager public serviceManagerImplementation;

    RegistryCoordinator public registryCoordinator;
    RegistryCoordinator public registryCoordinatorImplementation;

    StakeRegistry public stakeRegistry;
    StakeRegistry public stakeRegistryImplementation;

    BLSApkRegistry public blsApkRegistry;
    BLSApkRegistry public blsApkRegistryImplementation;

    IndexRegistry public indexRegistry;
    IndexRegistry public indexRegistryImplementation;

    function _deployAVS(
        address contractOwner,
        EmptyContract emptyContract,
        ProxyAdmin proxyAdmin,
        PauserRegistry pauserRegistry,
        Slasher slasher,
        IDelegationManager delegationManager,
        IStakeRegistry.StrategyParams[][] memory strategyParams,
        uint96[] memory minimumStakeForQuorum,
        IRegistryCoordinator.OperatorSetParam[] memory operatorSetParams
    ) internal {
        // make the CONTRACT_OWNER the owner of the serviceManager contract
        serviceManager = IServiceManager(
            address(
                new TransparentUpgradeableProxy(
                    address(emptyContract),
                    address(proxyAdmin),
                    ""
                )
            )
        );

        registryCoordinator = RegistryCoordinator(
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

        blsApkRegistry = BLSApkRegistry(
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
            delegationManager
        );

        registryCoordinatorImplementation = new RegistryCoordinator(
            serviceManager,
            stakeRegistry,
            blsApkRegistry,
            indexRegistry
        );
        blsApkRegistryImplementation = new BLSApkRegistry(registryCoordinator);

        indexRegistryImplementation = new IndexRegistry(registryCoordinator);

        operatorStateRetriever = new OperatorStateRetriever();

        serviceManagerImplementation = IServiceManager(
            address(
                new VitalServiceManager(
                    IDelegationManager(address(delegationManager)),
                    IRegistryCoordinator(address(registryCoordinator)),
                    IStakeRegistry(address(stakeRegistry))
                )
            )
        );

        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(blsApkRegistry))),
            address(blsApkRegistryImplementation)
        );

        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(indexRegistry))),
            address(indexRegistryImplementation)
        );

        proxyAdmin.upgrade(
            TransparentUpgradeableProxy(payable(address(stakeRegistry))),
            address(stakeRegistryImplementation)
        );

        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(registryCoordinator))),
            address(registryCoordinatorImplementation),
            abi.encodeWithSelector(
                RegistryCoordinator.initialize.selector,
                contractOwner,
                contractOwner,
                contractOwner,
                pauserRegistry,
                0 /*initialPausedStatus*/,
                operatorSetParams,
                minimumStakeForQuorum,
                strategyParams
            )
        );

        proxyAdmin.upgradeAndCall(
            TransparentUpgradeableProxy(payable(address(serviceManager))),
            address(serviceManagerImplementation),
            abi.encodeWithSelector(
                ServiceManagerBase.initialize.selector,
                contractOwner
            )
        );
    }
}

contract VitalServiceManagerTest is VitalAVSDeployer, EigenLayerDeployer {
    using BN254 for BN254.G1Point;

    address public constant ALICE = 0x1326324f5A9fb193409E10006e4EA41b970Df321;

    function setUp() public {
        vm.startPrank(CONTRACT_OWNER);

        address pauser = address(69);
        address unpauser = address(489);
        address[] memory pausers = new address[](1);
        pausers[0] = pauser;
        PauserRegistry pauserRegistry = new PauserRegistry(pausers, unpauser);
        ProxyAdmin proxyAdmin = new ProxyAdmin();

        EmptyContract emptyContract = new EmptyContract();
        Slasher slasher;
        StrategyManager strategyManager;
        IDelegationManager delegationManager;
        {
            (
                slasher,
                strategyManager,
                delegationManager
            ) = _deployEigenLayerContractsLocal(
                emptyContract,
                proxyAdmin,
                pauserRegistry
            );
        }

        {
            uint8 numQuorums = 1;

            // setup the dummy minimum stake for quorum
            uint96[] memory minimumStakeForQuorum = new uint96[](numQuorums);
            for (uint256 i; i < minimumStakeForQuorum.length; i++) {
                minimumStakeForQuorum[i] = 1000;
            }

            IRegistryCoordinator.OperatorSetParam[]
                memory operatorSetParams = new IRegistryCoordinator.OperatorSetParam[](
                    numQuorums
                );

            for (uint i = 0; i < numQuorums; i++) {
                // hard code these for now
                operatorSetParams[i] = IRegistryCoordinator.OperatorSetParam({
                    maxOperatorCount: 10,
                    kickBIPsOfOperatorStake: 15000,
                    kickBIPsOfTotalStake: 150
                });
            }

            // setup the dummy quorum strategies
            IStakeRegistry.StrategyParams[][]
                memory quorumStrategiesConsideredAndMultipliers = new IStakeRegistry.StrategyParams[][](
                    1
                );
            for (
                uint256 i = 0;
                i < quorumStrategiesConsideredAndMultipliers.length;
                i++
            ) {
                quorumStrategiesConsideredAndMultipliers[
                    i
                ] = new IStakeRegistry.StrategyParams[](1);
                quorumStrategiesConsideredAndMultipliers[i][0] = IStakeRegistry
                    .StrategyParams(wethStrat, uint96(1e18));
            }

            _deployAVS(
                CONTRACT_OWNER,
                emptyContract,
                proxyAdmin,
                pauserRegistry,
                slasher,
                delegationManager,
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
        uint256 amountToDeposit = 1000;

        vm.startPrank(CONTRACT_OWNER);
        weth.transfer(ALICE, amountToDeposit);
        vm.stopPrank();

        vm.startPrank(ALICE);

        IDelegationManager.OperatorDetails
            memory operatorDetails = IDelegationManager.OperatorDetails({
                earningsReceiver: ALICE,
                delegationApprover: address(0),
                stakerOptOutWindowBlocks: 0
            });
        string memory emptyStringForMetadataURI;
        delegationManager.registerAsOperator(
            operatorDetails,
            emptyStringForMetadataURI
        );

        weth.approve(address(strategyManager), amountToDeposit);

        strategyManager.depositIntoStrategy(wethStrat, weth, amountToDeposit);
        require(
            strategyManager.stakerStrategyShares(ALICE, wethStrat) ==
                amountToDeposit,
            "amountToDeposit mismatch"
        );

        require(
            stakeRegistry.weightOfOperatorForQuorum(0, ALICE) ==
                amountToDeposit,
            "weight mismatch"
        );

        uint256 privKey = 69;
        IBLSApkRegistry.PubkeyRegistrationParams
            memory pubkeyRegistrationParams;

        BN254.G1Point memory messageHash = registryCoordinator
            .pubkeyRegistrationMessageHash(ALICE);

        pubkeyRegistrationParams.pubkeyRegistrationSignature = BN254.scalar_mul(
            messageHash,
            privKey
        );
        pubkeyRegistrationParams.pubkeyG1 = BN254.generatorG1().scalar_mul(
            privKey
        );

        //privKey*G2
        pubkeyRegistrationParams.pubkeyG2.X[1] = uint256(
            0x2A3B3F7EF4F62985AF31809FDC531483E5F1CD67AA1BCF0F8AC0D17E158AA967
        );
        pubkeyRegistrationParams.pubkeyG2.X[0] = uint256(
            0x0BCB2B68B6C68A5AEA7FE75B5446C4CA410461FA226C2487D07EB2C504639CB5
        );
        pubkeyRegistrationParams.pubkeyG2.Y[1] = uint256(
            0x00C874E4FCFB88D5C98A0240BC6F7F37D45F2226CA147317B3A2B7243DDB6C1B
        );
        pubkeyRegistrationParams.pubkeyG2.Y[0] = uint256(
            0x0940E64478DB51FE630CC540DBEABEA34D072A54FD7C743056E18174F9A1B64E
        );

        ISignatureUtils.SignatureWithSaltAndExpiry
            memory signatureWithSaltAndExpiry;

        uint256 expiry = block.timestamp + 10;
        bytes32 salt = bytes32(uint256(keccak256("defaultSalt")));

        bytes32 digestHash = delegationManager
            .calculateOperatorAVSRegistrationDigestHash(
                ALICE,
                address(serviceManager),
                salt,
                expiry
            );

        (uint8 v, bytes32 r, bytes32 s) = vm.sign(privKey, digestHash);

        signatureWithSaltAndExpiry.expiry = expiry;
        signatureWithSaltAndExpiry.signature = abi.encodePacked(r, s, v);
        signatureWithSaltAndExpiry.salt = salt;

        bytes memory quorumNumbers = new bytes(1);
        quorumNumbers[0] = bytes1(0);
        string memory defaultSocket = "69.69.69.69:420";

        registryCoordinator.registerOperator(
            quorumNumbers,
            defaultSocket,
            pubkeyRegistrationParams,
            signatureWithSaltAndExpiry
        );
        assertEq(
            registryCoordinator.getOperatorId(ALICE),
            BN254.hashG1Point(pubkeyRegistrationParams.pubkeyG1),
            "operator ID mismatch"
        );
        vm.stopPrank();
    }
}
