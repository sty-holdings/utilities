package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/integrii/flaggy"
	"github.com/nats-io/nats.go"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"

	ctv "github.com/sty-holdings/sharedServices/v2024/constantsTypesVars"
	errs "github.com/sty-holdings/sharedServices/v2024/errorServices"
	fbs "github.com/sty-holdings/sharedServices/v2024/firebaseServices"
	jwts "github.com/sty-holdings/sharedServices/v2024/jwtServices"
	ns "github.com/sty-holdings/sharedServices/v2024/natsSerices"
)

type Config struct {
	FirebaseCredentials string `yaml:"firebase_credentials"` // The filename of your firebase credentials.
	NatsInfo            struct {
		CABundle       string `yaml:"ca_bundle"`       // The filename of your NATS SSL CA Bundle.
		CertificateKey string `yaml:"certificate_key"` // The filename of your NATS SSL certificate key (Private key).
		Certificate    string `yaml:"certificate"`     // The filename of your NATS SSL certificate.
		Credentials    string `yaml:"credentials"`     // The filename of your NATS credentials.
		Port           string `yaml:"port"`            // The NATS port to connect to the NATS Server.
		URL            string `yaml:"url"`             // The URL for the NATS Server.
	} `yaml:"nats_info"`
	TestingOn bool `yaml:"testing_on"` // This puts the server into testing mode.
}

var (
	action         string
	configFilename string
	question       string
	secretKey      string
	testingOn      bool
	userId         string
	utilityName    = "Submit Question"
	//
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
	flaggy.String(&action, "a", "answer", "G - Get My Answer | S - single question | T - Test Questions")
	flaggy.String(&configFilename, "c", "config", "The directory and filename of the configuration file.")
	flaggy.String(&question, "q", "question", "Submit a question. Only valid for 'G' and 'S' actions.")
	flaggy.String(&secretKey, "s", "secret", "Secret key (base64) to encrypt your message.")
	flaggy.Bool(&testingOn, "t", "testingOn", "This puts the server into testing mode.")
	flaggy.String(&userId, "u", "userid", "The userid from the identity provider (Firebase Auth or AWS Cognito).")

	// Set the version and parse all inputs into variables.
	flaggy.Parse()
}

//goland:noinspection GoBoolExpressions
func main() {

	var (
		errorInfo     errs.ErrorInfo
		tConfig       Config
		tInstanceName string
		tNATSConfig   ns.NATSConfiguration
		tNATSConnPtr  *nats.Conn
		tSubject      string
	)

	fmt.Println()

	switch action {
	case "G":
		tSubject = ctv.SUB_HAL_GET_MY_ANSWER
		checkNotEmpty(question, "You must provide a question.")
	case "S", "T":
		tSubject = ctv.SUB_GEMINI_ANALYZE_QUESTION
		if action == "S" { // Only check for question if action is "S"
			checkNotEmpty(question, "You must provide a question.")
		}
	default:
		flaggy.ShowHelpAndExit("The action is invalid.")
	}

	checkNotEmpty(configFilename, "You must provide a configuration filename.")
	checkNotEmpty(secretKey, "You must provide the registered secret key for the user.")

	if userId == ctv.VAL_EMPTY {
		userId = "6oEtPOwn9hN2gNRYmGDQY3QXwrF2"
		fmt.Printf("The default userId (%s) will be used.\n", userId)
	}

	if tConfig, errorInfo = loadConfig(configFilename); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		os.Exit(1)
	}

	if tInstanceName, errorInfo = ns.BuildInstanceName(ns.METHOD_DASHES, "submit-question"); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		os.Exit(1)
	}
	tNATSConfig = ns.NATSConfiguration{
		NATSCredentialsFilename: tConfig.NatsInfo.Credentials,
		NATSPort:                tConfig.NatsInfo.Port,
		NATSTLSInfo: jwts.TLSInfo{
			TLSCertFQN:       tConfig.NatsInfo.Certificate,
			TLSPrivateKeyFQN: tConfig.NatsInfo.CertificateKey,
			TLSCABundleFQN:   tConfig.NatsInfo.CABundle,
		},
		NATSURL: tConfig.NatsInfo.URL,
	}
	if tNATSConnPtr, errorInfo = ns.GetConnection(tInstanceName, tNATSConfig); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		os.Exit(1)
	}

	if action == "T" {
		processTrainingData(tConfig.FirebaseCredentials, tNATSConnPtr, tInstanceName)
		return
	}

	fmt.Println("----------------------------")
	fmt.Printf("Question: %s\n", question)
	if _, errorInfo = sendRequest(tNATSConnPtr, tInstanceName, tSubject, question); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
	}
}

func checkNotEmpty(value string, message string) {
	if value == ctv.VAL_EMPTY {
		flaggy.ShowHelpAndExit(message)
	}
}

func loadConfig(configFilename string) (config Config, errorInfo errs.ErrorInfo) {

	var (
		tConfigFile *os.File
	)

	if tConfigFile, errorInfo.Error = os.Open(configFilename); errorInfo.Error != nil {
		log.Fatal(errorInfo.Error)
	}
	defer func(tConfigFile *os.File) {
		err := tConfigFile.Close()
		if err != nil {
			log.Fatalf("error: %v", err)
		}
	}(tConfigFile)

	decoder := yaml.NewDecoder(tConfigFile)
	errorInfo.Error = decoder.Decode(&config)

	return
}

func processTrainingData(fbCredentials string, natsConnPtr *nats.Conn, instanceName string) {

	var (
		errorInfo    errs.ErrorInfo
		tAppPtr      *firebase.App
		tDocRefPtr   []*firestore.DocumentSnapshot
		tFSClientPtr *firestore.Client
		tCounter     int
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
		tCounter++
		fmt.Println("----------------------------")
		fmt.Printf("Counter: %d, Document ID: %s, Category: %s, Question: %s\n", tCounter, snapshot.Ref.ID, data[ctv.FN_CATEGORY], data[ctv.FN_QUESTION])

		if _, errorInfo = sendRequest(natsConnPtr, instanceName, ctv.SUB_GEMINI_ANALYZE_QUESTION, data[ctv.FN_QUESTION].(string)); errorInfo.Error != nil {
			errs.PrintErrorInfo(errorInfo)
			os.Exit(1)
		}
	}

	return
}

func sendRequest(natsConnPtr *nats.Conn, instanceName string, subject string, question string) (response string, errorInfo errs.ErrorInfo) {

	var (
		dMessage             string
		tSTYHClientIdB64     = "912c2c2c-a1f7-11ef-852b-85093fa0b49a"
		tQuestionJSON        []byte
		tEncryptedMessageB64 string
		tMessagePtr          *nats.Msg
		tReplyPtr            *nats.Msg
		tResponsePtr         *ns.NATSReply
	)

	if tEncryptedMessageB64, errorInfo = jwts.Encrypt(userId, secretKey, question); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
	}

	if tQuestionJSON, errorInfo.Error = json.Marshal(
		ns.AnalyzeQuestionRequest{Question: tEncryptedMessageB64},
	); errorInfo.Error != nil {
		errs.PrintError(errs.ErrMessageJSONInvalid, ctv.VAL_EMPTY)
		return
	}

	tMessagePtr = &nats.Msg{
		Header:  make(nats.Header),
		Data:    tQuestionJSON,
		Subject: subject,
	}
	tMessagePtr.Header.Add(ctv.FN_USERNAME, userId)
	tMessagePtr.Header.Add(ctv.FN_STYH_CLIENT_ID, tSTYHClientIdB64)

	startTime := time.Now()
	if tReplyPtr, errorInfo = ns.RequestWithHeader(natsConnPtr, instanceName, tMessagePtr, 5*time.Second); errorInfo.Error != nil {
		elapsedTime := time.Since(startTime)
		fmt.Printf("Elapsed time: %v\n", elapsedTime)
		errs.PrintErrorInfo(errorInfo)
		return
	}
	elapsedTime := time.Since(startTime)
	fmt.Printf("Elapsed time: %v\n", elapsedTime)

	if errorInfo = ns.UnmarshalMessageData("sendRequest", tReplyPtr, &tResponsePtr); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		return
	}

	if tResponsePtr.Response != nil {
		if dMessage, errorInfo = jwts.Decrypt(userId, secretKey, tResponsePtr.Response.(string)); errorInfo.Error != nil {
			errs.PrintErrorInfo(errorInfo)
			return
		}
		fmt.Printf("Decrypted Response: %s\n", dMessage)
	}
	fmt.Printf("Error Response: %s\n", tResponsePtr.ErrorInfo)

	return
}
