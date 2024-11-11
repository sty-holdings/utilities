package main

import (
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func main() {
	log.Printf("Usage: verify_certificate SERVER_NAME CERT.pem CHAIN.pem")

	serverName := os.Args[1]

	certPEM, err := os.ReadFile(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	rootPEM, err := os.ReadFile(os.Args[3])
	if err != nil {
		log.Fatal(err)
	}

	roots := x509.NewCertPool()
	ok := roots.AppendCertsFromPEM([]byte(rootPEM))
	if !ok {
		panic("failed to parse root certificate")
	}

	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		panic("failed to parse certificate PEM")
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("failed to parse certificate: " + err.Error())
	}

	opts := x509.VerifyOptions{
		Roots:         roots,
		DNSName:       serverName,
		Intermediates: x509.NewCertPool(),
	}

	if _, err := cert.Verify(opts); err != nil {
		panic("failed to verify certificate: " + err.Error())
	}

	log.Printf("verification succeeds")
}
