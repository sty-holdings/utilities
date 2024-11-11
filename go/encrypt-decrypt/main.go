package main

import (
	"fmt"
	"os"

	"github.com/integrii/flaggy"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	errs "github.com/sty-holdings/sharedServices/v2024/errorServices"
	jwts "github.com/sty-holdings/sharedServices/v2024/jwtServices"
)

// Add types to the request_reply_types.go or the data_structure_types.go file

var (
	// Add Variables here for the file (Remember, they are global)
	// Start up values for a service
	key         string
	message     string
	utilityName = "Encrypt/Decrypt Utility"
	testingOn   bool
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
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/styh-dev/albert"

	// Add a flag to the main program (this will be available in all subcommands as well).
	flaggy.String(&key, "k", "key", "The encryption / decryption key.")
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
	if (key == "" || message == "-t") && testingOn == false {
		flaggy.ShowHelpAndExit("")
	}

	if eMessage, errorInfo = jwts.Encrypt("Encrypt/Decrypt Utility", key, message); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		os.Exit(1)
	}

	fmt.Printf("Message:\t\t %s\n", message)
	fmt.Printf("Encrypted Message:\t %s\n", eMessage)

	if dMessage, errorInfo = jwts.Decrypt("Encrypt/Decrypt Utility", key, eMessage); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		os.Exit(1)
	}

	fmt.Printf("Decrypted Message:\t %s\n", dMessage)
}
