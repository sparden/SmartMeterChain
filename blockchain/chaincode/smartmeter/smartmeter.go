package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartMeterContract provides functions for managing smart meter data on the blockchain
type SmartMeterContract struct {
	contractapi.Contract
}

// InitLedger seeds the blockchain with sample meters and default tariffs
func (s *SmartMeterContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	// Seed default tariffs
	tariffs := []Tariff{
		{
			ID:       "TARIFF-DOMESTIC",
			Category: "domestic",
			Slabs: []TariffSlab{
				{MinUnits: 0, MaxUnits: 100, RatePerUnit: 3.50},
				{MinUnits: 101, MaxUnits: 200, RatePerUnit: 5.00},
				{MinUnits: 201, MaxUnits: 300, RatePerUnit: 6.50},
				{MinUnits: 301, MaxUnits: -1, RatePerUnit: 8.00},
			},
			FixedCharge:   50.00,
			EffectiveFrom: "2024-04-01",
			DocType:       "tariff",
		},
		{
			ID:       "TARIFF-COMMERCIAL",
			Category: "commercial",
			Slabs: []TariffSlab{
				{MinUnits: 0, MaxUnits: 100, RatePerUnit: 6.00},
				{MinUnits: 101, MaxUnits: 300, RatePerUnit: 8.00},
				{MinUnits: 301, MaxUnits: -1, RatePerUnit: 10.00},
			},
			FixedCharge:   150.00,
			EffectiveFrom: "2024-04-01",
			DocType:       "tariff",
		},
		{
			ID:       "TARIFF-INDUSTRIAL",
			Category: "industrial",
			Slabs: []TariffSlab{
				{MinUnits: 0, MaxUnits: 500, RatePerUnit: 7.50},
				{MinUnits: 501, MaxUnits: -1, RatePerUnit: 9.50},
			},
			FixedCharge:   300.00,
			EffectiveFrom: "2024-04-01",
			DocType:       "tariff",
		},
	}

	for _, tariff := range tariffs {
		tariffJSON, err := json.Marshal(tariff)
		if err != nil {
			return fmt.Errorf("failed to marshal tariff: %v", err)
		}
		err = ctx.GetStub().PutState(tariff.ID, tariffJSON)
		if err != nil {
			return fmt.Errorf("failed to put tariff: %v", err)
		}
	}

	// Seed sample meters
	meters := []Meter{
		{
			ID: "METER-MTR001", MeterID: "MTR-001", ConsumerID: "CONS-001",
			Location: "Block A, Sector 15, Noida", Type: "domestic",
			Status: "active", InstalledDate: "2024-01-15", DocType: "meter",
		},
		{
			ID: "METER-MTR002", MeterID: "MTR-002", ConsumerID: "CONS-002",
			Location: "Shop 42, Connaught Place, Delhi", Type: "commercial",
			Status: "active", InstalledDate: "2024-02-10", DocType: "meter",
		},
		{
			ID: "METER-MTR003", MeterID: "MTR-003", ConsumerID: "CONS-003",
			Location: "Plot 7, Industrial Area, Gurgaon", Type: "industrial",
			Status: "active", InstalledDate: "2024-03-01", DocType: "meter",
		},
	}

	for _, meter := range meters {
		meterJSON, err := json.Marshal(meter)
		if err != nil {
			return fmt.Errorf("failed to marshal meter: %v", err)
		}
		err = ctx.GetStub().PutState(meter.ID, meterJSON)
		if err != nil {
			return fmt.Errorf("failed to put meter: %v", err)
		}
	}

	return nil
}

// RegisterMeter registers a new smart meter on the blockchain
func (s *SmartMeterContract) RegisterMeter(ctx contractapi.TransactionContextInterface, meterJSON string) error {
	var meter Meter
	err := json.Unmarshal([]byte(meterJSON), &meter)
	if err != nil {
		return fmt.Errorf("failed to unmarshal meter: %v", err)
	}

	meter.DocType = "meter"
	meter.ID = "METER-" + meter.MeterID

	existing, err := ctx.GetStub().GetState(meter.ID)
	if err != nil {
		return fmt.Errorf("failed to read from world state: %v", err)
	}
	if existing != nil {
		return fmt.Errorf("meter %s already exists", meter.MeterID)
	}

	meterBytes, err := json.Marshal(meter)
	if err != nil {
		return fmt.Errorf("failed to marshal meter: %v", err)
	}

	return ctx.GetStub().PutState(meter.ID, meterBytes)
}

// SubmitReading submits a new meter reading and anchors its hash on the blockchain
func (s *SmartMeterContract) SubmitReading(ctx contractapi.TransactionContextInterface, readingJSON string) error {
	var reading MeterReading
	err := json.Unmarshal([]byte(readingJSON), &reading)
	if err != nil {
		return fmt.Errorf("failed to unmarshal reading: %v", err)
	}

	// Verify meter exists
	meterKey := "METER-" + reading.MeterID
	meterBytes, err := ctx.GetStub().GetState(meterKey)
	if err != nil {
		return fmt.Errorf("failed to read meter: %v", err)
	}
	if meterBytes == nil {
		return fmt.Errorf("meter %s does not exist", reading.MeterID)
	}

	// Compute hash
	hashInput := fmt.Sprintf("%s|%.4f|%s", reading.MeterID, reading.Reading, reading.Timestamp)
	hash := sha256.Sum256([]byte(hashInput))
	reading.Hash = fmt.Sprintf("%x", hash)
	reading.DocType = "reading"
	reading.Status = "valid"
	reading.ID = fmt.Sprintf("READING-%s-%s", reading.MeterID, reading.Timestamp)

	// Get previous reading hash for chain integrity
	var meter Meter
	err = json.Unmarshal(meterBytes, &meter)
	if err != nil {
		return fmt.Errorf("failed to unmarshal meter: %v", err)
	}
	reading.PreviousHash = meter.LastReadingHash

	// Update meter's last reading
	meter.LastReadingTime = reading.Timestamp
	meter.LastReadingHash = reading.Hash
	updatedMeter, err := json.Marshal(meter)
	if err != nil {
		return fmt.Errorf("failed to marshal updated meter: %v", err)
	}
	err = ctx.GetStub().PutState(meterKey, updatedMeter)
	if err != nil {
		return fmt.Errorf("failed to update meter: %v", err)
	}

	// Store reading
	readingBytes, err := json.Marshal(reading)
	if err != nil {
		return fmt.Errorf("failed to marshal reading: %v", err)
	}

	return ctx.GetStub().PutState(reading.ID, readingBytes)
}

// GetReadingsByMeter retrieves all readings for a specific meter
func (s *SmartMeterContract) GetReadingsByMeter(ctx contractapi.TransactionContextInterface, meterID string) ([]*MeterReading, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"reading","meterId":"%s"},"sort":[{"timestamp":"desc"}]}`, meterID)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()

	var readings []*MeterReading
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate: %v", err)
		}
		var reading MeterReading
		err = json.Unmarshal(queryResult.Value, &reading)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal reading: %v", err)
		}
		readings = append(readings, &reading)
	}

	return readings, nil
}

// GetMeter retrieves a single meter by its ID
func (s *SmartMeterContract) GetMeter(ctx contractapi.TransactionContextInterface, meterID string) (*Meter, error) {
	meterBytes, err := ctx.GetStub().GetState("METER-" + meterID)
	if err != nil {
		return nil, fmt.Errorf("failed to read meter: %v", err)
	}
	if meterBytes == nil {
		return nil, fmt.Errorf("meter %s does not exist", meterID)
	}

	var meter Meter
	err = json.Unmarshal(meterBytes, &meter)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal meter: %v", err)
	}

	return &meter, nil
}

// GetAllMeters retrieves all registered meters
func (s *SmartMeterContract) GetAllMeters(ctx contractapi.TransactionContextInterface) ([]*Meter, error) {
	queryString := `{"selector":{"docType":"meter"}}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer resultsIterator.Close()

	var meters []*Meter
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate: %v", err)
		}
		var meter Meter
		err = json.Unmarshal(queryResult.Value, &meter)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal meter: %v", err)
		}
		meters = append(meters, &meter)
	}

	return meters, nil
}

// DetectTamper checks if a reading hash matches the expected hash on-chain
func (s *SmartMeterContract) DetectTamper(ctx contractapi.TransactionContextInterface, meterID string, readingHash string) (*TamperAlert, error) {
	meterBytes, err := ctx.GetStub().GetState("METER-" + meterID)
	if err != nil {
		return nil, fmt.Errorf("failed to read meter: %v", err)
	}
	if meterBytes == nil {
		return nil, fmt.Errorf("meter %s does not exist", meterID)
	}

	var meter Meter
	err = json.Unmarshal(meterBytes, &meter)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal meter: %v", err)
	}

	if meter.LastReadingHash == "" {
		return nil, nil // No readings yet, nothing to compare
	}

	if meter.LastReadingHash != readingHash {
		alert := &TamperAlert{
			ID:           fmt.Sprintf("ALERT-%s-%d", meterID, time.Now().UnixMilli()),
			MeterID:      meterID,
			ExpectedHash: meter.LastReadingHash,
			ActualHash:   readingHash,
			DetectedAt:   time.Now().UTC().Format(time.RFC3339),
			Severity:     "critical",
			DocType:      "tamperAlert",
		}

		// Store alert on chain
		alertBytes, err := json.Marshal(alert)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal alert: %v", err)
		}
		err = ctx.GetStub().PutState(alert.ID, alertBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to store alert: %v", err)
		}

		// Mark meter as tampered
		meter.Status = "tampered"
		updatedMeter, err := json.Marshal(meter)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal meter: %v", err)
		}
		err = ctx.GetStub().PutState("METER-"+meterID, updatedMeter)
		if err != nil {
			return nil, fmt.Errorf("failed to update meter: %v", err)
		}

		return alert, nil
	}

	return nil, nil // Hashes match, no tamper detected
}
