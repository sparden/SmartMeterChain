#!/bin/bash
# SmartMeterChain — Chaincode Deployment Script
# Packages, installs, approves, and commits the smartmeter chaincode
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
NETWORK_DIR="$(dirname "$SCRIPT_DIR")"
CHAINCODE_DIR="${NETWORK_DIR}/../chaincode/smartmeter"
CHANNEL_NAME="smartmeterchannel"
CC_NAME="smartmeter"
CC_VERSION="1.0"
CC_SEQUENCE="1"
CC_LABEL="${CC_NAME}_${CC_VERSION}"
ORDERER_CA="${NETWORK_DIR}/../organizations/ordererOrganizations/smartmeterchain.com/orderers/orderer.smartmeterchain.com/tls/ca.crt"

export FABRIC_CFG_PATH="${NETWORK_DIR}"

echo "========================================="
echo " SmartMeterChain Chaincode Deployment"
echo " Name: ${CC_NAME} v${CC_VERSION}"
echo "========================================="

# Step 1: Package chaincode
echo "[1/5] Packaging chaincode..."
cd "${CHAINCODE_DIR}" && GO111MODULE=on go mod vendor && cd -
peer lifecycle chaincode package "${NETWORK_DIR}/${CC_NAME}.tar.gz" \
    --path "${CHAINCODE_DIR}" \
    --lang golang \
    --label "${CC_LABEL}"
echo "  -> Package created."

# Step 2: Install on Discom peer
echo "[2/5] Installing on Discom peer..."
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="DiscomMSP"
export CORE_PEER_TLS_ROOTCERT_FILE="${NETWORK_DIR}/../organizations/peerOrganizations/discom.smartmeterchain.com/peers/peer0.discom.smartmeterchain.com/tls/ca.crt"
export CORE_PEER_MSPCONFIGPATH="${NETWORK_DIR}/../organizations/peerOrganizations/discom.smartmeterchain.com/users/Admin@discom.smartmeterchain.com/msp"
export CORE_PEER_ADDRESS=localhost:7051

peer lifecycle chaincode install "${NETWORK_DIR}/${CC_NAME}.tar.gz"
echo "  -> Installed on Discom."

# Get package ID
CC_PACKAGE_ID=$(peer lifecycle chaincode queryinstalled --output json | \
    python3 -c "import sys,json; pkgs=json.load(sys.stdin)['installed_chaincodes']; print([p['package_id'] for p in pkgs if p['label']=='${CC_LABEL}'][0])")
echo "  -> Package ID: ${CC_PACKAGE_ID}"

# Step 3: Install on Regulator peer
echo "[3/5] Installing on Regulator peer..."
export CORE_PEER_LOCALMSPID="RegulatorMSP"
export CORE_PEER_TLS_ROOTCERT_FILE="${NETWORK_DIR}/../organizations/peerOrganizations/regulator.smartmeterchain.com/peers/peer0.regulator.smartmeterchain.com/tls/ca.crt"
export CORE_PEER_MSPCONFIGPATH="${NETWORK_DIR}/../organizations/peerOrganizations/regulator.smartmeterchain.com/users/Admin@regulator.smartmeterchain.com/msp"
export CORE_PEER_ADDRESS=localhost:9051

peer lifecycle chaincode install "${NETWORK_DIR}/${CC_NAME}.tar.gz"
echo "  -> Installed on Regulator."

# Step 4: Approve for both orgs
echo "[4/5] Approving chaincode..."

# Approve for Regulator (already set as current peer)
peer lifecycle chaincode approveformyorg \
    -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.smartmeterchain.com \
    --tls --cafile "${ORDERER_CA}" \
    --channelID "${CHANNEL_NAME}" \
    --name "${CC_NAME}" \
    --version "${CC_VERSION}" \
    --package-id "${CC_PACKAGE_ID}" \
    --sequence ${CC_SEQUENCE}
echo "  -> Approved by Regulator."

# Switch to Discom and approve
export CORE_PEER_LOCALMSPID="DiscomMSP"
export CORE_PEER_TLS_ROOTCERT_FILE="${NETWORK_DIR}/../organizations/peerOrganizations/discom.smartmeterchain.com/peers/peer0.discom.smartmeterchain.com/tls/ca.crt"
export CORE_PEER_MSPCONFIGPATH="${NETWORK_DIR}/../organizations/peerOrganizations/discom.smartmeterchain.com/users/Admin@discom.smartmeterchain.com/msp"
export CORE_PEER_ADDRESS=localhost:7051

peer lifecycle chaincode approveformyorg \
    -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.smartmeterchain.com \
    --tls --cafile "${ORDERER_CA}" \
    --channelID "${CHANNEL_NAME}" \
    --name "${CC_NAME}" \
    --version "${CC_VERSION}" \
    --package-id "${CC_PACKAGE_ID}" \
    --sequence ${CC_SEQUENCE}
echo "  -> Approved by Discom."

# Step 5: Commit chaincode
echo "[5/5] Committing chaincode definition..."
peer lifecycle chaincode commit \
    -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.smartmeterchain.com \
    --tls --cafile "${ORDERER_CA}" \
    --channelID "${CHANNEL_NAME}" \
    --name "${CC_NAME}" \
    --version "${CC_VERSION}" \
    --sequence ${CC_SEQUENCE} \
    --peerAddresses localhost:7051 \
    --tlsRootCertFiles "${NETWORK_DIR}/../organizations/peerOrganizations/discom.smartmeterchain.com/peers/peer0.discom.smartmeterchain.com/tls/ca.crt" \
    --peerAddresses localhost:9051 \
    --tlsRootCertFiles "${NETWORK_DIR}/../organizations/peerOrganizations/regulator.smartmeterchain.com/peers/peer0.regulator.smartmeterchain.com/tls/ca.crt"

echo ""
echo "========================================="
echo " Chaincode '${CC_NAME}' v${CC_VERSION} COMMITTED!"
echo " Channel: ${CHANNEL_NAME}"
echo "========================================="

# Verify
echo ""
echo "Verifying committed chaincode..."
peer lifecycle chaincode querycommitted --channelID "${CHANNEL_NAME}" --name "${CC_NAME}"
echo ""
echo "Testing with InitLedger..."
peer chaincode invoke \
    -o localhost:7050 \
    --ordererTLSHostnameOverride orderer.smartmeterchain.com \
    --tls --cafile "${ORDERER_CA}" \
    -C "${CHANNEL_NAME}" \
    -n "${CC_NAME}" \
    --peerAddresses localhost:7051 \
    --tlsRootCertFiles "${NETWORK_DIR}/../organizations/peerOrganizations/discom.smartmeterchain.com/peers/peer0.discom.smartmeterchain.com/tls/ca.crt" \
    --peerAddresses localhost:9051 \
    --tlsRootCertFiles "${NETWORK_DIR}/../organizations/peerOrganizations/regulator.smartmeterchain.com/peers/peer0.regulator.smartmeterchain.com/tls/ca.crt" \
    -c '{"function":"InitLedger","Args":[]}'

echo "Done! Chaincode is ready."
