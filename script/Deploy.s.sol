// SPDX-License-Identifier: BUSL-1.1
pragma solidity ^0.8.12;
// solhint-disable

import "forge-std/Script.sol";
import "../test/ServiceManager.test.sol";

// anvil --fork-url https://eth-goerli.g.alchemy.com/v2/<api-key>
// FILE='/home/x/z/avs/script/config/deploy.goerli.json' forge script script/Deploy.s.sol:Deploy --rpc-url http://127.0.0.1:8545 --broadcast -vvvv --slow
contract Deploy is Script, AVSDeployer {
    using BN254 for BN254.G1Point;

    // function _parseStakeRegistryParams(
    //     string memory data
    // )
    //     internal
    //     pure
    //     returns (
    //         uint96[] memory minimumStakeForQuourm,
    //         IVoteWeigher.StrategyAndWeightingMultiplier[][]
    //             memory strategyAndWeightingMultipliers
    //     )
    // {
    //     bytes memory stakesConfigsRaw = stdJson.parseRaw(
    //         data,
    //         ".minimumStakes"
    //     );
    //     minimumStakeForQuourm = abi.decode(stakesConfigsRaw, (uint96[]));

    //     bytes memory strategyConfigsRaw = stdJson.parseRaw(
    //         data,
    //         ".strategyWeights"
    //     );
    //     strategyAndWeightingMultipliers = abi.decode(
    //         strategyConfigsRaw,
    //         (IVoteWeigher.StrategyAndWeightingMultiplier[][])
    //     );
    // }

    function _parseRegistryCoordinatorParams(
        string memory data
    )
        internal
        returns (
            IRegistryCoordinator.OperatorSetParam[] memory operatorSetParams,
            address churner,
            address ejector
        )
    {
        bytes memory operatorConfigsRaw = stdJson.parseRaw(
            data,
            ".operatorSetParams"
        );
        operatorSetParams = abi.decode(
            operatorConfigsRaw,
            (IRegistryCoordinator.OperatorSetParam[])
        );

        churner = stdJson.readAddress(data, ".permissions.churner");
        ejector = stdJson.readAddress(data, ".permissions.ejector");
    }

    function run() public {
        string memory data = vm.readFile(vm.envString("FILE"));

        // READ JSON DATA
        bytes32 pk = stdJson.readBytes32(data, ".pk");
        vm.startBroadcast(uint256(pk));

        address admin = stdJson.readAddress(data, ".permissions.owner");

        {
            ProxyAdmin proxyAdmin = new ProxyAdmin();
            EmptyContract emptyContract = new EmptyContract();
            PauserRegistry pauserRegistry;
            {
                address[] memory pausers = new address[](1);
                pausers[0] = admin;
                pauserRegistry = new PauserRegistry(pausers, admin);
            }

            IDelegationManager delegation = IDelegationManager(
                stdJson.readAddress(data, ".delegation")
            );
            Slasher slasher = Slasher(stdJson.readAddress(data, ".slasher"));
            StrategyManager strategyManager = StrategyManager(
                stdJson.readAddress(data, ".strategyManager")
            );

            (
                IRegistryCoordinator.OperatorSetParam[]
                    memory operatorSetParams,
                address churner,
                address ejector
            ) = _parseRegistryCoordinatorParams(data);

            // setup the dummy quorum strategies

            IStakeRegistry.StrategyParams[][]
                memory quorumStrategiesConsideredAndMultipliers = new IStakeRegistry.StrategyParams[][](
                    1
                );

            quorumStrategiesConsideredAndMultipliers[
                0
            ] = new IStakeRegistry.StrategyParams[](2);

            quorumStrategiesConsideredAndMultipliers[0][0] = IStakeRegistry
                .StrategyParams(
                    IStrategy(
                        address(0xB613E78E2068d7489bb66419fB1cfa11275d14da)
                    ),
                    uint96(1070136092289993178)
                );
            quorumStrategiesConsideredAndMultipliers[0][1] = IStakeRegistry
                .StrategyParams(
                    IStrategy(
                        address(0x879944A8cB437a5f8061361f82A6d4EED59070b5)
                    ),
                    1071364636818145808
                );

            // setup the dummy minimum stake for quorum
            uint96[] memory minimumStakeForQuorum = new uint96[](1);
            for (uint256 i; i < minimumStakeForQuorum.length; i++) {
                minimumStakeForQuorum[i] = 1000;
            }

            _deployAVS(
                admin,
                emptyContract,
                proxyAdmin,
                pauserRegistry,
                slasher,
                delegation,
                quorumStrategiesConsideredAndMultipliers,
                minimumStakeForQuorum,
                operatorSetParams
            );

            // IDelegationManager.OperatorDetails
            //     memory operatorDetails = IDelegationManager.OperatorDetails({
            //         earningsReceiver: admin,
            //         delegationApprover: address(0),
            //         stakerOptOutWindowBlocks: 0
            //     });
            // string memory emptyStringForMetadataURI;

            // delegation.registerAsOperator(
            //     operatorDetails,
            //     emptyStringForMetadataURI
            // );
        }

        vm.stopBroadcast();
    }
}
