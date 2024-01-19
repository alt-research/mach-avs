// SPDX-License-Identifier: UNLICENSED
// SEE LICENSE IN https://files.altlayer.io/Alt-Research-License-1.md
// Copyright Alt Research Ltd. 2023. All rights reserved.
//
// You acknowledge and agree that Alt Research Ltd. ("Alt Research") (or Alt
// Research's licensors) own all legal rights, titles and interests in and to the
// work, software, application, source code, documentation and any other documents
pragma solidity ^0.8.23;
// solhint-disable

import "@openzeppelin/contracts/proxy/transparent/TransparentUpgradeableProxy.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import "@openzeppelin/contracts/token/ERC20/presets/ERC20PresetMinterPauser.sol";
import "forge-std/Test.sol";
import "../../src/stake/StakedMultiToken.sol";
import "../../src/stake/AVSStakedMultiToken.sol";

contract StakedMultiTokenMock is AVSStakedMultiToken {
    constructor(address avs_) AVSStakedMultiToken(avs_) {}

    function updateDistribution(uint256 id) external {
        _updateDistribution(id, totalSupply(id));
    }
}

// slither-disable-start all
contract StakedMultiTokenTest is Test {
    address constant DEPLOYER = address(1);
    address constant ALICE = address(2);
    address constant BOB = address(3);

    uint256 constant MAX_BPS = 1e4;
    uint256 constant PRECISION_FACTOR = 1 ether;
    uint256 constant INITIAL_EXCHANGE_RATE = 1 ether;
    uint256 constant EXCHANGE_RATE_UNIT = 1 ether;
    uint256 constant LOWER_BOUND = 1 ether;
    uint256 constant INITIAL_MAX_SLASHABLE_BPS = 3e3;
    uint256 constant COOLDOWN_SECONDS = 5; // 5 sec
    uint256 constant UNSTAKE_WINDOW = 3; // 3 sec
    uint256 constant INITIAL_FUND_AMOUNT = 10 ether;
    uint128 constant EMISSION_PER_SECOND = 10;
    uint128 constant DISTRIBUTION_DURATION_IN_SECOND = 300; // 5 min;

    StakedMultiTokenMock stk;
    ERC20PresetMinterPauser erc20;

    function setUp() external {
        vm.startPrank(DEPLOYER);

        erc20 = new ERC20PresetMinterPauser("Token", "TK");
        TransparentUpgradeableProxy proxy = new TransparentUpgradeableProxy(
            address(new StakedMultiTokenMock(ALICE)),
            address(new ProxyAdmin()),
            ""
        );
        stk = StakedMultiTokenMock(payable(proxy));

        address admin = DEPLOYER;
        stk.initialize(
            admin,
            EMISSION_PER_SECOND,
            DISTRIBUTION_DURATION_IN_SECOND,
            address(erc20),
            address(erc20),
            COOLDOWN_SECONDS,
            UNSTAKE_WINDOW,
            admin,
            admin,
            admin,
            admin,
            admin
        );

        erc20.mint(DEPLOYER, INITIAL_FUND_AMOUNT);
        erc20.approve(address(stk), INITIAL_FUND_AMOUNT);

        erc20.mint(ALICE, INITIAL_FUND_AMOUNT);
        erc20.mint(BOB, INITIAL_FUND_AMOUNT);
        vm.stopPrank();

        vm.startPrank(ALICE);
        erc20.approve(address(stk), INITIAL_FUND_AMOUNT);
        vm.stopPrank();

        vm.startPrank(BOB);
        erc20.approve(address(stk), INITIAL_FUND_AMOUNT);
        vm.stopPrank();

        vm.startPrank(DEPLOYER);
        stk.initializeAsset(0);
        stk.initializeAsset(1);
        vm.stopPrank();
    }

    function test_SetUp() external {
        require(erc20.balanceOf(ALICE) == INITIAL_FUND_AMOUNT);
        require(erc20.balanceOf(BOB) == INITIAL_FUND_AMOUNT);
        require(erc20.allowance(ALICE, address(stk)) == INITIAL_FUND_AMOUNT);
        require(erc20.allowance(BOB, address(stk)) == INITIAL_FUND_AMOUNT);

        require(stk.totalIDs() == 2);

        require(stk.exchangeRates(0) == INITIAL_EXCHANGE_RATE);
        require(stk.exchangeRates(1) == INITIAL_EXCHANGE_RATE);
        require(stk.isValidID(0) == true);
        require(stk.isValidID(1) == true);

        require(stk.exchangeRates(2) == 0);
        require(stk.isValidID(2) == false);

        vm.startPrank(DEPLOYER);
        vm.expectRevert(StakedMultiToken.AlreadyInitialized.selector);
        stk.initializeAsset(1);
        vm.stopPrank();
    }

    function test_RevertIf_MoreThanMaxSupply() external {
        uint256 maxSupply = stk.MAX_SUPPLY();

        vm.startPrank(DEPLOYER);
        erc20.mint(ALICE, maxSupply);
        vm.stopPrank();

        vm.startPrank(ALICE);
        erc20.approve(address(stk), maxSupply);

        vm.expectRevert(StakedMultiToken.MultiplicationOverflow.selector);
        stk.stake(ALICE, 1, maxSupply + 1);

        stk.stake(ALICE, 1, maxSupply);

        vm.expectRevert(StakedMultiToken.MaxSupply.selector);
        stk.stake(ALICE, 1, 1);

        vm.stopPrank();
    }

    function test_Stake() external {
        vm.startPrank(ALICE);
        // Reverts trying to stake with invalid id
        vm.expectRevert(StakedMultiToken.InvalidID.selector);
        stk.stake(ALICE, 2, 1 ether);

        // Reverts trying to stake 0 amount
        vm.expectRevert(StakedMultiToken.ZeroAmount.selector);
        stk.stake(ALICE, 1, 0 ether);

        // Reverts trying to activate cooldown with 0 staked amount
        vm.expectRevert(StakedMultiToken.InvalidBalanceOnCooldown.selector);
        stk.cooldown(1);

        // Alice stakes 10 TK
        stk.stake(ALICE, 1, 1 ether);

        require(erc20.balanceOf(ALICE) == INITIAL_FUND_AMOUNT - 1 ether);
        require(erc20.balanceOf(address(stk)) == 1 ether);
        require(stk.balanceOf(ALICE, 1) == 1 ether);

        // Alice stakes more TK, increasing total STK balance.
        stk.stake(ALICE, 1, 0.5 ether);

        require(erc20.balanceOf(ALICE) == INITIAL_FUND_AMOUNT - 1.5 ether);
        require(erc20.balanceOf(address(stk)) == 1.5 ether);
        require(stk.balanceOf(ALICE, 1) == 1.5 ether);

        vm.stopPrank();
    }

    function test_Claim_Rewards() external {
        vm.startPrank(ALICE);
        stk.stake(ALICE, 1, 1 ether);
        require(stk.totalSupply(1) == 1 ether);
        require(stk.balanceOf(ALICE, 1) == 1 ether);
        require(stk.distributionIndices(1) == 0);

        vm.warp(block.timestamp + 5);
        stk.updateDistribution(1);

        uint256 t0Rewards = 5 * EMISSION_PER_SECOND;

        require(stk.distributionIndices(1) == t0Rewards);
        require(stk.getAccruedRewards(ALICE, 1) == t0Rewards);

        // Claims full rewards
        stk.claimRewards(ALICE, 1, t0Rewards);
        require(
            erc20.balanceOf(ALICE) == INITIAL_FUND_AMOUNT - 1 ether + t0Rewards
        );
        require(stk.getAccruedRewards(ALICE, 1) == 0);
        require(stk.rewardsBalances(1, ALICE) == 0);

        vm.warp(block.timestamp + 10);

        uint256 t1Rewards = 10 * EMISSION_PER_SECOND;

        // Claims half rewards
        stk.claimRewards(ALICE, 1, t1Rewards / 2);
        require(
            erc20.balanceOf(ALICE) ==
                INITIAL_FUND_AMOUNT - 1 ether + t0Rewards + t1Rewards / 2
        );
        require(stk.getAccruedRewards(ALICE, 1) == 0);
        require(stk.rewardsBalances(1, ALICE) == t1Rewards / 2);

        vm.stopPrank();
    }

    function test_Claim_MaxRewards() external {
        vm.startPrank(ALICE);
        stk.stake(ALICE, 1, 1 ether);
        vm.warp(block.timestamp + 5);
        stk.updateDistribution(1);
        uint256 t0Rewards = 5 * EMISSION_PER_SECOND;

        require(stk.distributionIndices(1) == t0Rewards);
        require(stk.getAccruedRewards(ALICE, 1) == t0Rewards);

        // Claims more 10x rewards but gets only the max rewards
        stk.claimRewards(ALICE, 1, t0Rewards * 10);
        require(
            erc20.balanceOf(ALICE) == INITIAL_FUND_AMOUNT - 1 ether + t0Rewards
        );
        require(stk.getAccruedRewards(ALICE, 1) == 0);
        require(stk.rewardsBalances(1, ALICE) == 0);
        vm.stopPrank();
    }

    function test_SetEmissionPerSecond() external {
        vm.startPrank(ALICE);
        // Reverts when a non-emission admin tries to set emission per second
        vm.expectRevert(DistributionManager.NotEmissionAdmin.selector);
        stk.setEmissionPerSecond(0);
        vm.stopPrank();

        vm.startPrank(DEPLOYER);
        stk.setEmissionPerSecond(0);
        require(stk.emissionPerSecond() == 0);

        vm.warp(block.timestamp + 5);
        stk.updateDistribution(1);
        require(stk.distributionIndices(1) == 0);

        vm.stopPrank();
    }

    function test_SetCooldownSeconds() external {
        vm.startPrank(ALICE);
        // Reverts when a non-cooldown admin tries to set cooldown seconds
        vm.expectRevert(StakedMultiToken.NotCooldownAdmin.selector);
        stk.setCooldownSeconds(COOLDOWN_SECONDS * 2);
        vm.stopPrank();

        vm.startPrank(DEPLOYER);
        stk.setCooldownSeconds(COOLDOWN_SECONDS * 2);
        require(stk.cooldownSeconds() == COOLDOWN_SECONDS * 2);
        vm.stopPrank();
    }

    function test_SetUnstakeWindow() external {
        vm.startPrank(ALICE);
        // Reverts when a non-unstake-window admin tries to set unstake window
        vm.expectRevert(StakedMultiToken.NotUnstakeWindowAdmin.selector);
        stk.setUnstakeWindow(UNSTAKE_WINDOW * 2);
        vm.stopPrank();

        vm.startPrank(DEPLOYER);
        stk.setUnstakeWindow(UNSTAKE_WINDOW * 2);
        require(stk.unstakeWindow() == UNSTAKE_WINDOW * 2);
        vm.stopPrank();
    }

    function test_SetMaxSlashableBPS() external {
        vm.startPrank(ALICE);
        // Reverts when a non-slashing admin tries to set max slashable BPS
        vm.expectRevert(StakedMultiToken.NotSlashingAdmin.selector);
        stk.setMaxSlashableBPS(1);
        vm.stopPrank();

        vm.startPrank(DEPLOYER);
        // Reverts when a slashing admin tries to set invalid max slashable BPS
        vm.expectRevert(StakedMultiToken.InvalidBPS.selector);
        stk.setMaxSlashableBPS(0);

        vm.expectRevert(StakedMultiToken.InvalidBPS.selector);
        stk.setMaxSlashableBPS(10000);

        stk.setMaxSlashableBPS(10);
        require(stk.maxSlashableBPS() == 10);
        vm.stopPrank();
    }
}
// slither-disable-end all
