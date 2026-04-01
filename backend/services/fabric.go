package services

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/smartmeterchain/backend/config"
)

// FabricService wraps interactions with Hyperledger Fabric chaincode.
// When FABRIC_ENABLED=false, it operates in mock mode for development.
type FabricService struct {
	cfg     *config.Config
	enabled bool
}

func NewFabricService(cfg *config.Config) *FabricService {
	fs := &FabricService{
		cfg:     cfg,
		enabled: cfg.FabricEnabled,
	}
	if !fs.enabled {
		log.Println("[Fabric] Running in MOCK mode (set FABRIC_ENABLED=true for real blockchain)")
	} else {
		log.Println("[Fabric] Connecting to Hyperledger Fabric network...")
	}
	return fs
}

// SubmitTransaction submits a transaction to the chaincode
func (fs *FabricService) SubmitTransaction(fn string, args ...string) ([]byte, string, error) {
	if !fs.enabled {
		return fs.mockSubmit(fn, args...)
	}
	return fs.fabricSubmit(fn, args...)
}

// EvaluateTransaction queries the chaincode (read-only)
func (fs *FabricService) EvaluateTransaction(fn string, args ...string) ([]byte, error) {
	if !fs.enabled {
		return fs.mockEvaluate(fn, args...)
	}
	return fs.fabricEvaluate(fn, args...)
}

func (fs *FabricService) fabricSubmit(fn string, args ...string) ([]byte, string, error) {
	// Real Fabric Gateway SDK implementation
	// Uses cfg.FabricGateway, cfg.FabricChannel, cfg.FabricCC
	return nil, "", fmt.Errorf("fabric submit not implemented — enable mock mode for development")
}

func (fs *FabricService) fabricEvaluate(fn string, args ...string) ([]byte, error) {
	return nil, fmt.Errorf("fabric evaluate not implemented — enable mock mode for development")
}

// mockSubmit simulates chaincode transactions for development
func (fs *FabricService) mockSubmit(fn string, args ...string) ([]byte, string, error) {
	txID := fmt.Sprintf("mock-tx-%s-%d", fn, len(args))
	log.Printf("[Fabric-Mock] Submit: %s(%v) -> tx:%s", fn, args, txID)

	response := map[string]interface{}{
		"function": fn,
		"args":     args,
		"tx_id":    txID,
		"status":   "SUCCESS",
	}
	data, _ := json.Marshal(response)
	return data, txID, nil
}

// mockEvaluate simulates chaincode queries for development
func (fs *FabricService) mockEvaluate(fn string, args ...string) ([]byte, error) {
	log.Printf("[Fabric-Mock] Evaluate: %s(%v)", fn, args)

	response := map[string]interface{}{
		"function": fn,
		"args":     args,
		"status":   "SUCCESS",
		"result":   "mock-data",
	}
	data, _ := json.Marshal(response)
	return data, nil
}
