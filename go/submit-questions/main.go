package submit_questions

import (
	"fmt"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/integrii/flaggy"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	ctv "github.com/sty-holdings/sharedServices/v2024/constantsTypesVars"
	errs "github.com/sty-holdings/sharedServices/v2024/errorServices"
	fbs "github.com/sty-holdings/sharedServices/v2024/firebaseServices"
	fss "github.com/sty-holdings/sharedServices/v2024/firestoreServices"
)

// Add types to the request_reply_types.go or the data_structure_types.go file

var (
	// Add Variables here for the file (Remember, they are global)
	// Start up values for a service
	enterQuestion bool // False means that the utility will read the questions from Firestore training questions collection.
	question      string
	fbCredentials string
	utilityName   = "Submit Question"
	testingOn     bool
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
	flaggy.Bool(&enterQuestion, "e", "enter", "You want to submit a question. False means training questions will be processed.")
	flaggy.Bool(&testingOn, "t", "testingOn", "This puts the server into testing mode.")
	flaggy.String(&question, "q", "question", "Question your submitting.")
	flaggy.String(&fbCredentials, "c", "fbcreds", "The filename of your firebase credentials.")

	// Set the version and parse all inputs into variables.
	flaggy.Parse()
}

//goland:noinspection GoBoolExpressions
func main() {

	var (
	//errorInfo errs.ErrorInfo
	)

	if enterQuestion && question == ctv.VAL_EMPTY {
		flaggy.ShowHelpAndExit("You must provide a question.")
	}

	if enterQuestion {
		// Send single question
	}

	processTrainingData(fbCredentials)
}

func processTrainingData(fbCredentials string) {

	var (
		errorInfo    errs.ErrorInfo
		tAppPtr      *firebase.App
		tFSClientPtr *firestore.Client
		tDocRefPtr   []*firestore.DocumentSnapshot
	)

	if tAppPtr, _, errorInfo = fbs.GetFirebaseAppAuthConnection(fbCredentials); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		return
	}
	if tFSClientPtr, errorInfo = fss.GetFirestoreClientConnection(tAppPtr); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		return
	}

	if tDocRefPtr, errorInfo = fss.GetAllDocuments(tFSClientPtr, ctv.DATASTORE_TRAINING_QUESTIONS); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		return
	}

	for _, snapshot := range tDocRefPtr {
		if !snapshot.Exists() {
			fmt.Printf("Document does not exist\n")
			continue
		}

		data := snapshot.Data()
		// Access data fields using their key, e.g.:
		fmt.Printf("Document ID: %s, Field1: %v\n", snapshot.Ref.ID, data["Field1"])

		// You can also extract the data into a struct:
		type MyStruct struct {
			Field1 string `firestore:"Field1"`
			Field2 int    `firestore:"Field2"`
		}
		var myData MyStruct
		if err := snapshot.DataTo(&myData); err != nil {
			fmt.Printf("Error converting data to struct: %v\n", err)
			continue
		}
		fmt.Printf("Struct data: %+v\n", myData)
	}
	return
}
