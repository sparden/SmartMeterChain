# SmartMeterChain — Compliance & Regulatory Document

## 1. Blockchain India Challenge Compliance

### Non-Crypto Requirement
- Solution uses **Hyperledger Fabric** — a permissioned, non-crypto blockchain
- No tokens, cryptocurrencies, or mining involved
- All value transfers are recorded as INR amounts in traditional banking systems

### Open Source Compliance
- **Hyperledger Fabric** — Apache 2.0 License
- **Go Fiber** — MIT License
- **React/TanStack** — MIT License
- **Flutter** — BSD 3-Clause License
- **Rust dependencies** — MIT/Apache 2.0
- All licenses allow commercial use and modification

### Original Development
- Core chaincode (billing, tamper detection, dispute resolution) — 100% original
- Backend API and business logic — 100% original
- Frontend dashboard and mobile app — 100% original
- Only standard libraries and frameworks used as foundations

## 2. Indian Regulatory Compliance

### Electricity Act, 2003
- Compliant with Section 55 (meters and billing)
- Supports CERC/SERC tariff structures
- Enables consumer grievance redressal (Section 42)

### Smart Meter Standards
- **IS 16444**: Indian Standard for smart meters
- **IEC 62056**: International metering standard (DLMS/COSEM compatible)
- **IS 15959**: Data exchange standards for electricity meters
- API designed to ingest data from any IS 16444-compliant meter

### Data Protection
- **IT Act, 2000**: Personal data handling compliance
- **DPDP Act, 2023**: Data Protection and Digital Privacy compliance
  - Consumer consent for data collection
  - Right to erasure (off-chain data only; on-chain data is pseudonymized)
  - Data localization: all data stored on Indian servers
  - Purpose limitation: data used only for metering and billing
- **Aadhaar Act**: No Aadhaar data stored on blockchain; identity via separate KYC

### CERT-In Guidelines
- Security incident reporting procedures in place
- Vulnerability management for blockchain nodes
- Access control via X.509 certificates and JWT

## 3. Security Compliance

### Data Security
- **Encryption at rest**: SQLite with encryption extension
- **Encryption in transit**: TLS 1.3 for all API communication, TLS for Fabric
- **Hashing**: SHA-256 for data integrity verification
- **Key Management**: MSP-based PKI for blockchain identities

### Access Control
- Role-based access control (RBAC): Admin, Consumer, Regulator
- JWT authentication with 24-hour token expiry
- Fabric MSP for blockchain-level identity and authorization
- Endorsement policies require both Discom and Regulator approval

### Audit Trail
- Every transaction recorded on immutable blockchain ledger
- Off-chain audit log for API-level actions
- Tamper alerts with full forensic data

## 4. Privacy Architecture

### On-Chain Data (Immutable)
- Meter readings (pseudonymized meter IDs)
- Bill hashes (not full bill details)
- Tariff configurations
- Dispute status

### Off-Chain Data (Mutable, DPDP-compliant)
- Consumer personal information (name, email, phone)
- Detailed bill breakdowns
- User authentication credentials
- Analytics and dashboards

### Data Minimization
- Only essential data goes on-chain
- Personal data kept off-chain and deletable
- Blockchain contains hashes and references, not raw PII

## 5. Interoperability

### Standards Supported
- REST API (OpenAPI 3.0 compatible)
- JSON data format
- DLMS/COSEM meter protocol support (via API adapter)
- Hyperledger Fabric Gateway SDK

### Integration Capability
- Any MDMS via REST API
- Payment systems via webhook/callback
- Government platforms (UMANG, DigiLocker) via API gateway
- Cross-state federation via Fabric channels
