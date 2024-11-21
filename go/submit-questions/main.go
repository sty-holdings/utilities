package main

import (
	"encoding/json"
	"fmt"
	"os"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/integrii/flaggy"
	"github.com/nats-io/nats.go"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	ctv "github.com/sty-holdings/sharedServices/v2024/constantsTypesVars"
	errs "github.com/sty-holdings/sharedServices/v2024/errorServices"
	fbs "github.com/sty-holdings/sharedServices/v2024/firebaseServices"
	jwts "github.com/sty-holdings/sharedServices/v2024/jwtServices"
	ns "github.com/sty-holdings/sharedServices/v2024/natsSerices"
)

var (
	// Add Variables here for the file (Remember, they are global)
	// Start up values for a service
	fbCredentials   string
	natsCABundle    string
	natsCertificate string
	natsCertKey     string
	natsCredentials string
	natsPort        string
	natsURL         string
	question        string
	analyzeProcess  string
	secretKey       string
	utilityName     = "Submit Question"
	testingOn       bool
	//
	username = "6oEtPOwn9hN2gNRYmGDQY3QXwrF2"
)

func init() {

	appDescription := cases.Title(language.English).String(utilityName) + " process a question and return an answer.\n"
	// Set your program's name and description.  These appear in help output.
	flaggy.SetName("\n" + utilityName) // "\n" is added to the start of the name to make the output easier to read.
	flaggy.SetDescription(appDescription)

	// You can disable various things by changing bool on the default parser
	// (or your own parser if you have created one).
	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	// You can set a help prepend or append on the default parser.
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/sty-holdings/utilities"

	// Add a flag to the main program (this will be available in all subcommands as well).
	flaggy.String(&natsCABundle, "a", "ncabun", "The filename of your NATS SSL CA Bundle.")
	flaggy.String(&natsCertificate, "c", "ncert", "The filename of your NATS SSL certificate.")
	flaggy.String(&fbCredentials, "f", "fbcred", "The filename of your firebase credentials.")
	flaggy.String(&analyzeProcess, "g", "answer", "Do you want to analyze or process the question? 'A' | 'P'. 'P' required -f.")
	flaggy.String(&natsCertKey, "k", "nkey", "The filename of your NATS SSL certificate key (Private key).")
	flaggy.String(&natsCredentials, "n", "ncreds", "The filename of your NATS credentials.")
	flaggy.String(&natsPort, "p", "nport", "The NATS port to connect to the NATS Server.")
	flaggy.String(&question, "q", "question", "Submit a question.")
	flaggy.String(&secretKey, "s", "secret", "Secret key (base64) to encrypt your message.")
	flaggy.Bool(&testingOn, "t", "testingOn", "This puts the server into testing mode.")
	flaggy.String(&natsURL, "u", "nurl", "The the URL for the NATS Server.")

	// Set the version and parse all inputs into variables.
	flaggy.Parse()
}

//goland:noinspection GoBoolExpressions
func main() {

	var (
		errorInfo     errs.ErrorInfo
		tInstanceName string
		tNATSConfig   ns.NATSConfiguration
		tNATSConnPtr  *nats.Conn
	)

	fmt.Println()

	if question == ctv.VAL_EMPTY {
		flaggy.ShowHelpAndExit("You must provide a question.")
	}

	if fbCredentials == ctv.VAL_EMPTY && question == ctv.VAL_EMPTY {
		flaggy.ShowHelpAndExit("When you are not entering a question, you must provide a credentials filename.")
		os.Exit(1)
	}
	if natsCABundle == ctv.VAL_EMPTY {
		flaggy.ShowHelpAndExit("You must provide your NATS credentials filename.")
		os.Exit(1)
	}
	if natsCertificate == ctv.VAL_EMPTY {
		flaggy.ShowHelpAndExit("You must provide your NATS credentials filename.")
		os.Exit(1)
	}
	if natsCertKey == ctv.VAL_EMPTY {
		flaggy.ShowHelpAndExit("You must provide your NATS credentials filename.")
		os.Exit(1)
	}
	if natsCredentials == ctv.VAL_EMPTY {
		flaggy.ShowHelpAndExit("You must provide your NATS credentials filename.")
		os.Exit(1)
	}
	if natsPort == ctv.VAL_EMPTY {
		flaggy.ShowHelpAndExit("You must provide a port number.")
		os.Exit(1)
	}
	if natsURL == ctv.VAL_EMPTY {
		flaggy.ShowHelpAndExit("You must provide the NATS Server URL.")
		os.Exit(1)
	}

	if tInstanceName, errorInfo = ns.BuildInstanceName(ns.METHOD_DASHES, "submit-question"); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		os.Exit(1)
	}
	tNATSConfig = ns.NATSConfiguration{
		NATSCredentialsFilename: natsCredentials,
		NATSPort:                natsPort,
		NATSTLSInfo: jwts.TLSInfo{
			TLSCertFQN:       natsCertificate,
			TLSPrivateKeyFQN: natsCertKey,
			TLSCABundleFQN:   natsCABundle,
		},
		NATSURL: natsURL,
	}
	if tNATSConnPtr, errorInfo = ns.GetConnection(tInstanceName, tNATSConfig); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		os.Exit(1)
	}

	if errorInfo = determineProcess(tNATSConnPtr, tInstanceName, analyzeProcess, question); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		os.Exit(1)
	}
}

func determineProcess(natsConnPtr *nats.Conn, instanceName string, analyzeProcess string, question string) (errorInfo errs.ErrorInfo) {

	var (
		//tInstanceName string
		tResponse string
	)

	if question == ctv.VAL_EMPTY && fbCredentials != ctv.VAL_EMPTY {
		processTrainingData(fbCredentials, natsConnPtr, instanceName)
	}

	if tResponse, errorInfo = sendRequest(natsConnPtr, instanceName, analyzeProcess, question); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		return
	}
	fmt.Printf("Reply: %s\n", tResponse)

	return
}

func processTrainingData(fbCredentials string, natsConnPtr *nats.Conn, instanceName string) {

	var (
		errorInfo    errs.ErrorInfo
		tAppPtr      *firebase.App
		tDocRefPtr   []*firestore.DocumentSnapshot
		tFSClientPtr *firestore.Client
		tResponse    string
	)

	fmt.Println()

	if tAppPtr, _, errorInfo = fbs.GetFirebaseAppAuthConnection(fbCredentials); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		return
	}
	if tFSClientPtr, errorInfo = fbs.GetFirestoreClientConnection(tAppPtr); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		return
	}
	if tDocRefPtr, errorInfo = fbs.GetAllDocuments(tFSClientPtr, ctv.DATASTORE_TRAINING_QUESTIONS); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		return
	}

	for _, snapshot := range tDocRefPtr {
		//
		// Get Question
		if !snapshot.Exists() {
			fmt.Printf("Document does not exist\n")
			continue
		}
		data := snapshot.Data()
		// Access data fields using their key, e.g.:
		fmt.Printf("Document ID: %s, Category: %s, Question: %s\n", snapshot.Ref.ID, data[ctv.FN_CATEGORY], data[ctv.FN_QUESTION])

		if tResponse, errorInfo = sendRequest(natsConnPtr, instanceName, "P", data[ctv.FN_QUESTION].(string)); errorInfo.Error != nil {
			errs.PrintErrorInfo(errorInfo)
			os.Exit(1)
		}

		fmt.Printf("Reply: %+v\n", tResponse)
	}

	return
}

func sendRequest(natsConnPtr *nats.Conn, instanceName string, analyzeProcess string, question string) (response string, errorInfo errs.ErrorInfo) {

	var (
		dMessage             string
		tSTYHClientIdB64     = "912c2c2c-a1f7-11ef-852b-85093fa0b49a"
		tQuestionJSON        []byte
		tEncryptedMessageB64 string
		tMessagePtr          *nats.Msg
		tReplyPtr            *nats.Msg
		tResponsePtr         *ns.NATSReply
	)

	if tEncryptedMessageB64, errorInfo = jwts.Encrypt(username, secretKey, question); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
	}

	if tQuestionJSON, errorInfo.Error = json.Marshal(
		ns.AnalyzeQuestionRequest{Question: tEncryptedMessageB64},
	); errorInfo.Error != nil {
		errs.PrintError(errs.ErrMessageJSONInvalid, ctv.VAL_EMPTY)
		return
	}

	tMessagePtr = &nats.Msg{
		Header: make(nats.Header),
		Data:   tQuestionJSON,
	}
	tMessagePtr.Header.Add(ctv.FN_USERNAME, username)
	tMessagePtr.Header.Add(ctv.FN_STYH_CLIENT_ID, tSTYHClientIdB64)

	if analyzeProcess == "P" {
		tMessagePtr.Subject = ctv.SUB_GEMINI_GET_MY_ANSWER
	} else {
		tMessagePtr.Subject = ctv.SUB_GEMINI_ANALYZE_QUESTION
	}

	if tReplyPtr, errorInfo = ns.RequestWithHeader(natsConnPtr, instanceName, tMessagePtr, 3); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		return
	}

	if errorInfo = ns.UnmarshalMessageData("sendRequest", tReplyPtr, &tResponsePtr); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		return
	}

	if dMessage, errorInfo = jwts.Decrypt(username, secretKey, tResponsePtr.Response.(string)); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		return
	}
	fmt.Printf("Decrypted Response: %s\n", dMessage)

	return
}
