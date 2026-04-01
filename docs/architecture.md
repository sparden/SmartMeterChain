# SmartMeterChain — Architecture Document

## 1. System Overview

SmartMeterChain is a 5-layer architecture designed for India's smart meter ecosystem:

```
Layer 5: Presentation   (TanStack Start Web + Flutter Mobile)
Layer 4: API Gateway     (Go Fiber REST API with JWT auth)
Layer 3: Business Logic  (Services: Billing, Meter, Alerts, Disputes)
Layer 2: Data Layer      (Hyperledger Fabric + SQLite off-chain cache)
Layer 1: IoT / Ingestion (Rust Simulator / Real Smart Meters via MQTT/HTTP)
```

## 2. Blockchain Architecture

### Network Topology
- **Orderer**: Single Raft orderer (scalable to multi-orderer in production)
- **Organizations**:
  - **DiscomMSP**: Power distribution company (DISCOM)
  - **RegulatorMSP**: State Electricity Regulatory Commission (SERC)
- **Channel**: `smartmeterchannel` — shared ledger between both orgs
- **State DB**: CouchDB (enables rich queries on meter data)

### Chaincode (Smart Contracts)
Written in Go, deployed on Hyperledger Fabric 2.5:

| Contract | Functions | Purpose |
|----------|-----------|---------|
| smartmeter.go | RegisterMeter, SubmitReading, DetectTamper | Meter lifecycle and data ingestion |
| billing.go | GenerateBill, GetBill, VerifyBill | Slab-based billing with hash verification |
| tariff.go | SetTariff, GetTariff, GetAllTariffs | Tariff slab management |
| dispute.go | FileDispute, ResolveDispute | On-chain dispute resolution |

### Data Flow
```
Smart Meter → HTTP POST → Go API → Fabric SDK → Chaincode → World State (CouchDB)
                                  → SQLite Cache (off-chain fast queries)
```

## 3. Tamper Detection Algorithm

Multi-signal anomaly detection:
1. **Reverse Flow**: Reading < Previous Reading (meter bypass)
2. **Spike Detection**: Reading > 500% of previous (illegal tapping)
3. **Large Jump**: Absolute delta > 1000 kWh (hardware malfunction)
4. **Voltage Anomaly**: Voltage outside 207-253V range (± 10% of 230V)

Alerts are generated in real-time and stored both on-chain and off-chain.

## 4. Billing Engine

Slab-based tariff calculation following Indian electricity board standards:

```
For each tariff slab (ordered by slab_start):
  units_in_slab = min(remaining_units, slab_width)
  cost += units_in_slab * rate_per_unit
  Add fixed_charge for highest applicable slab
```

Bill hash = SHA-256(bill_id + units_used + amount) — stored on-chain for verification.

## 5. Security Architecture

- **Authentication**: JWT tokens with 24h expiry
- **Authorization**: Role-based (admin, consumer, regulator)
- **Data Integrity**: SHA-256 hashing of all readings and bills
- **Transport**: TLS for all Fabric communication
- **MSP**: X.509 certificate-based identity for blockchain participants

## 6. Scalability Strategy

- **Horizontal**: Add peer nodes per DISCOM; add organizations for new states
- **Vertical**: CouchDB indexes for rich queries; SQLite for high-speed off-chain reads
- **Federation**: Each state DISCOM runs its own peer; regulator has read access to all
- **Target**: 10M+ meters across multiple DISCOMs (1000 TPS with Fabric 2.5 optimizations)

## 7. Deployment Architecture

```
Docker Compose (Development):
  ├── Fabric Orderer
  ├── Fabric Peer (Discom)
  ├── Fabric Peer (Regulator)
  ├── CouchDB x2
  ├── Go Fiber API
  └── Frontend (static build)

Production (Kubernetes):
  ├── Fabric Operator for peer/orderer management
  ├── API replicas behind load balancer
  ├── PostgreSQL (replacing SQLite for prod)
  └── CloudFront/CDN for frontend
```
