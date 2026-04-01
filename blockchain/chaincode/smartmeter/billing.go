package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// GenerateBill creates a bill for a consumer based on readings within a billing period
func (s *SmartMeterContract) GenerateBill(ctx contractapi.TransactionContextInterface, consumerID string, meterID string, periodStart string, periodEnd string) (*Bill, error) {
	// Get tariff for this meter's category
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

	// Get tariff
	tariffKey := "TARIFF-" + capitalize(meter.Type)
	tariffBytes, err := ctx.GetStub().GetState(tariffKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read tariff: %v", err)
	}
	if tariffBytes == nil {
		return nil, fmt.Errorf("tariff for category %s not found", meter.Type)
	}

	var tariff Tariff
	err = json.Unmarshal(tariffBytes, &tariff)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tariff: %v", err)
	}

	// Query readings for the billing period
	queryString := fmt.Sprintf(
		`{"selector":{"docType":"reading","meterId":"%s","timestamp":{"$gte":"%s","$lte":"%s"}},"sort":[{"timestamp":"asc"}]}`,
		meterID, periodStart, periodEnd,
	)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to query readings: %v", err)
	}
	defer resultsIterator.Close()

	// Calculate total units consumed
	var firstReading, lastReading float64
	hasReadings := false
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
		if !hasReadings {
			firstReading = reading.Reading
			hasReadings = true
		}
		lastReading = reading.Reading
	}

	if !hasReadings {
		return nil, fmt.Errorf("no readings found for meter %s in the billing period", meterID)
	}

	unitsConsumed := lastReading - firstReading
	if unitsConsumed < 0 {
		unitsConsumed = 0
	}

	// Apply slab-based tariff calculation
	slabBreakdown, energyCharge := calculateSlabCharges(unitsConsumed, tariff.Slabs)
	totalAmount := energyCharge + tariff.FixedCharge

	// Generate bill
	now := time.Now().UTC().Format(time.RFC3339)
	billNumber := fmt.Sprintf("BILL-%s-%s", consumerID, periodEnd[:7])

	hashInput := fmt.Sprintf("%s|%s|%.2f|%s|%s", billNumber, consumerID, totalAmount, periodStart, periodEnd)
	hashBytes := sha256.Sum256([]byte(hashInput))

	bill := &Bill{
		ID:                 "BILL-" + billNumber,
		BillNumber:         billNumber,
		ConsumerID:         consumerID,
		MeterID:            meterID,
		BillingPeriodStart: periodStart,
		BillingPeriodEnd:   periodEnd,
		UnitsConsumed:      unitsConsumed,
		SlabBreakdown:      slabBreakdown,
		FixedCharge:        tariff.FixedCharge,
		TotalAmount:        totalAmount,
		Hash:               fmt.Sprintf("%x", hashBytes),
		Status:             "generated",
		GeneratedAt:        now,
		DocType:            "bill",
	}

	billBytes, err := json.Marshal(bill)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal bill: %v", err)
	}

	err = ctx.GetStub().PutState(bill.ID, billBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to store bill: %v", err)
	}

	return bill, nil
}

// calculateSlabCharges applies Indian slab-based tariff to units consumed
func calculateSlabCharges(units float64, slabs []TariffSlab) ([]SlabCharge, float64) {
	var breakdown []SlabCharge
	var totalCharge float64
	remaining := units

	for _, slab := range slabs {
		if remaining <= 0 {
			break
		}

		slabRange := slab.MaxUnits - slab.MinUnits
		if slab.MaxUnits == -1 {
			slabRange = remaining // Unlimited slab
		}

		unitsInSlab := remaining
		if unitsInSlab > slabRange {
			unitsInSlab = slabRange
		}

		charge := unitsInSlab * slab.RatePerUnit
		totalCharge += charge

		breakdown = append(breakdown, SlabCharge{
			SlabName:    fmt.Sprintf("%.0f-%.0f units", slab.MinUnits, slab.MaxUnits),
			Units:       unitsInSlab,
			RatePerUnit: slab.RatePerUnit,
			Amount:      charge,
		})

		remaining -= unitsInSlab
	}

	return breakdown, totalCharge
}

// GetBill retrieves a bill by its ID
func (s *SmartMeterContract) GetBill(ctx contractapi.TransactionContextInterface, billID string) (*Bill, error) {
	billBytes, err := ctx.GetStub().GetState("BILL-" + billID)
	if err != nil {
		return nil, fmt.Errorf("failed to read bill: %v", err)
	}
	if billBytes == nil {
		return nil, fmt.Errorf("bill %s does not exist", billID)
	}

	var bill Bill
	err = json.Unmarshal(billBytes, &bill)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal bill: %v", err)
	}

	return &bill, nil
}

// GetBillsByConsumer retrieves all bills for a consumer
func (s *SmartMeterContract) GetBillsByConsumer(ctx contractapi.TransactionContextInterface, consumerID string) ([]*Bill, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"bill","consumerId":"%s"},"sort":[{"generatedAt":"desc"}]}`, consumerID)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to query bills: %v", err)
	}
	defer resultsIterator.Close()

	var bills []*Bill
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate: %v", err)
		}
		var bill Bill
		err = json.Unmarshal(queryResult.Value, &bill)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal bill: %v", err)
		}
		bills = append(bills, &bill)
	}

	return bills, nil
}

// VerifyBill verifies the integrity of a bill by recomputing its hash
func (s *SmartMeterContract) VerifyBill(ctx contractapi.TransactionContextInterface, billID string) (bool, error) {
	bill, err := s.GetBill(ctx, billID)
	if err != nil {
		return false, err
	}

	hashInput := fmt.Sprintf("%s|%s|%.2f|%s|%s", bill.BillNumber, bill.ConsumerID, bill.TotalAmount, bill.BillingPeriodStart, bill.BillingPeriodEnd)
	hashBytes := sha256.Sum256([]byte(hashInput))
	computedHash := fmt.Sprintf("%x", hashBytes)

	return computedHash == bill.Hash, nil
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	upper := s[0]
	if upper >= 'a' && upper <= 'z' {
		upper -= 32
	}
	return string(upper) + s[1:]
}
