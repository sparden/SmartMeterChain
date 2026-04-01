package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SetTariff creates or updates a tariff for a consumer category
func (s *SmartMeterContract) SetTariff(ctx contractapi.TransactionContextInterface, tariffJSON string) error {
	var tariff Tariff
	err := json.Unmarshal([]byte(tariffJSON), &tariff)
	if err != nil {
		return fmt.Errorf("failed to unmarshal tariff: %v", err)
	}

	tariff.DocType = "tariff"
	tariff.ID = "TARIFF-" + strings.ToUpper(tariff.Category)

	// Validate slabs
	if len(tariff.Slabs) == 0 {
		return fmt.Errorf("tariff must have at least one slab")
	}

	for i, slab := range tariff.Slabs {
		if slab.RatePerUnit <= 0 {
			return fmt.Errorf("slab %d rate must be positive", i)
		}
		if slab.MaxUnits != -1 && slab.MaxUnits <= slab.MinUnits {
			return fmt.Errorf("slab %d: maxUnits must be greater than minUnits", i)
		}
	}

	tariffBytes, err := json.Marshal(tariff)
	if err != nil {
		return fmt.Errorf("failed to marshal tariff: %v", err)
	}

	return ctx.GetStub().PutState(tariff.ID, tariffBytes)
}

// GetTariff retrieves the tariff for a specific category
func (s *SmartMeterContract) GetTariff(ctx contractapi.TransactionContextInterface, category string) (*Tariff, error) {
	tariffKey := "TARIFF-" + strings.ToUpper(category)
	tariffBytes, err := ctx.GetStub().GetState(tariffKey)
	if err != nil {
		return nil, fmt.Errorf("failed to read tariff: %v", err)
	}
	if tariffBytes == nil {
		return nil, fmt.Errorf("tariff for category %s not found", category)
	}

	var tariff Tariff
	err = json.Unmarshal(tariffBytes, &tariff)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal tariff: %v", err)
	}

	return &tariff, nil
}

// GetAllTariffs retrieves all defined tariffs
func (s *SmartMeterContract) GetAllTariffs(ctx contractapi.TransactionContextInterface) ([]*Tariff, error) {
	queryString := `{"selector":{"docType":"tariff"}}`

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to query tariffs: %v", err)
	}
	defer resultsIterator.Close()

	var tariffs []*Tariff
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate: %v", err)
		}
		var tariff Tariff
		err = json.Unmarshal(queryResult.Value, &tariff)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal tariff: %v", err)
		}
		tariffs = append(tariffs, &tariff)
	}

	return tariffs, nil
}
