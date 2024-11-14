package main

import (
	"fmt"
	"os"

	"github.com/integrii/flaggy"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	ctv "github.com/sty-holdings/sharedServices/v2024/constantsTypesVars"
	errs "github.com/sty-holdings/sharedServices/v2024/errorServices"
	jwts "github.com/sty-holdings/sharedServices/v2024/jwtServices"
)

// Add types to the request_reply_types.go or the data_structure_types.go file

var (
	// Add Variables here for the file (Remember, they are global)
	// Start up values for a service
	key            string
	encryptMessage bool
	decryptMessage bool
	message        string
	utilityName    = "Encrypt/Decrypt Utility"
	testingOn      bool
)

func init() {

	appDescription := cases.Title(language.English).String(utilityName) + " will encrypt and decrypt a message using a key.\n"
	// Set your program's name and description.  These appear in help output.
	flaggy.SetName("\n" + utilityName) // "\n" is added to the start of the name to make the output easier to read.
	flaggy.SetDescription(appDescription)

	// You can disable various things by changing bool on the default parser
	// (or your own parser if you have created one).
	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	// You can set a help prepend or append on the default parser.
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/sty-holdings/utilities"

	// Add a flag to the main program (this will be available in all subcommands as well).
	flaggy.String(&key, "k", "key", "The encryption / decryption key.")
	flaggy.Bool(&encryptMessage, "e", "emsg", "Indicates that the message should be encrypted.")
	flaggy.Bool(&decryptMessage, "d", "dmsg", "Indicates that the message should be decrypted.")
	flaggy.String(
		&message, "m", "msg", "The message to encrypt/decrypt with.",
	)
	flaggy.Bool(&testingOn, "t", "testingOn", "This puts the server into testing mode.")

	// Set the version and parse all inputs into variables.
	flaggy.Parse()
}

//goland:noinspection GoBoolExpressions
func main() {

	var (
		errorInfo errs.ErrorInfo
		dMessage  string
		eMessage  string
	)

	// Has the config file location and name been provided, if not, return help.
	if (key == ctv.VAL_EMPTY || message == "-t" || message == ctv.VAL_EMPTY) && testingOn == false {
		flaggy.ShowHelpAndExit("")
	}

	if encryptMessage == false && decryptMessage == false {
		flaggy.ShowHelpAndExit("")
	}

	if decryptMessage && encryptMessage == false {
		eMessage = message
	}

	if encryptMessage {
		if eMessage, errorInfo = jwts.Encrypt(message, key, "Encrypt/Decrypt Utility"); errorInfo.Error != nil {
			errs.PrintErrorInfo(errorInfo)
			os.Exit(1)
		}
		fmt.Printf("Message:\t\t %s\n", message)
		fmt.Printf("Encrypted Message:\t %s\n", eMessage)
	}

	if decryptMessage {
		if dMessage, errorInfo = jwts.Decrypt(message, key, "Encrypt/Decrypt Utility"); errorInfo.Error != nil {
			errs.PrintErrorInfo(errorInfo)
			os.Exit(1)
		}

		fmt.Printf("Decrypted Message:\t %s\n", dMessage)
	}
}
