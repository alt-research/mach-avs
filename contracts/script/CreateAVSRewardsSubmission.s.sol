// SPDX-License-Identifier: UNLICENSED
pragma solidity ^0.8.12;

import "eigenlayer-core/contracts/interfaces/IRewardsCoordinator.sol";
import "eigenlayer-middleware/interfaces/IServiceManager.sol";

import "forge-std/Script.sol";
import "@openzeppelin/contracts/proxy/transparent/ProxyAdmin.sol";
import {MachServiceManager} from "../src/core/MachServiceManager.sol";

import {ERC20} from "@openzeppelin/contracts/token/ERC20/ERC20.sol";

contract AltLayerInu is ERC20 {
    constructor() ERC20("AltLayer Inu", "ALTINU") {}

    function mint(address account, uint256 amount) external {
        _mint(account, amount);
    }

    function burn(address account, uint256 amount) external {
        _burn(account, amount);
    }
}

// forge script ./script/CreateAVSRewardsSubmission.s.sol:CreateAVSRewardsSubmission \
//     --private-key $PK \
//     --rpc-url $URL --broadcast --slow
contract CreateAVSRewardsSubmission is Script {
    function run() external {
        vm.startBroadcast();
        IRewardsCoordinator.StrategyAndMultiplier[] memory strategiesAndMultipliers =
            new IRewardsCoordinator.StrategyAndMultiplier[](11);

        strategiesAndMultipliers[0] = IRewardsCoordinator.StrategyAndMultiplier(
            IStrategy(address(0x05037A81BD7B4C9E0F7B430f1F2A22c31a2FD943)), 1 ether
        );
        strategiesAndMultipliers[1] = IRewardsCoordinator.StrategyAndMultiplier(
            IStrategy(address(0x31B6F59e1627cEfC9fA174aD03859fC337666af7)), 1 ether
        );
        strategiesAndMultipliers[2] = IRewardsCoordinator.StrategyAndMultiplier(
            IStrategy(address(0x3A8fBdf9e77DFc25d09741f51d3E181b25d0c4E0)), 1 ether
        );
        strategiesAndMultipliers[3] = IRewardsCoordinator.StrategyAndMultiplier(
            IStrategy(address(0x46281E3B7fDcACdBa44CADf069a94a588Fd4C6Ef)), 1 ether
        );
        strategiesAndMultipliers[4] = IRewardsCoordinator.StrategyAndMultiplier(
            IStrategy(address(0x70EB4D3c164a6B4A5f908D4FBb5a9cAfFb66bAB6)), 1 ether
        );
        strategiesAndMultipliers[5] = IRewardsCoordinator.StrategyAndMultiplier(
            IStrategy(address(0x7673a47463F80c6a3553Db9E54c8cDcd5313d0ac)), 1 ether
        );
        strategiesAndMultipliers[6] = IRewardsCoordinator.StrategyAndMultiplier(
            IStrategy(address(0x7D704507b76571a51d9caE8AdDAbBFd0ba0e63d3)), 1 ether
        );
        strategiesAndMultipliers[7] = IRewardsCoordinator.StrategyAndMultiplier(
            IStrategy(address(0x80528D6e9A2BAbFc766965E0E26d5aB08D9CFaF9)), 1 ether
        );
        strategiesAndMultipliers[8] = IRewardsCoordinator.StrategyAndMultiplier(
            IStrategy(address(0x9281ff96637710Cd9A5CAcce9c6FAD8C9F54631c)), 1 ether
        );
        strategiesAndMultipliers[9] = IRewardsCoordinator.StrategyAndMultiplier(
            IStrategy(address(0xaccc5A86732BE85b5012e8614AF237801636F8e5)), 1 ether
        );
        strategiesAndMultipliers[10] = IRewardsCoordinator.StrategyAndMultiplier(
            IStrategy(address(0xbeaC0eeEeeeeEEeEeEEEEeeEEeEeeeEeeEEBEaC0)), 1 ether
        );

        IRewardsCoordinator.RewardsSubmission[] memory rewardsSubmissions =
            new IRewardsCoordinator.RewardsSubmission[](1) ;

        AltLayerInu rewardToken = new AltLayerInu();
        uint256 amount = 1000000 ether;

        uint32 CALCULATION_INTERVAL_SECONDS = 7 days;
        address serviceManager = 0xAE9a4497dee2540DaF489BeddB0706128a99ec63;
        uint32 startTimestamp = uint32(block.timestamp) - uint32(block.timestamp % CALCULATION_INTERVAL_SECONDS);

        rewardsSubmissions[0] = IRewardsCoordinator.RewardsSubmission({
            strategiesAndMultipliers: strategiesAndMultipliers,
            token: IERC20(rewardToken),
            amount: amount,
            startTimestamp: startTimestamp,
            duration: uint32(28 days)
        });
        rewardToken.mint(0x31Cc55D177824193A5fa2BF34da8AfAfbd366111, amount);
        rewardToken.approve(address(serviceManager), amount);
        IRewardsCoordinator(serviceManager).createAVSRewardsSubmission(rewardsSubmissions);

        vm.stopBroadcast();
    }
}
