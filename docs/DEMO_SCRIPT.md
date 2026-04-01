# SmartMeterChain — Competition Demo Script (5 minutes)

## Setup (before demo starts)
- Backend running on port 3000
- Frontend running on port 3001
- Rust simulator running (generating readings every 30s)
- Open frontend in browser, logged out

---

## Act 1: Introduction (30s)

> "SmartMeterChain is a blockchain-powered smart meter ecosystem built for Indian DISCOMs. It uses Hyperledger Fabric to create tamper-proof records of electricity consumption, automated slab-based billing, and transparent dispute resolution — all compliant with the Electricity Act 2003."

**Show:** Login page with the SmartMeterChain branding.

---

## Act 2: Admin Login + Dashboard (45s)

1. Login as **admin / admin123**
2. Point out the **Dashboard**:
   - "4 active smart meters across Delhi, Mumbai, Bangalore"
   - "Real-time readings being ingested from our Rust-based DLMS simulator"
   - "Today's readings count increasing live" (auto-refreshes every 15s)
   - "Tamper alerts section monitors for anomalies like reverse flow and consumption spikes"
3. Scroll to **Blockchain Secured** badge: "Every reading is SHA-256 hashed and anchored on Hyperledger Fabric"

---

## Act 3: Smart Meters (30s)

1. Click **Meters** in sidebar
2. Show 4 meter cards: "We have domestic, commercial, and industrial meters"
3. Click **Register Meter** button
4. Fill: ID=`SM-CHN-001`, Consumer=`consumer2`, Location=`Chennai, T Nagar`, Type=`domestic`
5. Submit: "New meter registered and ready for data ingestion"

---

## Act 4: Billing — The Core Feature (60s)

1. Click **Billing** in sidebar
2. Click **Generate Bill**
3. Select meter `SM-DEL-001`, Period: `2026-03-01` to `2026-03-31`
4. Submit: "Bill generated using Indian slab-based tariff calculation"
   - "Domestic: 0-100 units @ Rs 3, 100-300 @ Rs 5.50, 300+ @ Rs 8"
   - "Fixed charges added per applicable slab"
5. Point out the bill in the table: Bill ID, units, amount, **pending** status
6. Click **Pay** button: "Consumer pays the bill, status changes to paid"
7. Note the chain link icon: "This bill's hash is recorded on Hyperledger Fabric"

---

## Act 5: Blockchain Verification (45s)

1. Click **Blockchain** in sidebar
2. Copy a Bill ID from the billing page
3. Paste into the verification box and click **Verify**
4. Show result: "Verified on Blockchain — the off-chain hash matches the on-chain record"
5. Explain: "If anyone tampers with the bill amount or units, the hash won't match"
6. Point out the **How Blockchain Verification Works** section

---

## Act 6: Dispute Resolution (30s)

1. Click **Disputes** in sidebar
2. Click **File Dispute**
3. Select a bill, enter reason: "Meter reading seems unusually high for March"
4. Submit: "Dispute filed by consumer"
5. Click **Resolve Dispute** → Enter: "Reading verified against on-chain data, within expected range"
6. Show dispute status changed to **resolved**

---

## Act 7: Role-Based Access (30s)

1. Click **Logout**
2. Login as **regulator / regulator123**
3. Show sidebar: "Regulator has read-only access — can see Dashboard, Meters, Tariffs, Alerts, Blockchain"
4. Show "Read-only access (Regulator)" label in sidebar
5. "Regulators can audit all data but cannot modify anything — as per SERC guidelines"

---

## Act 8: Tariffs + Alerts (30s)

1. Logout, login as admin again
2. Click **Tariffs**: Show slab structure for domestic/commercial/industrial
3. Click **Alerts**: "Our anomaly detection flags reverse flow, consumption spikes over 500%, and voltage irregularities"
4. If any alerts visible: show acknowledge button

---

## Closing (30s)

> "SmartMeterChain demonstrates a complete, production-ready smart meter ecosystem:
> - **Tamper-proof** readings via Hyperledger Fabric
> - **Automated billing** with Indian slab-based tariffs
> - **Real-time monitoring** with anomaly detection
> - **Transparent disputes** with blockchain-backed verification
> - **Multi-role access** for DISCOMs, consumers, and regulators
>
> Built with Go, React, Rust, and Flutter — deployable on any cloud infrastructure."

---

## Backup: If Live Demo Fails
- Show pre-recorded demo video
- Show architecture diagram from `docs/architecture.md`
- Walk through code structure and blockchain chaincode

## Key Numbers to Mention
- 4 smart meters simulated across 3 cities
- 30-second reading intervals (configurable)
- 3 tariff categories with slab-based pricing
- 2% anomaly injection rate for tamper detection testing
- SHA-256 hashing for all readings and bills
- 2-org Fabric network (DiscomMSP + RegulatorMSP)
