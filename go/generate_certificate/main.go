// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Generate a self-signed X.509 certificate for a TLS server. Outputs to
// 'cert.pem' and 'key.pem' and will overwrite existing files.

package main

import (
	"fmt"
	"regexp"

	cmd "albert/Utilities/generate_certificate/src"
	"albert/constants"
	"github.com/integrii/flaggy"
)

//goland:noinspection ALL
const (
	VERSION          = "2023.1.0"
	APPLICATION_NAME = "Generate Certificate"
)

var (
	host     string
	validFor string
	selfCA   bool
	rsaBits  = 4096
	//
	keyFileName  string
	certFileName string
)

func init() {

	appDescription := "Generate certificate will create a csr and key file that are self-signed.\n" +
		"\nVersion: \n" +
		constants.FOUR_SPACES + "- " + VERSION + "\n" +
		"\nConstraints: \n" +
		constants.FOUR_SPACES + "- There is no log for this utility. All messages are output to the console.\n" +
		constants.FOUR_SPACES + "- Only Self-Signed RSA certificates are created.\n" +
		"\nNotes:\n" +
		constants.FOUR_SPACES + "Key files will be output with no extension for the private key, .pub for the public key, and .pem for the cert.\n" +
		constants.FOUR_SPACES + "The files permissions are set to 0744.\n" +
		"\nFor more info, see link below:\n"

	// Set your program's name and description.  These appear in help output.
	flaggy.SetName("\n" + APPLICATION_NAME) // "\n" is added to the start of the name to make the output easier to read.
	flaggy.SetDescription(appDescription)

	// You can disable various things by changing bool on the default parser
	// (or your own parser if you have created one).
	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	// You can set a help prepend or append on the default parser.
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/styh-dev/albert/utilitiles/generate_certificate"

	// Add a flag to the main program (this will be available in all subcommands as well).
	flaggy.String(&host, "n", "hostname", "REQUIRED: The DNS name of the system where the certificate will be installed.")
	flaggy.String(&validFor, "v", "valid_for", "REQUIRED: The length of time the certificate is valid, expressed in <digits><unit>."+
		"\n\t\t\t<digits> is any number from 1 to 10"+"\n\t\t\t<period> can be d | D for days, m | M for months, or y | Y for years.")
	flaggy.String(&keyFileName, "k", "key_name", "REQUIRED: The directory and filename of the out key file. DO NOT provide an extension to the name.")
	flaggy.String(&certFileName, "c", "cert_name", "REQUIRED: The directory and filename of the out certificate file. DO NOT provide an extension to the name.")
	flaggy.Bool(&selfCA, "s", "self_CA", "Will this certification be its own Certificate Authority. The default is false.")
	flaggy.Int(&rsaBits, "r", "rsa_bits", "Size of RSA key to generate. The value must be 1024 or higher when supplied. The default is 4096.")

	// Set the version and parse all inputs into variables.
	flaggy.SetVersion(VERSION)
	flaggy.Parse()
}

func main() {

	if host == constants.EMPTY || validFor == constants.EMPTY || keyFileName == constants.EMPTY || certFileName == constants.EMPTY {
		flaggy.ShowHelpAndExit(fmt.Sprintln(constants.COLOR_RED, "ERROR: Please review the usage and supplied all required arguments."))
	}

	if match, _ := regexp.MatchString("^([1-9]|10)[dDmMyY]", validFor); match == false {
		flaggy.ShowHelpAndExit(fmt.Sprintln(constants.COLOR_RED, "ERROR: Please review the format for the 'valid_for' argument."))
	}

	if rsaBits < 1024 {
		flaggy.ShowHelpAndExit(fmt.Sprintln(constants.COLOR_RED, "ERROR: 'rsa_bits' must be 1024 or higher."))
	}

	cmd.Run(coreJWT.GenerateCertificate{
		CertFileName:       certFileName,
		Certificate:        nil,
		Host:               host,
		PublicKey:          nil,
		PrivateKey:         nil,
		PrivateKeyFileName: keyFileName,
		RSABits:            rsaBits,
		SelfCA:             selfCA,
		ValidFor:           validFor,
	})
}
