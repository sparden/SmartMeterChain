package main

import (
	"log"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

func main() {
	smartMeterChaincode, err := contractapi.NewChaincode(&SmartMeterContract{})
	if err != nil {
		log.Fatalf("Error creating SmartMeter chaincode: %v", err)
	}

	if err := smartMeterChaincode.Start(); err != nil {
		log.Fatalf("Error starting SmartMeter chaincode: %v", err)
	}
}
