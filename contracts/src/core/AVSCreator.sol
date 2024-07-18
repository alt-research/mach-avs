// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";

import {IRegistryCoordinator} from "eigenlayer-middleware/interfaces/IRegistryCoordinator.sol";
import {OperatorStateRetriever} from "eigenlayer-middleware/OperatorStateRetriever.sol";
import {PauserRegistry, IPauserRegistry} from "eigenlayer-contracts/src/contracts/permissions/PauserRegistry.sol";
import {IStakeRegistry} from "eigenlayer-middleware/interfaces/IStakeRegistry.sol";
import {RegistryCoordinator} from "eigenlayer-middleware/RegistryCoordinator.sol";

contract EmptyContract {}

contract AVSCreatorStorage {
    address public immutable emptyContract;
    address public immutable delegationManager;
    address public immutable avsDirectory;
    OperatorStateRetriever public immutable operatorStateRetriever;

    bytes public indexRegistryBytecode;
    bytes public stakeRegistryBytecode;
    bytes public apkRegistryBytecode;
    bytes public registryCoordinatorBytecode;

    constructor(address delegationManager_, address avsDirectory_) {
        delegationManager = delegationManager_;
        avsDirectory = avsDirectory_;
        emptyContract = address(new EmptyContract());
        operatorStateRetriever = new OperatorStateRetriever();
    }
}
                                                                                                    
                                                        
contract AVSCreator is AVSCreatorStorage {
    error AlreadySet();

    event AVSCreated(
        ProxyAdmin ProxyAdmin,
        PauserRegistry pauserRegistry,
        address indexRegistryProxy,
        address stakeRegistryProxy,
        address apkRegistryProxy,
        address registryCoordinatorProxy,
        address serviceManagerProxy
    );

    modifier setOnce(bytes memory bytecode) {
        if (bytecode.length != 0) {
            revert AlreadySet();
        }
        _;
    }

    constructor(address delegationManager_, address avsDirectory_)
        AVSCreatorStorage(delegationManager_, avsDirectory_)
    {}

    function setIndexRegistryBytecode(bytes calldata bytecode) external setOnce(indexRegistryBytecode) {
        indexRegistryBytecode = bytecode;
    }

    function setStakeRegistryBytecode(bytes calldata bytecode) external setOnce(stakeRegistryBytecode) {
        stakeRegistryBytecode = bytecode;
    }

    function setApkRegistryBytecode(bytes calldata bytecode) external setOnce(apkRegistryBytecode) {
        apkRegistryBytecode = bytecode;
    }

    function setRegistryCoordinatorBytecode(bytes calldata bytecode) external setOnce(registryCoordinatorBytecode) {
        registryCoordinatorBytecode = bytecode;
    }

    function createAVS(
        address initialOwner_,
        address churnApprover_,
        address ejector_,
        IRegistryCoordinator.OperatorSetParam[] memory operatorSetParams_,
        uint96[] memory minimumStakes_,
        IStakeRegistry.StrategyParams[][] memory strategyParams_
    ) external {
        bytes memory registryCoordinatorInit;
        PauserRegistry pauserRegistry;
        {
            {
                address[] memory pausers = new address[](1);
                pausers[0] = initialOwner_;
                pauserRegistry = new PauserRegistry(pausers, initialOwner_);
            }

            registryCoordinatorInit = abi.encodeWithSelector(
                RegistryCoordinator.initialize.selector,
                initialOwner_,
                churnApprover_,
                ejector_,
                pauserRegistry,
                0, /*initialPausedStatus*/
                operatorSetParams_,
                minimumStakes_,
                strategyParams_
            );
        }

        ProxyAdmin proxyAdmin = new ProxyAdmin();

        // Deploy proxies for each contract
        address indexRegistryProxy = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));
        address stakeRegistryProxy = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));
        address apkRegistryProxy = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));
        address registryCoordinatorProxy =
            address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));
        address serviceManagerProxy = address(new TransparentUpgradeableProxy(emptyContract, address(proxyAdmin), ""));

        emit AVSCreated(
            proxyAdmin,
            pauserRegistry,
            indexRegistryProxy,
            stakeRegistryProxy,
            apkRegistryProxy,
            registryCoordinatorProxy,
            serviceManagerProxy
        );

        // Deploy the actual implementation contracts

        {
            address indexRegistryImpl =
                _createImplementation(abi.encodePacked(indexRegistryBytecode, abi.encode(registryCoordinatorProxy)));

            address stakeRegistryImpl = _createImplementation(
                abi.encodePacked(stakeRegistryBytecode, abi.encode(registryCoordinatorProxy, delegationManager))
            );

            address apkRegistryImpl =
                _createImplementation(abi.encodePacked(apkRegistryBytecode, abi.encode(registryCoordinatorProxy)));

            address registryCoordinatorImpl = _createImplementation(
                abi.encodePacked(
                    registryCoordinatorBytecode,
                    abi.encode(serviceManagerProxy, stakeRegistryProxy, apkRegistryProxy, indexRegistryProxy)
                )
            );

            proxyAdmin.upgrade(TransparentUpgradeableProxy(payable(indexRegistryProxy)), indexRegistryImpl);
            proxyAdmin.upgrade(TransparentUpgradeableProxy(payable(stakeRegistryProxy)), stakeRegistryImpl);
            proxyAdmin.upgrade(TransparentUpgradeableProxy(payable(apkRegistryProxy)), apkRegistryImpl);
            proxyAdmin.upgradeAndCall(
                TransparentUpgradeableProxy(payable(address(registryCoordinatorProxy))),
                address(registryCoordinatorImpl),
                registryCoordinatorInit
            );
        }

        proxyAdmin.transferOwnership(initialOwner_);
    }

    function _createImplementation(bytes memory bytecodeWithConstructor) internal returns (address) {
        address addr;
        assembly {
            addr := create(0, add(bytecodeWithConstructor, 0x20), mload(bytecodeWithConstructor))
            if iszero(extcodesize(addr)) { revert(0, 0) }
        }
        return addr;
    }
}
