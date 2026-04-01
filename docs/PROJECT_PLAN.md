# Project Plan: SmartMeterChain — Blockchain-Powered Smart Meter Ecosystem

## Context
Full prototype for **Blockchain India Challenge (CDAC/MeitY)** — Use Case #9: Power/Smart Meter Ecosystem. Non-crypto, permissioned blockchain solution.

**Project location**: `/Users/akashgupta/SmartMeterChain/`

---

## Final Tech Stack

| Component | Technology |
|---|---|
| Blockchain | Hyperledger Fabric 2.5 + Go chaincode |
| Backend API | Go (Fiber v2) |
| Frontend | TanStack Start (Vite + TanStack Router + TanStack Query + Tailwind v4) |
| Mobile | Flutter + BLoC + Dio |
| Simulator | Rust (tokio + reqwest) |
| Database | SQLite (via GORM) for off-chain cache |
| Deploy | Docker Compose (Fabric + API), Cloudflare Pages (frontend) |

---

## Project Structure

```
SmartMeterChain/
├── README.md
├── docker-compose.yaml
├── docs/
│   ├── architecture.md
│   ├── business-model.md
│   └── compliance.md
│
├── blockchain/                        # Hyperledger Fabric
│   ├── network/
│   │   ├── docker-compose-fabric.yaml
│   │   ├── configtx.yaml
│   │   ├── crypto-config.yaml
│   │   └── scripts/
│   │       ├── setup-network.sh
│   │       └── deploy-chaincode.sh
│   └── chaincode/smartmeter/          # Go chaincode
│       ├── go.mod
│       ├── main.go                    # Chaincode entry, Init + Invoke router
│       ├── smartmeter.go              # RegisterMeter, SubmitReading, GetReadings, DetectTamper
│       ├── billing.go                 # GenerateBill, auto-calc from readings + tariff slabs
│       ├── tariff.go                  # SetTariff, GetTariff (domestic/commercial/industrial slabs)
│       ├── dispute.go                 # FileDispute, ResolveDispute, GetDisputes
│       └── models.go                  # MeterReading, Bill, Consumer, Tariff, Dispute structs
│
├── backend/                           # Go Fiber API
│   ├── go.mod
│   ├── go.sum
│   ├── main.go                        # Fiber server entry
│   ├── config/
│   │   ├── config.go                  # App config (env vars)
│   │   └── fabric.go                  # Fabric gateway connection profile
│   ├── handlers/
│   │   ├── meter.go                   # Meter CRUD + data ingestion
│   │   ├── billing.go                 # Bill generation + history
│   │   ├── consumer.go                # Consumer management
│   │   ├── tariff.go                  # Tariff CRUD
│   │   ├── dispute.go                 # Dispute filing + resolution
│   │   └── dashboard.go              # Analytics + aggregation
│   ├── services/
│   │   ├── fabric.go                  # Fabric Gateway SDK wrapper
│   │   ├── meter.go                   # Meter business logic
│   │   ├── billing.go                 # Billing calculations
│   │   └── alert.go                   # Tamper alerts + notifications
│   ├── models/
│   │   └── models.go                  # GORM models (User, Meter, ReadingCache, etc.)
│   ├── middleware/
│   │   └── auth.go                    # JWT auth middleware
│   ├── database/
│   │   └── db.go                      # SQLite + GORM setup
│   └── utils/
│       ├── hash.go                    # SHA-256 hashing
│       └── response.go               # Standard API responses
│
├── frontend/                          # TanStack Start
│   ├── package.json
│   ├── app.config.ts                  # TanStack Start config
│   ├── tailwind.config.ts
│   ├── tsconfig.json
│   ├── app/
│   │   ├── client.tsx                 # Client entry
│   │   ├── router.tsx                 # TanStack Router config
│   │   ├── routeTree.gen.ts           # Auto-generated route tree
│   │   ├── routes/
│   │   │   ├── __root.tsx             # Root layout (sidebar + header)
│   │   │   ├── index.tsx              # Login page
│   │   │   ├── dashboard.tsx          # Main dashboard
│   │   │   ├── meters.tsx             # Meter management
│   │   │   ├── billing.tsx            # Billing overview
│   │   │   ├── consumers.tsx          # Consumer management
│   │   │   ├── disputes.tsx           # Dispute resolution
│   │   │   ├── tariffs.tsx            # Tariff config
│   │   │   ├── regulator.tsx          # Regulator read-only view
│   │   │   └── explorer.tsx           # Blockchain tx explorer
│   │   ├── components/
│   │   │   ├── Sidebar.tsx
│   │   │   ├── Header.tsx
│   │   │   ├── StatCard.tsx
│   │   │   ├── ConsumptionChart.tsx   # Recharts
│   │   │   ├── BillingTable.tsx
│   │   │   ├── MeterCard.tsx
│   │   │   ├── TamperAlert.tsx
│   │   │   └── BlockchainBadge.tsx
│   │   ├── lib/
│   │   │   ├── api.ts                 # TanStack Query hooks + fetch client
│   │   │   └── types.ts              # Shared TypeScript types
│   │   └── styles/
│   │       └── globals.css            # Tailwind imports
│   └── public/
│       └── favicon.ico
│
├── mobile/                            # Flutter consumer app
│   ├── pubspec.yaml
│   ├── lib/
│   │   ├── main.dart
│   │   ├── models/
│   │   │   ├── meter_reading.dart
│   │   │   ├── bill.dart
│   │   │   └── consumer.dart
│   │   ├── services/
│   │   │   ├── api_service.dart       # Dio REST client
│   │   │   └── auth_service.dart
│   │   ├── screens/
│   │   │   ├── login_screen.dart
│   │   │   ├── home_screen.dart       # Consumption summary + current bill
│   │   │   ├── bills_screen.dart      # Bill history
│   │   │   ├── consumption_screen.dart # Interactive chart
│   │   │   ├── dispute_screen.dart    # File + track disputes
│   │   │   └── verify_screen.dart     # Verify bill on blockchain
│   │   ├── widgets/
│   │   │   ├── consumption_chart.dart
│   │   │   ├── bill_card.dart
│   │   │   └── blockchain_badge.dart  # "Verified on Blockchain" badge
│   │   └── bloc/
│   │       ├── auth/
│   │       │   ├── auth_bloc.dart
│   │       │   ├── auth_event.dart
│   │       │   └── auth_state.dart
│   │       ├── meter/
│   │       │   ├── meter_bloc.dart
│   │       │   ├── meter_event.dart
│   │       │   └── meter_state.dart
│   │       └── billing/
│   │           ├── billing_bloc.dart
│   │           ├── billing_event.dart
│   │           └── billing_state.dart
│   └── test/
│
└── simulator/                         # Rust meter simulator
    ├── Cargo.toml
    ├── src/
    │   ├── main.rs                    # CLI entry + async runtime
    │   ├── meter.rs                   # DLMS-like reading generator
    │   ├── config.rs                  # Meter profile configs
    │   └── api_client.rs             # POST readings to Go backend
    └── config.yaml                    # Meter IDs, intervals, ranges
```

---

## Implementation Steps (35 steps, 5 phases)

### Phase 1: Blockchain Layer (Steps 1-6)
1. Create full project directory structure + README
2. Write Fabric network configs (docker-compose, configtx, crypto-config for 2 orgs: UtilityOrg, RegulatorOrg)
3. Write `models.go` — MeterReading, Bill, Consumer, Tariff, Dispute structs with JSON tags
4. Write `smartmeter.go` — RegisterMeter, SubmitReading (hash + anchor), GetReadingsByMeter, DetectTamper
5. Write `tariff.go` (slab-based Indian tariff: domestic/commercial/industrial) + `billing.go` (GenerateBill from readings + tariff, GetBillHistory) + `dispute.go` (FileDispute, ResolveDispute)
6. Write `main.go` (chaincode Init/Invoke router) + setup/deploy scripts

### Phase 2: Go Backend API (Steps 7-14)
7. Initialize Go module, install Fiber v2 + GORM + Fabric Gateway SDK
8. Write models (GORM) + database setup (SQLite)
9. Write Fabric gateway service (connect to network, submit/evaluate transactions)
10. Write meter handlers (POST /meters, POST /meters/:id/readings, GET /meters/:id/readings, GET /meters/:id/tamper-check)
11. Write billing handlers (POST /billing/generate, GET /billing/:consumerId, GET /billing/:billId/verify)
12. Write consumer + tariff + dispute handlers
13. Write dashboard handler (aggregated stats, consumption trends, alert counts)
14. Write JWT auth middleware + hash utilities + config

### Phase 3: TanStack Start Frontend (Steps 15-24)
15. Initialize TanStack Start project (npm create @tanstack/start), add Tailwind v4 + Recharts
16. Root layout (__root.tsx) — sidebar navigation + header with role badge (Utility/Regulator)
17. Login page with role selection
18. Dashboard page — 4 stat cards (active meters, monthly consumption, revenue, tamper alerts) + consumption line chart + recent alerts feed
19. Meters page — table with status, register new meter form, detail view
20. Billing page — table with bills, generate bill action, "Verify on Blockchain" button showing on-chain hash match
21. Consumer management page — CRUD table
22. Dispute resolution page — file dispute form, dispute list with status, resolve action
23. Tariff configuration page — slab editor (add/edit/delete tariff slabs per category)
24. Regulator view (read-only dashboard) + Blockchain explorer page (list transactions, show raw on-chain data)

### Phase 4: Flutter Mobile App (Steps 25-30)
25. Initialize Flutter project, add dependencies (dio, flutter_bloc, fl_chart, equatable)
26. Write models (MeterReading, Bill, Consumer) + API service (Dio client)
27. Auth flow — login screen + AuthBloc
28. Home screen — consumption summary card, current bill card, quick action buttons
29. Bills screen — list with BillCards, detail view, payment status indicator
30. Consumption chart screen (fl_chart) + Dispute screen + Verify screen (blockchain hash verification)

### Phase 5: Simulator + Documentation (Steps 31-35)
31. Rust simulator — async meter reading generator with realistic consumption patterns (day/night/peak profiles)
32. Config system — multiple meter profiles (residential, commercial, industrial), configurable intervals
33. README — complete setup guide (prerequisites, docker, run commands, demo walkthrough)
34. Architecture document — system diagram, data flow, blockchain interaction model
35. Business model + compliance doc — Indian Electricity Act 2003 alignment, CERC/SERC compliance notes, proliferation strategy across States/UTs

---

## Key Design Decisions
- **Non-crypto**: No tokens/cryptocurrency — pure permissioned Fabric with Go chaincode
- **Go everywhere on backend**: Chaincode + API in Go = shared models, single language expertise
- **Off-chain + on-chain**: Readings stored in SQLite; SHA-256 hashes anchored on Fabric for integrity
- **Indian tariff model**: Slab-based (0-100 units @ Rs 3, 101-200 @ Rs 5, 200+ @ Rs 7 etc.) per category
- **Tamper detection**: Hash of submitted reading vs on-chain anchored hash — mismatch = tamper alert
- **Batch anchoring**: Hourly aggregates written to blockchain to manage throughput
- **Rust simulator**: DLMS-like profiles with day/night/peak patterns — deep-tech showcase for jury

## Verification / Demo Flow
1. `docker-compose up` — starts Fabric network (2 peers, 1 orderer, 2 CAs)
2. `./blockchain/network/scripts/deploy-chaincode.sh` — deploys smartmeter chaincode
3. `cd backend && go run main.go` — starts Go Fiber API on :3000
4. `cd frontend && npm run dev` — starts TanStack Start dashboard on :3001
5. `cd simulator && cargo run` — generates meter readings every 5 seconds, posts to API
6. Open dashboard — see live data flowing, auto-generated bills, tamper detection alerts
7. Run Flutter app — consumer views bills, consumption charts, verifies on blockchain
