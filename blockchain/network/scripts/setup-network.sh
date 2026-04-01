#!/bin/bash
# SmartMeterChain — Network Setup Script
# Generates crypto material, creates channel, and joins peers
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
NETWORK_DIR="$(dirname "$SCRIPT_DIR")"
CHANNEL_NAME="smartmeterchannel"
FABRIC_CFG_PATH="${NETWORK_DIR}"

export FABRIC_CFG_PATH

echo "========================================="
echo " SmartMeterChain Network Setup"
echo "========================================="

# Step 1: Generate crypto material
echo "[1/5] Generating crypto material..."
if [ -d "${NETWORK_DIR}/../organizations" ]; then
    echo "  -> Crypto material already exists, skipping."
else
    cryptogen generate --config="${NETWORK_DIR}/crypto-config.yaml" \
        --output="${NETWORK_DIR}/../organizations"
    echo "  -> Done."
fi

# Step 2: Start Docker containers
echo "[2/5] Starting Fabric network containers..."
docker compose -f "${NETWORK_DIR}/docker-compose-fabric.yaml" up -d
echo "  -> Waiting 5s for containers to initialize..."
sleep 5

# Step 3: Create channel
echo "[3/5] Creating channel '${CHANNEL_NAME}'..."
export OSN_TLS_CA_ROOT_CERT="${NETWORK_DIR}/../organizations/ordererOrganizations/smartmeterchain.com/orderers/orderer.smartmeterchain.com/tls/ca.crt"
export ADMIN_TLS_SIGN_CERT="${NETWORK_DIR}/../organizations/ordererOrganizations/smartmeterchain.com/orderers/orderer.smartmeterchain.com/tls/server.crt"
export ADMIN_TLS_PRIVATE_KEY="${NETWORK_DIR}/../organizations/ordererOrganizations/smartmeterchain.com/orderers/orderer.smartmeterchain.com/tls/server.key"

configtxgen -profile SmartMeterOrdererGenesis \
    -outputBlock "${NETWORK_DIR}/genesis.block" \
    -channelID "${CHANNEL_NAME}"

osnadmin channel join --channelID "${CHANNEL_NAME}" \
    --config-block "${NETWORK_DIR}/genesis.block" \
    -o localhost:7053 \
    --ca-file "${OSN_TLS_CA_ROOT_CERT}" \
    --client-cert "${ADMIN_TLS_SIGN_CERT}" \
    --client-key "${ADMIN_TLS_PRIVATE_KEY}"

echo "  -> Channel created."

# Step 4: Join Discom peer
echo "[4/5] Joining Discom peer to channel..."
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="DiscomMSP"
export CORE_PEER_TLS_ROOTCERT_FILE="${NETWORK_DIR}/../organizations/peerOrganizations/discom.smartmeterchain.com/peers/peer0.discom.smartmeterchain.com/tls/ca.crt"
export CORE_PEER_MSPCONFIGPATH="${NETWORK_DIR}/../organizations/peerOrganizations/discom.smartmeterchain.com/users/Admin@discom.smartmeterchain.com/msp"
export CORE_PEER_ADDRESS=localhost:7051

peer channel join -b "${NETWORK_DIR}/genesis.block"
echo "  -> Discom peer joined."

# Step 5: Join Regulator peer
echo "[5/5] Joining Regulator peer to channel..."
export CORE_PEER_LOCALMSPID="RegulatorMSP"
export CORE_PEER_TLS_ROOTCERT_FILE="${NETWORK_DIR}/../organizations/peerOrganizations/regulator.smartmeterchain.com/peers/peer0.regulator.smartmeterchain.com/tls/ca.crt"
export CORE_PEER_MSPCONFIGPATH="${NETWORK_DIR}/../organizations/peerOrganizations/regulator.smartmeterchain.com/users/Admin@regulator.smartmeterchain.com/msp"
export CORE_PEER_ADDRESS=localhost:9051

peer channel join -b "${NETWORK_DIR}/genesis.block"
echo "  -> Regulator peer joined."

echo ""
echo "========================================="
echo " Network is UP!"
echo " Channel: ${CHANNEL_NAME}"
echo " Orgs: DiscomMSP, RegulatorMSP"
echo "========================================="
