// Package cmd
package src

import (
	"crypto"
	"crypto/rand"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"os"

	"crypto/x509"

	"albert/constants"
	"albert/core/coreError"

	"albert/core/coreJWT"

	// "math/big"
	// "net"
	// "os"
	"strconv"
	"strings"
	"time"
)

// Run
func Run(certificateInfo coreJWT.GenerateCertificate) {

	var (
		errorInfo coreError.ErrorInfo
	)

	certificateInfo.PrivateKey, certificateInfo.PublicKey, errorInfo = generateKeyFiles(certificateInfo.RSABits, certificateInfo.PrivateKeyFileName)

	if errorInfo.Error == nil {
		errorInfo = generateCertificateFile(certificateInfo)
	}

	if errorInfo.Error == nil {
		fmt.Println(constants.COLOR_GREEN, "The keys and certificate has been generated successfully.")
		fmt.Println(constants.COLOR_GREEN, fmt.Sprintf("\tCertificate Info: %v", certificateInfo.CertFileName+".pem"))
		fmt.Println(constants.COLOR_GREEN, fmt.Sprintf("\tPrivate Key File: %v", certificateInfo.PrivateKeyFileName))
		fmt.Println(constants.COLOR_GREEN, fmt.Sprintf("\tPublic Key File: %v", certificateInfo.PrivateKeyFileName+".pub"))
	}
}

// generateCertificateFile is output in the DER format.
func generateCertificateFile(certificateInfo coreJWT.GenerateCertificate) (errorInfo coreError.ErrorInfo) {

	var (
		tHost         []string
		tSerialNumber *big.Int
	)

	tPeriod, tDuration := parseValidFor(certificateInfo.ValidFor)

	if tSerialNumber, errorInfo = setSerialNumber(); errorInfo.Error != nil {
		fmt.Printf("%vERROR: %v", constants.COLOR_RED, errorInfo.Error.Error())
	}

	if errorInfo.Error == nil {
		tCertificateTemplate := x509.Certificate{
			SerialNumber: tSerialNumber,
			Subject: pkix.Name{
				Country:       []string{"US"},
				Organization:  []string{"STY Holdings Inc"},
				Locality:      []string{"California"},
				Province:      []string{constants.EMPTY},
				StreetAddress: []string{"San Francisco Bay Area"},
				// CommonName:    "STY Holdings Inc",
			},
			NotBefore:             time.Now(),
			NotAfter:              setCertificateExpiry(tPeriod, tDuration),
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
			DNSNames:              append(tHost, certificateInfo.Host),
			IsCA:                  certificateInfo.SelfCA,
		}

		if certificateInfo.SelfCA {
			tCertificateTemplate.IsCA = true
			tCertificateTemplate.KeyUsage |= x509.KeyUsageCertSign
		}

		if certificateInfo.Certificate, errorInfo.Error = x509.CreateCertificate(rand.Reader, &tCertificateTemplate, &tCertificateTemplate, certificateInfo.PublicKey, certificateInfo.PrivateKey); errorInfo.Error != nil {
			fmt.Printf("%vERROR: %v\n", constants.COLOR_RED, errorInfo.Error.Error())
			fmt.Println(constants.COLOR_RED, "Failed to create certificate: %v", errorInfo.Error)
			os.Exit(1)
		}

	}

	if errorInfo.Error == nil {
		errorInfo = writeOutPEMFile(certificateInfo.CertFileName+".pem", certificateInfo.Certificate, constants.CERTIFICATE)
	}

	return
}

// generateKeyFiles
func generateKeyFiles(rsaBits int, privateKeyFileName string) (privateKey crypto.PrivateKey, publicKey crypto.PublicKey, errorInfo coreError.ErrorInfo) {

	var (
		tKeyUsage          = x509.KeyUsageDigitalSignature
		tPublicKeyFileName = privateKeyFileName + ".pub"
	)

	if privateKey, publicKey, errorInfo = coreJWT.GenerateRSAKey(rsaBits); errorInfo.Error == nil {
		// Only RSA subject keys should have the Key Encipherment Key Usage bits set. In
		// the context of TLS this Key Usage is particular to RSA key exchange and
		// authentication.
		tKeyUsage |= x509.KeyUsageKeyEncipherment
		errorInfo = writePrivateKeys(privateKeyFileName, privateKey)
		if errorInfo.Error == nil {
			errorInfo = writePublicKeys(tPublicKeyFileName, publicKey)
		}
	}

	return
}

// parseValidFor
func parseValidFor(validFor string) (period string, duration int) {

	var (
		index      int
		tCharacter rune
		tDuration  string
	)

	tValidFor := strings.ToUpper(validFor)
	for index, tCharacter = range tValidFor {
		if index == len(tValidFor)-1 {
			period = string(tCharacter)
		} else {
			tDuration = tDuration + string(tCharacter)
		}
	}

	duration, _ = strconv.Atoi(tDuration)

	return
}

// setCertificateExpiry
func setCertificateExpiry(period string, duration int) (expiry time.Time) {

	switch strings.ToUpper(period) {
	case constants.DAY:
		expiry = time.Now().AddDate(0, 0, duration)
	case constants.MONTH:
		expiry = time.Now().AddDate(0, duration, 0)
	case constants.YEAR:
		expiry = time.Now().AddDate(duration, 0, 0)
	}

	return
}

// setSerialNumber
func setSerialNumber() (tSerialNumber *big.Int, errorInfo coreError.ErrorInfo) {

	// Local-sensitive hash (LSH) is used to create the serial number
	tSerialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	if tSerialNumber, errorInfo.Error = rand.Int(rand.Reader, tSerialNumberLimit); errorInfo.Error != nil {
		log.Fatalf("Failed to generate serial number: %v", errorInfo.Error)
	}

	return
}

// writeOutPEMFile
func writeOutPEMFile(fullyQualifiedName string, pemData []byte, keyType string) (errorInfo coreError.ErrorInfo) {

	var (
		keyOut *os.File
	)

	if keyOut, errorInfo.Error = os.Create(fullyQualifiedName); errorInfo.Error != nil {
		fmt.Println(constants.COLOR_RED, fmt.Sprintf("%vERROR: Unable to create %v. ERROR: %v", constants.COLOR_RED, fullyQualifiedName, errorInfo.Error.Error()))
	} else {
		defer func(keyOut *os.File) {
			if err := keyOut.Close(); err != nil {
				fmt.Println(constants.COLOR_RED, fmt.Sprintf("%vERROR: Unable to close %v. ERROR: %v", constants.COLOR_RED, fullyQualifiedName, err.Error()))
			}
		}(keyOut)
	}

	if errorInfo.Error == nil {
		if err := pem.Encode(keyOut, &pem.Block{Type: keyType, Bytes: pemData}); err != nil {
			fmt.Println(constants.COLOR_RED, fmt.Sprintf("%vERROR: Failed to write data to %v. ERROR: %v", constants.COLOR_RED, fullyQualifiedName, err.Error()))
		}
	}

	if errorInfo.Error == nil {
		if errorInfo.Error = os.Chmod(fullyQualifiedName, 0744); errorInfo.Error != nil {
			fmt.Println(constants.COLOR_RED, fmt.Sprintf("%vERROR: Unable to set permissions to %v. ERROR: %v", constants.COLOR_RED, fullyQualifiedName, errorInfo.Error.Error()))
		}
	}

	return
}

// writePrivateKeys
func writePrivateKeys(fullyQualifiedName string, privateKey crypto.PrivateKey) (errorInfo coreError.ErrorInfo) {

	var (
		tMarshalledPrivateKey []byte
	)

	if tMarshalledPrivateKey, errorInfo.Error = x509.MarshalPKCS8PrivateKey(privateKey); errorInfo.Error != nil {
		fmt.Println(constants.COLOR_RED, fmt.Sprintf("ERROR: Unable to marshal private key: %v", errorInfo.Error.Error()))
	}

	if errorInfo.Error == nil {
		errorInfo = writeOutPEMFile(fullyQualifiedName, tMarshalledPrivateKey, constants.CERT_PRIVATE_KEY)
	}

	return
}

// writePublicKeys
func writePublicKeys(fullyQualifiedName string, publicKey crypto.PublicKey) (errorInfo coreError.ErrorInfo) {

	var (
		tMarshalledPublicKey []byte
	)

	if tMarshalledPublicKey, errorInfo.Error = x509.MarshalPKIXPublicKey(publicKey); errorInfo.Error != nil {
		fmt.Println(constants.COLOR_RED, fmt.Sprintf("ERROR: Unable to marshal public key: %v", errorInfo.Error.Error()))
	}

	errorInfo = writeOutPEMFile(fullyQualifiedName, tMarshalledPublicKey, constants.CERT_PUBLIC_KEY)

	return
}
