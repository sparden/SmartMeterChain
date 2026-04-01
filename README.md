# SmartMeterChain

**Blockchain-Powered Smart Meter Ecosystem** for the **Blockchain India Challenge (C-DAC/MeitY)** — Use Case #9: Power/Smart Meter Ecosystem.

A non-crypto, permissioned blockchain solution that brings transparency, tamper-resistance, and trust to India's smart meter infrastructure.

## Problem Statement

India's power distribution ecosystem suffers from:
- Meter tampering and energy theft (~INR 1.5 lakh crore annual losses)
- Opaque billing processes with no consumer verifiability
- Slow dispute resolution (30-90 days average)
- Siloed data across DISCOMs, regulators, and consumers
- Lack of real-time monitoring for anomaly detection

## Solution

SmartMeterChain uses **Hyperledger Fabric** to create an immutable, shared ledger across all stakeholders (DISCOMs, regulators, consumers) with:

1. **Immutable Meter Readings** — Every reading is hashed and recorded on-chain
2. **Transparent Billing** — Bills generated from on-chain data using slab-based tariffs
3. **Bill Verification** — Any consumer can verify their bill against blockchain records
4. **Tamper Detection** — AI-based anomaly detection for spike, reverse flow, and voltage anomalies
5. **Dispute Resolution** — On-chain dispute filing with complete audit trail
6. **Real-time Dashboard** — Live monitoring for DISCOMs and regulators

## Architecture

```
+-------------------+     +------------------+     +-------------------+
|   Smart Meters    | --> |   Rust Simulator  | --> |   Go Fiber API    |
|   (IoT Devices)   |     |   (Data Ingest)   |     |   (Backend)       |
+-------------------+     +------------------+     +-------------------+
                                                           |
                                              +------------+------------+
                                              |                         |
                                    +---------v---------+    +----------v--------+
                                    | Hyperledger Fabric |    | SQLite (Off-chain |
                                    | (On-chain ledger)  |    |  cache via GORM)  |
                                    +-------------------+    +-------------------+
                                              |
                              +---------------+---------------+
                              |                               |
                    +---------v---------+           +---------v---------+
                    | TanStack Start    |           | Flutter Mobile    |
                    | (Web Dashboard)   |           | (Consumer App)    |
                    +-------------------+           +-------------------+
```

## Tech Stack

| Component | Technology |
|-----------|-----------|
| Blockchain | Hyperledger Fabric 2.5 + Go Chaincode |
| Backend API | Go (Fiber v2) + GORM + SQLite |
| Web Frontend | TanStack Start + React 19 + Tailwind v4 + Recharts |
| Mobile App | Flutter + BLoC + Dio |
| Simulator | Rust (Tokio + Reqwest) |
| Deployment | Docker Compose |

## Project Structure

```
SmartMeterChain/
├── blockchain/          # Hyperledger Fabric network + Go chaincode
│   ├── chaincode/       # Smart contracts (meter, billing, tariff, dispute)
│   └── network/         # Docker compose, crypto config, setup scripts
├── backend/             # Go Fiber REST API
├── frontend/            # TanStack Start web dashboard
├── mobile/              # Flutter consumer mobile app
├── simulator/           # Rust smart meter data simulator
├── docs/                # Architecture, business model, compliance docs
└── docker-compose.yaml  # Root orchestration
```

## Quick Start

### Prerequisites
- Docker & Docker Compose
- Go 1.22+
- Node.js 20+
- Rust 1.75+
- Flutter 3.19+

### 1. Start Fabric Network
```bash
cd blockchain/network
chmod +x scripts/*.sh
./scripts/setup-network.sh
./scripts/deploy-chaincode.sh
```

### 2. Start Backend API
```bash
cd backend
go mod tidy
go run main.go
# API runs on http://localhost:3000
```

### 3. Start Frontend
```bash
cd frontend
npm install
npm run dev
# Dashboard at http://localhost:3001
```

### 4. Run Simulator
```bash
cd simulator
cargo run -- --register --interval 30
```

### 5. Flutter Mobile
```bash
cd mobile
flutter pub get
flutter run
```

## Demo Credentials

| Role | Username | Password |
|------|----------|----------|
| Admin (DISCOM) | admin | admin123 |
| Consumer | consumer1 | consumer123 |
| Regulator | regulator | regulator123 |

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | /api/v1/auth/login | Authenticate user |
| GET | /api/v1/meters | List all meters |
| POST | /api/v1/meters/readings | Submit meter reading |
| GET | /api/v1/bills | Get bills |
| POST | /api/v1/bills/generate | Generate bill from readings |
| GET | /api/v1/verify/bill/:id | Verify bill on blockchain |
| POST | /api/v1/disputes | File billing dispute |
| GET | /api/v1/dashboard/stats | Dashboard analytics |

## Key Features for Challenge Evaluation

- **Non-Crypto Blockchain**: Permissioned Hyperledger Fabric (no tokens/cryptocurrency)
- **Deep Tech**: AI-based tamper detection, SHA-256 hash verification, slab-based billing engine
- **Multi-Stakeholder**: DISCOM, Regulator, and Consumer roles with RBAC
- **Scalable Architecture**: Modular with API integration support, CouchDB state database
- **Indian Context**: Domestic/Commercial/Industrial tariff slabs, INR billing, Indian smart meter standards
- **Complete Documentation**: Architecture, business model, compliance, and deployment guides

## License

Developed for the Blockchain India Challenge. Joint IPR with MeitY as per challenge terms.
