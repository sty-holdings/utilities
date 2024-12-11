/*
This utility administrates the Category SaaS Providers that are supported by DaveKnows.
*/
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"github.com/integrii/flaggy"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"

	ctv "github.com/sty-holdings/sharedServices/v2024/constantsTypesVars"
	errs "github.com/sty-holdings/sharedServices/v2024/errorServices"
	fbs "github.com/sty-holdings/sharedServices/v2024/firebaseServices"
)

type Config struct {
	FirebaseCredentials string `yaml:"firebase_credentials"` // The filename of your firebase credentials.
	TestingOn           bool   `yaml:"testing_on"`           // This puts the server into testing mode.
}

var (
	action         string
	configFilename string
	category       string
	subCategory    string
	saasProvider   string
	utilityName    = "Category SaaS Providers Admin"
	saasProviders  = make(map[string]map[string]map[string]string)
	//
)

func init() {

	appDescription := cases.Title(language.English).String(utilityName) + " process adding, deleting a Category Sub-Category SaaS Provider or listing all. Adding categories\n" +
		" and sub-categories will happen when the first SaaS Provider is added. Deleting categories and sub-categories are not supported."
	// Set your program's name and description.  These appear in help output.
	flaggy.SetName("\n" + utilityName) // "\n" is added to the start of the name to make the output easier to read.
	flaggy.SetDescription(appDescription)

	// You can disable various things by changing bool on the default parser
	// (or your own parser if you have created one).
	flaggy.DefaultParser.ShowHelpOnUnexpected = true

	// You can set a help prepend or append on the default parser.
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/sty-holdings/utilities"

	// Add a flag to the main program (this will be available in all subcommands as well).
	flaggy.String(&action, "a", "action", "The action is adding or deleting a SaaS Provider to a category / sub-category or listing all. (A = add | D = delete | L = List")
	flaggy.String(&configFilename, "c", "config", "The directory and filename of the configuration file.")
	flaggy.String(&category, "g", "category", "The category to process against.")
	flaggy.String(&subCategory, "s", "secret", "The sub-category to process against.")
	flaggy.String(&saasProvider, "p", "saaspro", "The SaaS Provider to add or delete")

	// Set the version and parse all inputs into variables.
	flaggy.Parse()
}

//goland:noinspection GoBoolExpressions
func main() {

	var (
		errorInfo    errs.ErrorInfo
		tAppPtr      *firebase.App
		tConfig      Config
		tData        []byte
		tFSClientPtr *firestore.Client
	)

	fmt.Println()
	action = strings.ToUpper(action)

	checkNotEmpty(action, "You must either add or delete a SaaS provider.")
	checkNotEmpty(configFilename, "You must provide a configuration filename.")
	if action == "A" || action == "D" {
		checkNotEmpty(category, "You must provide a category. ")
		checkNotEmpty(subCategory, "You must provide a sub-category. ")
		checkNotEmpty(saasProvider, "You must provide a SaaS provider. ")
	}

	if tConfig, errorInfo = loadConfig(configFilename); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		os.Exit(1)
	}

	if tAppPtr, _, errorInfo = fbs.GetFirebaseAppAuthConnection(tConfig.FirebaseCredentials); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		return
	}
	if tFSClientPtr, errorInfo = fbs.GetFirestoreClientConnection(tAppPtr); errorInfo.Error != nil {
		errs.PrintErrorInfo(errorInfo)
		return
	}

	if saasProviders, errorInfo = pullSupportSaaSProviders(tFSClientPtr); errorInfo.Error != nil {
		if !errors.Is(errs.ErrDocumentNotFound, errorInfo.Error) {
			errs.PrintErrorInfo(errorInfo)
			return
		}
	}

	switch action {
	case "A":
		processAddition(category, subCategory, saasProvider)
	case "D":
		processDeletion(category, subCategory, saasProvider)
	case "L":
		processList()
	default:
		flaggy.ShowHelpAndExit("You have selected an invalid action.")
	}

	tData, errorInfo.Error = json.Marshal(saasProviders)
	if errorInfo = fbs.SetDocument(tFSClientPtr, ctv.DATASTORE_REFERENCE_DATA, ctv.REF_SUPPORT_SAAS_PROVIDERS, map[any]interface{}{ctv.FN_JSON_STRING: string(tData)}); errorInfo.Error != nil {
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

func processAddition(category string, subCategory string, saasProvider string) {

	var (
		ok bool
	)

	if _, ok = saasProviders[category]; !ok {
		saasProviders[category] = make(map[string]map[string]string)
	}
	if _, ok = saasProviders[category][subCategory]; !ok {
		saasProviders[category][subCategory] = make(map[string]string)
	}

	saasProviders[category][subCategory][saasProvider] = ctv.TXT_NULL

	return
}

func processDeletion(category string, subCategory string, saasProvider string) {

	var (
		ok bool
	)

	if _, ok = saasProviders[category]; ok {
		if _, ok = saasProviders[category][subCategory]; ok {
			delete(saasProviders[category][subCategory], saasProvider)
		}
	}

	return
}

func processList() {

	var (
		errorInfo errs.ErrorInfo
		tData     []byte
	)

	tData, errorInfo.Error = json.MarshalIndent(saasProviders, ctv.VAL_EMPTY, ctv.SPACES_FOUR)
	fmt.Println(string(tData))

	return
}

// pullSupportSaaSProviders - will read the Firestore datastore and run the data as a tri-map structure.
//
//	Customer Messages: None
//	Errors: None
//	Verifications: None
func pullSupportSaaSProviders(tFSClientPtr *firestore.Client) (saasProviders map[string]map[string]map[string]string, errorInfo errs.ErrorInfo) {

	var (
		tData                interface{}
		tDocumentSnapshotPtr *firestore.DocumentSnapshot
	)

	if tDocumentSnapshotPtr, errorInfo = fbs.GetDocumentById(tFSClientPtr, ctv.DATASTORE_REFERENCE_DATA, ctv.REF_SUPPORT_SAAS_PROVIDERS); errorInfo.Error != nil {
		return
	}

	if tData, errorInfo.Error = tDocumentSnapshotPtr.DataAt(ctv.FN_JSON_STRING); errorInfo.Error != nil {
		return
	}
	saasProviders = make(map[string]map[string]map[string]string)
	errorInfo.Error = json.Unmarshal([]byte(tData.(string)), &saasProviders)

	return
}
