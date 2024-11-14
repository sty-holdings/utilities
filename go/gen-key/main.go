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
	key         bool
	utilityName = "Gen Key Utility"
	testingOn   bool
)

func init() {

	appDescription := cases.Title(language.English).String(utilityName) + " will generate a key.\n"
	// Set your program's name and description.  These appear in help output.
	flaggy.SetName("\n" + utilityName) // "\n" is added to the start of the name to make the output easier to read.
	flaggy.SetDescription(appDescription)

	// You can disable various things by changing bool on the default parser
	// (or your own parser if you have created one).
	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	// You can set a help prepend or append on the default parser.
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/sty-holdings/utilities"

	// Add a flag to the main program (this will be available in all subcommands as well).
	flaggy.Bool(&key, "k", "key", "Generate a key.")
	flaggy.Bool(&testingOn, "t", "testingOn", "This puts the server into testing mode.")

	// Set the version and parse all inputs into variables.
	flaggy.Parse()
}

//goland:noinspection GoBoolExpressions
func main() {

	var (
		errorInfo    errs.ErrorInfo
		symmetricKey string
	)

	if symmetricKey = jwts.GenerateSymmetricKey(); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		os.Exit(1)
	}
	fmt.Printf("\nKey (32 byte Base64): %s\n", symmetricKey)
}
