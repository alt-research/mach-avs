# Contracts Makefile - For MACH AVS tooling

# Default output directory for AVS details
OUTPUT_DIR ?= ../docs/avs-details

# RPC URLs
MAINNET_RPC_URL ?= https://ethereum.publicnode.com
HOLESKY_RPC_URL ?= https://ethereum-holesky.publicnode.com

# Mainnet Registry Coordinator addresses
ALTLAYER_MAINNET_RC = 0x561be1AB42170a19f31645F774e6e3862B2139AA
XTERIO_MAINNET_RC = 0x6e32dD4169825337B54236b9355910b83E51e5a8
DODOCHAIN_MAINNET_RC = 0xCa7f0a65f82230691580F41876161E6282460F29
CYBER_MAINNET_RC = 0x118610D207A32f10F4f7C3a1FEFac5b3327c2bad
GM_MAINNET_RC = 0xDc7e13f1A2Bd0346D741E932CaA341AB6EC21e92
SONEIUM_MAINNET_RC = 0x2c23CF71C023CBA700E379c8E73f040c70211D67

# Holesky Registry Coordinator addresses
ALTLAYER_HOLESKY_RC = 0x1eA7D160d325B289bF981e0D7aB6Bf3261a0FFf2

.PHONY: all scan clean scan-mainnet scan-testnet storage-report

storage-report:
	@echo "📝 Generating storage layout report..."
	@./lib/storage-layout-reporter/report.sh --output=docs/storage --source="src" --exclude="src/interfaces src/error src/core/groth16"
	@echo "✅ Storage report generated in docs/storage directory"

all: build scan

build:
	forge build

# Create output directories
$(OUTPUT_DIR)/mainnet:
	mkdir -p $(OUTPUT_DIR)/mainnet

$(OUTPUT_DIR)/holesky:
	mkdir -p $(OUTPUT_DIR)/holesky

# Scan mainnet AVS deployments
scan-mainnet: $(OUTPUT_DIR)/mainnet
	$(eval MAINNET_BLOCK_NUMBER := $(shell cast block-number --rpc-url $(MAINNET_RPC_URL)))
	@echo "Scanning mainnet AVS deployments using block number: $(MAINNET_BLOCK_NUMBER)"
	@echo "Scanning AltLayer MACH AVS on Mainnet..."
	OUTPUT_DIR=$(OUTPUT_DIR)/mainnet OUTPUT_FILE=altlayer.md REGISTRY_COORDINATOR=$(ALTLAYER_MAINNET_RC) RPC_URL=$(MAINNET_RPC_URL) BLOCK_NUMBER=$(MAINNET_BLOCK_NUMBER) ../scripts/scan-avs.sh
	
	@echo "Scanning Xterio MACH AVS on Mainnet..."
	OUTPUT_DIR=$(OUTPUT_DIR)/mainnet OUTPUT_FILE=xterio.md REGISTRY_COORDINATOR=$(XTERIO_MAINNET_RC) RPC_URL=$(MAINNET_RPC_URL) BLOCK_NUMBER=$(MAINNET_BLOCK_NUMBER) ../scripts/scan-avs.sh
	
	@echo "Scanning DODOchain MACH AVS on Mainnet..."
	OUTPUT_DIR=$(OUTPUT_DIR)/mainnet OUTPUT_FILE=dodochain.md REGISTRY_COORDINATOR=$(DODOCHAIN_MAINNET_RC) RPC_URL=$(MAINNET_RPC_URL) BLOCK_NUMBER=$(MAINNET_BLOCK_NUMBER) ../scripts/scan-avs.sh
	
	@echo "Scanning Cyber MACH AVS on Mainnet..."
	OUTPUT_DIR=$(OUTPUT_DIR)/mainnet OUTPUT_FILE=cyber.md REGISTRY_COORDINATOR=$(CYBER_MAINNET_RC) RPC_URL=$(MAINNET_RPC_URL) BLOCK_NUMBER=$(MAINNET_BLOCK_NUMBER) ../scripts/scan-avs.sh
	
	@echo "Scanning GM Network MACH AVS on Mainnet..."
	OUTPUT_DIR=$(OUTPUT_DIR)/mainnet OUTPUT_FILE=gm-network.md REGISTRY_COORDINATOR=$(GM_MAINNET_RC) RPC_URL=$(MAINNET_RPC_URL) BLOCK_NUMBER=$(MAINNET_BLOCK_NUMBER) ../scripts/scan-avs.sh
	
	@echo "Scanning Soneium MACH AVS on Mainnet..."
	OUTPUT_DIR=$(OUTPUT_DIR)/mainnet OUTPUT_FILE=soneium.md REGISTRY_COORDINATOR=$(SONEIUM_MAINNET_RC) RPC_URL=$(MAINNET_RPC_URL) BLOCK_NUMBER=$(MAINNET_BLOCK_NUMBER) ../scripts/scan-avs.sh

# Scan testnet AVS deployments
scan-testnet: $(OUTPUT_DIR)/holesky
	$(eval HOLESKY_BLOCK_NUMBER := $(shell cast block-number --rpc-url $(HOLESKY_RPC_URL)))
	@echo "Scanning Holesky testnet AVS deployments using block number: $(HOLESKY_BLOCK_NUMBER)"
	@echo "Scanning AltLayer MACH AVS on Holesky..."
	OUTPUT_DIR=$(OUTPUT_DIR)/holesky OUTPUT_FILE=altlayer.md REGISTRY_COORDINATOR=$(ALTLAYER_HOLESKY_RC) RPC_URL=$(HOLESKY_RPC_URL) BLOCK_NUMBER=$(HOLESKY_BLOCK_NUMBER) ../scripts/scan-avs.sh

# Scan all AVS deployments
scan: scan-mainnet scan-testnet
	@echo "Scanning completed for all MACH AVS deployments."

# Clean output files
clean:
	rm -rf $(OUTPUT_DIR)/mainnet/*.md
	rm -rf $(OUTPUT_DIR)/holesky/*.md

# JS-based scanning alternative
scan-js:
	@echo "Scanning using node.js script..."
	REGISTRY_COORDINATOR=$(ALTLAYER_MAINNET_RC) RPC_URL=$(MAINNET_RPC_URL) node script/gen-avs-doc.js $(OUTPUT_DIR)/mainnet/altlayer-js.md

help:
	@echo "Available targets:"
	@echo "  build         - Build contracts using forge"
	@echo "  scan          - Run AVS scanning for all networks"
	@echo "  scan-mainnet  - Run AVS scanning for Mainnet"
	@echo "  scan-testnet  - Run AVS scanning for Holesky testnet"
	@echo "  scan-js       - Run AVS scanning using the JS implementation"
	@echo "  clean         - Remove all scan output files"
	@echo "  storage-report - Generate storage layout report"
	@echo ""
	@echo "Parameters:"
	@echo "  BLOCK_NUMBER  - Specify a block number for consistent data (default: latest)"
	@echo "  Example: make scan-mainnet BLOCK_NUMBER=17500000"
