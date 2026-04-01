package config

import (
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewGrpcConnection creates a gRPC connection to the Fabric gateway peer.
func NewGrpcConnection(cfg *Config) (*grpc.ClientConn, error) {
	tlsCertPath := cfg.FabricCertPath
	if tlsCertPath == "" {
		return nil, fmt.Errorf("FABRIC_CERT_PATH is not set")
	}

	certificate, err := os.ReadFile(tlsCertPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read TLS cert: %w", err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(certificate) {
		return nil, fmt.Errorf("failed to add server certificate to pool")
	}

	transportCredentials := credentials.NewClientTLSFromCert(certPool, "peer0.discom.smartmeterchain.com")

	conn, err := grpc.NewClient(cfg.FabricGateway, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection: %w", err)
	}

	return conn, nil
}
