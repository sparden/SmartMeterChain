package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// FileDispute creates a new dispute for a bill
func (s *SmartMeterContract) FileDispute(ctx contractapi.TransactionContextInterface, disputeJSON string) error {
	var dispute Dispute
	err := json.Unmarshal([]byte(disputeJSON), &dispute)
	if err != nil {
		return fmt.Errorf("failed to unmarshal dispute: %v", err)
	}

	// Verify the bill exists
	billBytes, err := ctx.GetStub().GetState("BILL-" + dispute.BillID)
	if err != nil {
		return fmt.Errorf("failed to read bill: %v", err)
	}
	if billBytes == nil {
		return fmt.Errorf("bill %s does not exist", dispute.BillID)
	}

	// Update bill status to disputed
	var bill Bill
	err = json.Unmarshal(billBytes, &bill)
	if err != nil {
		return fmt.Errorf("failed to unmarshal bill: %v", err)
	}
	bill.Status = "disputed"
	updatedBill, err := json.Marshal(bill)
	if err != nil {
		return fmt.Errorf("failed to marshal bill: %v", err)
	}
	err = ctx.GetStub().PutState("BILL-"+dispute.BillID, updatedBill)
	if err != nil {
		return fmt.Errorf("failed to update bill: %v", err)
	}

	dispute.DocType = "dispute"
	dispute.Status = "open"
	dispute.FiledAt = time.Now().UTC().Format(time.RFC3339)
	dispute.ID = fmt.Sprintf("DISPUTE-%s-%d", dispute.BillID, time.Now().UnixMilli())

	disputeBytes, err := json.Marshal(dispute)
	if err != nil {
		return fmt.Errorf("failed to marshal dispute: %v", err)
	}

	return ctx.GetStub().PutState(dispute.ID, disputeBytes)
}

// ResolveDispute resolves or rejects a dispute
func (s *SmartMeterContract) ResolveDispute(ctx contractapi.TransactionContextInterface, disputeID string, resolution string, status string) error {
	disputeBytes, err := ctx.GetStub().GetState(disputeID)
	if err != nil {
		return fmt.Errorf("failed to read dispute: %v", err)
	}
	if disputeBytes == nil {
		return fmt.Errorf("dispute %s does not exist", disputeID)
	}

	var dispute Dispute
	err = json.Unmarshal(disputeBytes, &dispute)
	if err != nil {
		return fmt.Errorf("failed to unmarshal dispute: %v", err)
	}

	if dispute.Status != "open" && dispute.Status != "investigating" {
		return fmt.Errorf("dispute %s is already %s", disputeID, dispute.Status)
	}

	if status != "resolved" && status != "rejected" && status != "investigating" {
		return fmt.Errorf("invalid status: %s (must be resolved, rejected, or investigating)", status)
	}

	dispute.Status = status
	dispute.Resolution = resolution
	if status == "resolved" || status == "rejected" {
		dispute.ResolvedAt = time.Now().UTC().Format(time.RFC3339)
	}

	updatedDispute, err := json.Marshal(dispute)
	if err != nil {
		return fmt.Errorf("failed to marshal dispute: %v", err)
	}

	return ctx.GetStub().PutState(disputeID, updatedDispute)
}

// GetDispute retrieves a single dispute
func (s *SmartMeterContract) GetDispute(ctx contractapi.TransactionContextInterface, disputeID string) (*Dispute, error) {
	disputeBytes, err := ctx.GetStub().GetState(disputeID)
	if err != nil {
		return nil, fmt.Errorf("failed to read dispute: %v", err)
	}
	if disputeBytes == nil {
		return nil, fmt.Errorf("dispute %s does not exist", disputeID)
	}

	var dispute Dispute
	err = json.Unmarshal(disputeBytes, &dispute)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal dispute: %v", err)
	}

	return &dispute, nil
}

// GetDisputesByConsumer retrieves all disputes filed by a consumer
func (s *SmartMeterContract) GetDisputesByConsumer(ctx contractapi.TransactionContextInterface, consumerID string) ([]*Dispute, error) {
	queryString := fmt.Sprintf(`{"selector":{"docType":"dispute","consumerId":"%s"},"sort":[{"filedAt":"desc"}]}`, consumerID)

	resultsIterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to query disputes: %v", err)
	}
	defer resultsIterator.Close()

	var disputes []*Dispute
	for resultsIterator.HasNext() {
		queryResult, err := resultsIterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to iterate: %v", err)
		}
		var dispute Dispute
		err = json.Unmarshal(queryResult.Value, &dispute)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal dispute: %v", err)
		}
		disputes = append(disputes, &dispute)
	}

	return disputes, nil
}
