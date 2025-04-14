#!/bin/bash

# Script to scan AVS contracts and generate configuration details
# Usage: REGISTRY_COORDINATOR=0x1234... RPC_URL=https://... [BLOCK_NUMBER=latest] ./scan-avs.sh

set -e

# Check for required environment variables
if [ -z "$REGISTRY_COORDINATOR" ]; then
    echo "Error: REGISTRY_COORDINATOR environment variable is required"
    exit 1
fi

if [ -z "$RPC_URL" ]; then
    echo "Error: RPC_URL environment variable is required"
    exit 1
fi

# Check for output directory from args or environment variable
OUTPUT_DIR=${OUTPUT_DIR:-"./docs"}
OUTPUT_FILE=${OUTPUT_FILE:-"AVSDetails.md"}

# Create output directory if it doesn't exist
mkdir -p "$OUTPUT_DIR"

# Full path to output file
OUTPUT_MD_FILE="$OUTPUT_DIR/$OUTPUT_FILE"

echo "Output will be written to: $OUTPUT_MD_FILE"

# Define storage slots
IMPLEMENTATION_SLOT="0x360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc"
ADMIN_SLOT="0xb53127684a568b3173ae13b9f8a6016e243e63b6e8ee1178d6a717850b5d6103"

echo "Scanning AVS contracts starting from RegistryCoordinator: $REGISTRY_COORDINATOR"

# Get block number from environment or use latest
if [ -z "$BLOCK_NUMBER" ]; then
    echo "No specific block number provided, using latest block"
    BLOCK_NUMBER=$(cast block-number --rpc-url $RPC_URL)
else
    echo "Using provided block number: $BLOCK_NUMBER"
fi

# Define block parameter for all calls
BLOCK_PARAM="--block $BLOCK_NUMBER"

# Get current block info
BLOCK_TIMESTAMP=$(cast block $BLOCK_NUMBER --rpc-url $RPC_URL | grep timestamp | awk '{print $2}')
CHAIN_ID=$(cast chain-id --rpc-url $RPC_URL)

echo "Network info:"
echo "  Chain ID: $CHAIN_ID"
echo "  Block: $BLOCK_NUMBER"
echo "  Timestamp: $BLOCK_TIMESTAMP"

# Fetch contract addresses from the registry coordinator (with block param)
PAUSER_REGISTRY=$(cast call $REGISTRY_COORDINATOR "pauserRegistry()(address)" --rpc-url $RPC_URL $BLOCK_PARAM)

# Try avs() first, fall back to serviceManager() if it fails
SERVICE_MANAGER=$(cast call $REGISTRY_COORDINATOR "avs()(address)" --rpc-url $RPC_URL $BLOCK_PARAM 2>/dev/null || echo "")
if [ -z "$SERVICE_MANAGER" ]; then
  echo "avs() call failed, trying serviceManager() as fallback..."
  SERVICE_MANAGER=$(cast call $REGISTRY_COORDINATOR "serviceManager()(address)" --rpc-url $RPC_URL $BLOCK_PARAM)
  if [ -z "$SERVICE_MANAGER" ]; then
    echo "Error: Could not retrieve ServiceManager address using either method"
    exit 1
  else
    echo "Successfully retrieved ServiceManager using serviceManager() method"
  fi
else
  echo "Successfully retrieved ServiceManager using avs() method"
fi

STAKE_REGISTRY=$(cast call $REGISTRY_COORDINATOR "stakeRegistry()(address)" --rpc-url $RPC_URL $BLOCK_PARAM)
INDEX_REGISTRY=$(cast call $REGISTRY_COORDINATOR "indexRegistry()(address)" --rpc-url $RPC_URL $BLOCK_PARAM)
BLS_APK_REGISTRY=$(cast call $REGISTRY_COORDINATOR "blsApkRegistry()(address)" --rpc-url $RPC_URL $BLOCK_PARAM)

echo "Found Proxies:"
echo "  PauserRegistry: $PAUSER_REGISTRY"
echo "  ServiceManager: $SERVICE_MANAGER"
echo "  StakeRegistry: $STAKE_REGISTRY"
echo "  IndexRegistry: $INDEX_REGISTRY"
echo "  BLSApkRegistry: $BLS_APK_REGISTRY"

# Fetch implementation addresses for all proxies (with block param)
REGISTRY_COORDINATOR_IMPL=$(cast storage $REGISTRY_COORDINATOR $IMPLEMENTATION_SLOT --rpc-url $RPC_URL $BLOCK_PARAM)
SERVICE_MANAGER_IMPL=$(cast storage $SERVICE_MANAGER $IMPLEMENTATION_SLOT --rpc-url $RPC_URL $BLOCK_PARAM)
STAKE_REGISTRY_IMPL=$(cast storage $STAKE_REGISTRY $IMPLEMENTATION_SLOT --rpc-url $RPC_URL $BLOCK_PARAM)
INDEX_REGISTRY_IMPL=$(cast storage $INDEX_REGISTRY $IMPLEMENTATION_SLOT --rpc-url $RPC_URL $BLOCK_PARAM)
BLS_APK_REGISTRY_IMPL=$(cast storage $BLS_APK_REGISTRY $IMPLEMENTATION_SLOT --rpc-url $RPC_URL $BLOCK_PARAM)

# Extract implementation addresses by taking the last 40 characters (20 bytes) of the bytes32 value
REGISTRY_COORDINATOR_IMPL="0x$(echo $REGISTRY_COORDINATOR_IMPL | cut -c 27-66)"
SERVICE_MANAGER_IMPL="0x$(echo $SERVICE_MANAGER_IMPL | cut -c 27-66)"
STAKE_REGISTRY_IMPL="0x$(echo $STAKE_REGISTRY_IMPL | cut -c 27-66)"
INDEX_REGISTRY_IMPL="0x$(echo $INDEX_REGISTRY_IMPL | cut -c 27-66)"
BLS_APK_REGISTRY_IMPL="0x$(echo $BLS_APK_REGISTRY_IMPL | cut -c 27-66)"

echo "Found implementations:"
echo "  RegistryCoordinator: $REGISTRY_COORDINATOR_IMPL"
echo "  ServiceManager: $SERVICE_MANAGER_IMPL"
echo "  StakeRegistry: $STAKE_REGISTRY_IMPL"
echo "  IndexRegistry: $INDEX_REGISTRY_IMPL"
echo "  BLSApkRegistry: $BLS_APK_REGISTRY_IMPL"

# Get ProxyAdmin (with block param)
PROXY_ADMIN=$(cast storage $REGISTRY_COORDINATOR $ADMIN_SLOT --rpc-url $RPC_URL $BLOCK_PARAM)
PROXY_ADMIN="0x$(echo $PROXY_ADMIN | cut -c 27-66)"
echo "ProxyAdmin: $PROXY_ADMIN"

# Get operator set count (with block param)
OPERATOR_SET_COUNT=$(cast call $REGISTRY_COORDINATOR "quorumCount()(uint8)" --rpc-url $RPC_URL $BLOCK_PARAM)
echo "Operator set count: $OPERATOR_SET_COUNT"

# Determine network name based on chain ID
if [ "$CHAIN_ID" = "1" ]; then
    NETWORK_NAME="Ethereum Mainnet"
    EXPLORER_URL="https://etherscan.io/address/"
elif [ "$CHAIN_ID" = "17000" ]; then
    NETWORK_NAME="Holesky Testnet"
    EXPLORER_URL="https://holesky.etherscan.io/address/"
else
    NETWORK_NAME="Network $CHAIN_ID"
    EXPLORER_URL=""
fi

# Format timestamp as ISO date
BLOCK_DATE=$(date -d @$BLOCK_TIMESTAMP -u +"%Y-%m-%dT%H:%M:%S.000Z")

# Write contracts data for markdown
cat >$OUTPUT_MD_FILE <<EOF
<!-- filepath: $OUTPUT_MD_FILE -->
# AVS Configuration Details

*Generated at block \`$BLOCK_NUMBER\` (\`$BLOCK_DATE\`) on \`$NETWORK_NAME\` (Chain ID: \`$CHAIN_ID\`)*

## Contract Addresses

| Contract | Proxy | Implementation |
| -------- | ------- | -------------- |
| Registry Coordinator | [\`$REGISTRY_COORDINATOR\`](${EXPLORER_URL}$REGISTRY_COORDINATOR) | [\`$REGISTRY_COORDINATOR_IMPL\`](${EXPLORER_URL}$REGISTRY_COORDINATOR_IMPL) |
| Service Manager | [\`$SERVICE_MANAGER\`](${EXPLORER_URL}$SERVICE_MANAGER) | [\`$SERVICE_MANAGER_IMPL\`](${EXPLORER_URL}$SERVICE_MANAGER_IMPL) |
| Stake Registry | [\`$STAKE_REGISTRY\`](${EXPLORER_URL}$STAKE_REGISTRY) | [\`$STAKE_REGISTRY_IMPL\`](${EXPLORER_URL}$STAKE_REGISTRY_IMPL) |
| Index Registry | [\`$INDEX_REGISTRY\`](${EXPLORER_URL}$INDEX_REGISTRY) | [\`$INDEX_REGISTRY_IMPL\`](${EXPLORER_URL}$INDEX_REGISTRY_IMPL) |
| BLS APK Registry | [\`$BLS_APK_REGISTRY\`](${EXPLORER_URL}$BLS_APK_REGISTRY) | [\`$BLS_APK_REGISTRY_IMPL\`](${EXPLORER_URL}$BLS_APK_REGISTRY_IMPL) |
| Pauser Registry | [\`$PAUSER_REGISTRY\`](${EXPLORER_URL}$PAUSER_REGISTRY) | - |
| Proxy Admin | [\`$PROXY_ADMIN\`](${EXPLORER_URL}$PROXY_ADMIN) | - |

## Operator Set Configuration

**Total operator sets:** \`$OPERATOR_SET_COUNT\`

EOF

# For each operator set, collect strategy information for markdown
for q in $(seq 0 $(($OPERATOR_SET_COUNT - 1))); do
    echo "Processing operator set $q for markdown"

    MIN_STAKE=$(cast call $STAKE_REGISTRY "minimumStakeForQuorum(uint8)(uint96)" $q --rpc-url $RPC_URL $BLOCK_PARAM | tr -d ' ')
    STRATEGY_COUNT=$(cast call $STAKE_REGISTRY "strategyParamsLength(uint8)(uint256)" $q --rpc-url $RPC_URL $BLOCK_PARAM)

    # Simply display the raw value instead of trying to format with bc
    echo "  Minimum stake: $MIN_STAKE"
    echo "  Strategy count: $STRATEGY_COUNT"

    # First write the operator set header
    cat >>$OUTPUT_MD_FILE <<EOF

### Operator Set $q

**Minimum stake**: \`$MIN_STAKE\` Wei

| Strategy | Multiplier | Underlying Token |
| -------- | ---------- | ---------------- |
EOF

    # For each strategy in the operator set
    for i in $(seq 0 $(($STRATEGY_COUNT - 1))); do
        # Use calldata encoding to get strategyParamsByIndex (with block param)
        STRATEGY_PARAMS=$(cast call $STAKE_REGISTRY "strategyParamsByIndex(uint8,uint256)((address,uint96))" $q $i --rpc-url $RPC_URL $BLOCK_PARAM)

        # Parse the result - format is (address,uint96)
        STRATEGY_ADDRESS=$(echo $STRATEGY_PARAMS | cut -d ',' -f1 | tr -d '(' | tr -d ' ')
        STRATEGY_MULTIPLIER=$(echo $STRATEGY_PARAMS | cut -d ',' -f2 | tr -d ')' | tr -d ' ')

        # Try to get the underlying token for the strategy (with block param)
        UNDERLYING_TOKEN="0x0000000000000000000000000000000000000000"
        if [ "$STRATEGY_ADDRESS" != "0x0000000000000000000000000000000000000000" ]; then
            # Use silent failure to handle strategies that don't have underlyingToken() function
            UNDERLYING_TOKEN_RESULT=$(cast call $STRATEGY_ADDRESS "underlyingToken()(address)" --rpc-url $RPC_URL $BLOCK_PARAM 2>/dev/null || echo "0x0000000000000000000000000000000000000000")
            if [ "$UNDERLYING_TOKEN_RESULT" != "" ]; then
                UNDERLYING_TOKEN=$UNDERLYING_TOKEN_RESULT
            fi
        fi

        # Format multiplier more cleanly - without the [1e18] notation
        MULTIPLIER_STR="\`$STRATEGY_MULTIPLIER\`"

        # Write strategy directly to the file without using STRATEGY_INFO variable
        if [ "$UNDERLYING_TOKEN" = "0x0000000000000000000000000000000000000000" ]; then
            echo "| [\`$STRATEGY_ADDRESS\`](${EXPLORER_URL}$STRATEGY_ADDRESS) | $MULTIPLIER_STR | None |" >> "$OUTPUT_MD_FILE"
        else
            echo "| [\`$STRATEGY_ADDRESS\`](${EXPLORER_URL}$STRATEGY_ADDRESS) | $MULTIPLIER_STR | [\`$UNDERLYING_TOKEN\`](${EXPLORER_URL}$UNDERLYING_TOKEN) |" >> "$OUTPUT_MD_FILE"
        fi
    done

    # Add an empty line after the operator set
    echo "" >> "$OUTPUT_MD_FILE"
done

echo "Markdown file generated successfully at $OUTPUT_MD_FILE"
