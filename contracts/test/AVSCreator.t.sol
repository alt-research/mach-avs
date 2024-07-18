// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

import "forge-std/Test.sol";
import "../src/core/AVSCreator.sol";
import {IndexRegistry} from "eigenlayer-middleware/IndexRegistry.sol";
import {BLSApkRegistry} from "eigenlayer-middleware/BLSApkRegistry.sol";
import {RegistryCoordinator} from "eigenlayer-middleware/RegistryCoordinator.sol";
import {MachServiceManager} from "../src/core/MachServiceManager.sol";
import {StakeRegistry, IStrategy} from "eigenlayer-middleware/StakeRegistry.sol";

contract AVSCreatorTest is Test {
    AVSCreator public avsCreator;
    address delegationManager = address(0x1);
    address avsDirectory = address(0x2);

    event AVSCreated(
        ProxyAdmin ProxyAdmin,
        PauserRegistry pauserRegistry,
        address indexRegistryProxy,
        address stakeRegistryProxy,
        address apkRegistryProxy,
        address registryCoordinatorProxy,
        address serviceManagerProxy
    );

    function setUp() public {
        avsCreator = new AVSCreator(delegationManager, avsDirectory);
    }

    // Test initial contract setup
    function testInitialSetup() public {
        assertEq(avsCreator.delegationManager(), delegationManager);
        assertEq(avsCreator.avsDirectory(), avsDirectory);
    }

    // Test setting bytecode functions
    function testSetIndexRegistryBytecode() public {
        bytes memory bytecode = "dummyBytecode";
        avsCreator.setIndexRegistryBytecode(bytecode);
        assertEq(avsCreator.indexRegistryBytecode(), bytecode);

        // Test to ensure it cannot be set again
        vm.expectRevert(AVSCreator.AlreadySet.selector);
        avsCreator.setIndexRegistryBytecode(bytecode);
    }

    function testSetStakeRegistryBytecode() public {
        bytes memory bytecode = "dummyBytecode";
        avsCreator.setStakeRegistryBytecode(bytecode);
        assertEq(avsCreator.stakeRegistryBytecode(), bytecode);

        vm.expectRevert(AVSCreator.AlreadySet.selector);
        avsCreator.setStakeRegistryBytecode(bytecode);
    }

    function testSetApkRegistryBytecode() public {
        bytes memory bytecode = "dummyBytecode";
        avsCreator.setApkRegistryBytecode(bytecode);
        assertEq(avsCreator.apkRegistryBytecode(), bytecode);

        vm.expectRevert(AVSCreator.AlreadySet.selector);
        avsCreator.setApkRegistryBytecode(bytecode);
    }

    function testSetRegistryCoordinatorBytecode() public {
        bytes memory bytecode = "dummyBytecode";
        avsCreator.setRegistryCoordinatorBytecode(bytecode);
        assertEq(avsCreator.registryCoordinatorBytecode(), bytecode);

        vm.expectRevert(AVSCreator.AlreadySet.selector);
        avsCreator.setRegistryCoordinatorBytecode(bytecode);
    }

    // Test the creation of the AVS
    function testCreateAVS() public {
        // Ensure bytecode is set before creating AVS
        avsCreator.setIndexRegistryBytecode(type(IndexRegistry).creationCode);
        avsCreator.setStakeRegistryBytecode(type(StakeRegistry).creationCode);
        avsCreator.setApkRegistryBytecode(type(BLSApkRegistry).creationCode);
        avsCreator.setRegistryCoordinatorBytecode(type(RegistryCoordinator).creationCode);

        IRegistryCoordinator.OperatorSetParam[] memory operatorSetParams =
            new IRegistryCoordinator.OperatorSetParam[](1);
        operatorSetParams[0] = IRegistryCoordinator.OperatorSetParam({
            maxOperatorCount: 50,
            kickBIPsOfOperatorStake: 11000,
            kickBIPsOfTotalStake: 1001
        });

        uint96[] memory minimumStakes = new uint96[](1);
        minimumStakes[0] = 0;

        IStakeRegistry.StrategyParams[][] memory strategyParams = new IStakeRegistry.StrategyParams[][](1);
        strategyParams[0] = new IStakeRegistry.StrategyParams[](3);
        strategyParams[0][0] = IStakeRegistry.StrategyParams({
            strategy: IStrategy(address(0xbeaC0eeEeeeeEEeEeEEEEeeEEeEeeeEeeEEBEaC0)),
            multiplier: 1000000000000000000
        });
        strategyParams[0][1] = IStakeRegistry.StrategyParams({
            strategy: IStrategy(address(0x7D704507b76571a51d9caE8AdDAbBFd0ba0e63d3)),
            multiplier: 997992210000000000
        });
        strategyParams[0][2] = IStakeRegistry.StrategyParams({
            strategy: IStrategy(address(0x3A8fBdf9e77DFc25d09741f51d3E181b25d0c4E0)),
            multiplier: 1104234999999999999
        });
        // Expect the Created event
        vm.expectEmit(true, true, true, false);
        emit AVSCreated(
            ProxyAdmin(address(0)),
            PauserRegistry(address(0)),
            address(0),
            address(0),
            address(0),
            address(0),
            address(0)
        );
        avsCreator.createAVS(
            address(this), address(this), address(this), operatorSetParams, minimumStakes, strategyParams
        );

        // You can add assertions here to verify the state after createAVS is called
    }
}
