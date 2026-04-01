package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port           string
	JWTSecret      string
	DBPath         string
	FabricEnabled  bool
	FabricMSPID    string
	FabricCertPath string
	FabricKeyPath  string
	FabricGateway  string
	FabricChannel  string
	FabricCC       string
}

func Load() *Config {
	return &Config{
		Port:           getEnv("PORT", "3000"),
		JWTSecret:      getEnv("JWT_SECRET", "smartmeterchain-dev-secret-change-in-prod"),
		DBPath:         getEnv("DB_PATH", "./smartmeter.db"),
		FabricEnabled:  getEnvBool("FABRIC_ENABLED", false),
		FabricMSPID:    getEnv("FABRIC_MSPID", "DiscomMSP"),
		FabricCertPath: getEnv("FABRIC_CERT_PATH", ""),
		FabricKeyPath:  getEnv("FABRIC_KEY_PATH", ""),
		FabricGateway:  getEnv("FABRIC_GATEWAY", "localhost:7051"),
		FabricChannel:  getEnv("FABRIC_CHANNEL", "smartmeterchannel"),
		FabricCC:       getEnv("FABRIC_CHAINCODE", "smartmeter"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if val := os.Getenv(key); val != "" {
		b, err := strconv.ParseBool(val)
		if err == nil {
			return b
		}
	}
	return fallback
}
