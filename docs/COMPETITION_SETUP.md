# SmartMeterChain — Competition Setup Guide

## Quick Start (Any Linux/Mac with Docker)

```bash
git clone https://github.com/sparden/SmartMeterChain.git
cd SmartMeterChain

# Install deps (pick one based on platform)
bash deploy/setup-oracle-cloud.sh     # Oracle Cloud / Ubuntu VPS
bash deploy/setup-github-codespaces.sh # GitHub Codespaces
bash deploy/setup-gitpod.sh            # Gitpod / Ona

# Start Fabric network (shows 6 blockchain containers)
cd blockchain/network
bash scripts/setup-network.sh
cd ../..

# Start backend (auto-detects Fabric)
cd backend && go run main.go &

# Start frontend
cd ../frontend && npm install && npm run dev &

# Start simulator (generates live meter data)
cd ../simulator && cargo run --release -- --register &
```

## What Judges See

### Fabric Network (6 containers)
```
docker ps
```
Shows:
- `orderer.smartmeterchain.com` — Raft consensus orderer
- `peer0.discom.smartmeterchain.com` — DISCOM organization peer
- `peer0.regulator.smartmeterchain.com` — Regulator organization peer
- `couchdb.discom` — DISCOM state database
- `couchdb.regulator` — Regulator state database
- `cli` — Fabric CLI tools

### Channel Verification
```bash
export PATH=$PATH:$HOME/fabric/bin
export FABRIC_CFG_PATH=blockchain/network
export CORE_PEER_TLS_ENABLED=true
export CORE_PEER_LOCALMSPID="DiscomMSP"
export CORE_PEER_ADDRESS=localhost:7051
export CORE_PEER_TLS_ROOTCERT_FILE=blockchain/organizations/peerOrganizations/discom.smartmeterchain.com/peers/peer0.discom.smartmeterchain.com/tls/ca.crt
export CORE_PEER_MSPCONFIGPATH=blockchain/organizations/peerOrganizations/discom.smartmeterchain.com/users/Admin@discom.smartmeterchain.com/msp

peer channel list
# Output: smartmeterchannel

peer channel getinfo -c smartmeterchannel
# Output: Blockchain info: height:X, currentBlockHash:..., previousBlockHash:...
```

### CouchDB (World State)
Open http://localhost:5984/_utils/ (admin/adminpw) to see the Fabric world state database.

## Reset Everything
```bash
cd blockchain/network
docker compose -f docker-compose-fabric.yaml down -v
rm -rf ../organizations genesis.block
bash scripts/setup-network.sh
```

## Architecture for Judges

```
  Smart Meters (Rust Simulator)
        |
        v
  Go Fiber REST API (port 3000)
        |
   +---------+---------+
   |                   |
   v                   v
SQLite Cache     Hyperledger Fabric
(fast reads)     (tamper-proof records)
                       |
              +--------+--------+
              |                 |
         DiscomMSP        RegulatorMSP
         (peer0)            (peer0)
              |                 |
          CouchDB           CouchDB
```

## Demo Credentials
- Admin: admin / admin123
- Consumer: consumer1 / consumer123
- Regulator: regulator / regulator123

## Key Talking Points
1. **Non-crypto blockchain** — Hyperledger Fabric is permissioned, no cryptocurrency
2. **Indian tariff model** — Domestic/Commercial/Industrial slab-based billing per CERC
3. **Tamper detection** — SHA-256 hash comparison between off-chain and on-chain
4. **Multi-org** — DISCOM and Regulator as separate Fabric organizations
5. **Real-time** — Rust simulator generates DLMS-like meter readings every 30 seconds
6. **Compliance** — Electricity Act 2003, CERC/SERC guidelines
