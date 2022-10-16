package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/hex"
	"encoding/pem"
	"math/big"
	"net"
	"net/http"
	"os"
	"time"
)

type Handler struct {
	Func HandlerFunc
	Path string
}

type HandlerFunc func(w http.ResponseWriter, r *http.Request)

func cert(hosts []string) ([]byte, []byte, error) {
	// Create new private key object
	private, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	// Create random serial number
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	// Hash serial number for org name
	s := sha256.Sum256(serialNumber.Bytes())
	org := hex.EncodeToString(s[:])[:25]

	// Create cert object with our info. Valid for a year.
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{org},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour),

		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Set some configs for self signing
	template.IsCA = true
	template.KeyUsage = x509.KeyUsageCertSign

	certBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, private.Public(), private)
	if err != nil {
		return nil, nil, err
	}

	cert := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	privateBytes, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		return nil, nil, err
	}

	//pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privateBytes})
	key := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateBytes,
	})
	// Add hosts
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}
	return cert, key, nil
}

// NewHandler creates a handler object which forces a relationship between a handler and a path
func NewHandler(p string, f HandlerFunc) Handler {
	return Handler{
		Path: p,
		Func: f,
	}
}

// Server starts a new TLS 1.3 server on the specified host addresses/names and handles each handlerß
func Server(hosts []string, handlers []Handler) error {
	// Generate new cert on the fly
	cert, key, err := cert(hosts)
	if err != nil {
		return err
	}

	// Write cert
	certFile, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}

	err = os.WriteFile(certFile.Name(), cert, 766)
	if err != nil {
		return err
	}

	// Write key
	keyFile, err := os.CreateTemp("", "")
	if err != nil {
		return err
	}
	err = os.WriteFile(keyFile.Name(), key, 766)
	if err != nil {
		return err
	}

	// Create new multiplexer to handle connections
	mux := http.NewServeMux()

	// Register any handlers
	for _, h := range handlers {
		mux.HandleFunc(h.Path, h.Func)
	}

	// With TLS 1.3 all cipher suites are considered secure and not configurable
	cfg := &tls.Config{
		MinVersion: tls.VersionTLS13,
	}

	// Setup new server with our multiplexer
	srv := &http.Server{
		Addr:         ":443",
		Handler:      mux,
		TLSConfig:    cfg,
		TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	}

	// ListenAndServe will never return nil error. ErrServerClosed is safe run. ß
	err = srv.ListenAndServeTLS(certFile.Name(), keyFile.Name())
	if err == http.ErrServerClosed {
		err = nil
	}

	// Clean up
	err = os.Remove(certFile.Name())
	if err != nil {
		return err
	}

	return err
}
