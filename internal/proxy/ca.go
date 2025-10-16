package proxy

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

// CA represents a Certificate Authority
type CA struct {
	Cert       *x509.Certificate
	PrivateKey *rsa.PrivateKey
	certPath   string
	keyPath    string
}

// GetCADir returns the mozzy CA directory
func GetCADir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".mozzy")
}

// GetCA loads or generates the CA certificate
func GetCA() (*CA, error) {
	caDir := GetCADir()
	certPath := filepath.Join(caDir, "ca-cert.pem")
	keyPath := filepath.Join(caDir, "ca-key.pem")

	// Check if CA already exists
	if fileExists(certPath) && fileExists(keyPath) {
		return loadCA(certPath, keyPath)
	}

	// Generate new CA
	return generateCA(certPath, keyPath)
}

// loadCA loads an existing CA from disk
func loadCA(certPath, keyPath string) (*CA, error) {
	// Load certificate
	certPEM, err := os.ReadFile(certPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	block, _ := pem.Decode(certPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode CA cert PEM")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA cert: %w", err)
	}

	// Load private key
	keyPEM, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA key: %w", err)
	}

	keyBlock, _ := pem.Decode(keyPEM)
	if keyBlock == nil {
		return nil, fmt.Errorf("failed to decode CA key PEM")
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse CA key: %w", err)
	}

	color.Green("✓ Loaded CA certificate from %s", certPath)

	return &CA{
		Cert:       cert,
		PrivateKey: privateKey,
		certPath:   certPath,
		keyPath:    keyPath,
	}, nil
}

// generateCA creates a new CA certificate
func generateCA(certPath, keyPath string) (*CA, error) {
	color.Yellow("📜 Generating new CA certificate...")

	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create CA certificate template
	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   "Mozzy Proxy CA",
			Organization: []string{"Mozzy"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0), // Valid for 10 years
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  true,
		MaxPathLen:            0,
		MaxPathLenZero:        true,
	}

	// Self-sign the certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, template, template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(certPath), 0755); err != nil {
		return nil, fmt.Errorf("failed to create CA directory: %w", err)
	}

	// Save certificate
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if err := os.WriteFile(certPath, certPEM, 0644); err != nil {
		return nil, fmt.Errorf("failed to write CA cert: %w", err)
	}

	// Save private key with restricted permissions
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})
	if err := os.WriteFile(keyPath, keyPEM, 0600); err != nil {
		return nil, fmt.Errorf("failed to write CA key: %w", err)
	}

	color.Green("✓ CA certificate generated and saved")
	color.Cyan("  Certificate: %s", certPath)
	color.Cyan("  Private Key: %s", keyPath)
	fmt.Println()

	return &CA{
		Cert:       cert,
		PrivateKey: privateKey,
		certPath:   certPath,
		keyPath:    keyPath,
	}, nil
}

// GenerateServerCert generates a certificate for a specific host
func (ca *CA) GenerateServerCert(host string) (*x509.Certificate, *rsa.PrivateKey, error) {
	// Generate private key for the server
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate server key: %w", err)
	}

	// Create certificate template
	serialNumber, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   host,
			Organization: []string{"Mozzy Proxy"},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().AddDate(1, 0, 0), // Valid for 1 year
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		DNSNames:    []string{host},
	}

	// Sign with CA
	certBytes, err := x509.CreateCertificate(rand.Reader, template, ca.Cert, &privateKey.PublicKey, ca.PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create server certificate: %w", err)
	}

	cert, err := x509.ParseCertificate(certBytes)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse server certificate: %w", err)
	}

	return cert, privateKey, nil
}

// ExportCert exports the CA certificate in PEM format
func (ca *CA) ExportCert() ([]byte, error) {
	return pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: ca.Cert.Raw,
	}), nil
}

// GetInfo returns information about the CA certificate
func (ca *CA) GetInfo() string {
	return fmt.Sprintf(`Mozzy Proxy Certificate Authority

Subject: %s
Issuer:  %s
Serial:  %s

Valid From: %s
Valid To:   %s

Certificate: %s
Private Key: %s`,
		ca.Cert.Subject.CommonName,
		ca.Cert.Issuer.CommonName,
		ca.Cert.SerialNumber.String(),
		ca.Cert.NotBefore.Format("2006-01-02 15:04:05"),
		ca.Cert.NotAfter.Format("2006-01-02 15:04:05"),
		ca.certPath,
		ca.keyPath,
	)
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
