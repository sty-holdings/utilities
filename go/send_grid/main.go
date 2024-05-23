package main

import (
	"encoding/json"
	"fmt"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"log"
	"os"
)

// Retrievealltransactionaltemplates : Retrieve all transactional templates (legacy & dynamic).
// GET /templates
func Retrievealltransactionaltemplates() {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	host := "https://api.sendgrid.com"
	request := sendgrid.GetRequest(apiKey, "/v3/templates", host)
	request.Method = "GET"
	queryParams := make(map[string]string)
	queryParams["generations"] = "legacy,dynamic"
	request.QueryParams = queryParams
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

// Retrieveasingletransactionaltemplate : Retrieve a single transactional template.
// GET /templates/{template_id}
func Retrieveasingletransactionaltemplate(apiKey string) {
	host := "https://api.sendgrid.com"
	request := sendgrid.GetRequest(apiKey, "/v3/templates/{template_id}", host)
	request.Method = "GET"
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

// Retrieveaspecifictransactionaltemplateversion : Retrieve a specific transactional template version.
// GET /templates/{template_id}/versions/{version_id}
func Retrieveaspecifictransactionaltemplateversion() {
	apiKey := os.Getenv("SENDGRID_API_KEY")
	host := "https://api.sendgrid.com"
	request := sendgrid.GetRequest(apiKey, "/v3/templates/{template_id}/versions/{version_id}", host)
	request.Method = "GET"
	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}

func main() {

	b, err := os.ReadFile("/Users/syacko/workspace/styh-dev/src/albert/keys/development/.keys/savup-development-sendgrid.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	type SGKey struct {
		SendgridKey string `json:"sendgrid_key"`
	}

	var x SGKey

	_ = json.Unmarshal(b, &x)

	tDTD := make(map[string]interface{})
	tDTD["su_first_name"] = "Scott"

	host := "https://api.sendgrid.com"
	request := sendgrid.GetRequest(x.SendgridKey, "/v3/mail/send", host)
	// request := sendgrid.GetRequest(string(b), "/v3/mail/send", host)
	request.Method = "POST"
	m := mail.NewV3Mail()

	address := "verification@sty-holdings.com"
	name := "SavUp Verification"
	e := mail.NewEmail(name, address)
	m.SetFrom(e)
	m.Subject = "Your Example Verification"

	p1 := mail.NewPersonalization()
	tos1 := []*mail.Email{
		mail.NewEmail("Scott Yacko", "scott@yackofamily.com"),
	}
	p1.AddTos(tos1...)
	tos2 := []*mail.Email{
		mail.NewEmail("Scott (gmail) Yacko", "syacko@gmail.com"),
	}
	p1.AddTos(tos2...)
	// ccs1 := []*mail.Email{
	// 	mail.NewEmail("Scott Yacko", "syacko@gmail.com"),
	// }
	// p1.AddCCs(ccs1...)
	p1.DynamicTemplateData = tDTD
	// bccs1 := []*mail.Email{
	// 	mail.NewEmail("Jim Doe", "james_doe@example.com"),
	// }
	// p1.AddBCCs(bccs1...)
	m.AddPersonalizations(p1)

	c1 := mail.NewContent("text/html", "<p>Hello from Twilio SendGrid!</p><p>Sending with the email service trusted by developers and marketers for <strong>time-savings</strong>, <strong>scalability</strong>, and <strong>delivery expertise</strong>.</p><p>%open-track%</p>")
	m.AddContent(c1)

	// a1 := mail.NewAttachment()
	// a1.SetContent("PCFET0NUWVBFIGh0bWw+CjxodG1sIGxhbmc9ImVuIj4KCiAgICA8aGVhZD4KICAgICAgICA8bWV0YSBjaGFyc2V0PSJVVEYtOCI+CiAgICAgICAgPG1ldGEgaHR0cC1lcXVpdj0iWC1VQS1Db21wYXRpYmxlIiBjb250ZW50PSJJRT1lZGdlIj4KICAgICAgICA8bWV0YSBuYW1lPSJ2aWV3cG9ydCIgY29udGVudD0id2lkdGg9ZGV2aWNlLXdpZHRoLCBpbml0aWFsLXNjYWxlPTEuMCI+CiAgICAgICAgPHRpdGxlPkRvY3VtZW50PC90aXRsZT4KICAgIDwvaGVhZD4KCiAgICA8Ym9keT4KCiAgICA8L2JvZHk+Cgo8L2h0bWw+Cg==")
	// a1.SetFilename("index.html")
	// a1.SetType("text/html")
	// a1.SetDisposition("attachment")
	// m.AddAttachment(a1)
	//
	// m.AddCategories("cake")
	// m.AddCategories("pie")
	// m.AddCategories("baking")
	// m.SetSendAt(1617260400)

	// asm := mail.NewASM()
	// asm.SetGroupID(25384)
	// asm.AddGroupsToDisplay(25384)
	// m.SetASM(asm)

	// mailSettings := mail.NewMailSettings()
	// bypassListManagementSetting := mail.NewSetting(false)
	// mailSettings.SetBypassListManagement(bypassListManagementSetting)
	// footerSetting := mail.NewFooterSetting()
	// footerSetting.SetEnable(false)
	// mailSettings.SetFooter(footerSetting)
	// sandboxModeSetting := mail.NewSetting(false)
	// mailSettings.SetSandboxMode(sandboxModeSetting)
	// m.SetMailSettings(mailSettings)

	// trackingSettings := mail.NewTrackingSettings()
	// clickTrackingSetting := mail.NewClickTrackingSetting()
	// clickTrackingSetting.SetEnable(true)
	// clickTrackingSetting.SetEnableText(false)
	// trackingSettings.SetClickTracking(clickTrackingSetting)
	// openTrackingSetting := mail.NewOpenTrackingSetting()
	// openTrackingSetting.SetEnable(true)
	// openTrackingSetting.SetSubstitutionTag("%open-track%")
	// trackingSettings.SetOpenTracking(openTrackingSetting)
	// subscriptionTrackingSetting := mail.NewSubscriptionTrackingSetting()
	// subscriptionTrackingSetting.SetEnable(false)
	// trackingSettings.SetSubscriptionTracking(subscriptionTrackingSetting)
	// m.SetTrackingSettings(trackingSettings)

	replyToEmail := mail.NewEmail("SavUp Support", "support@sty-holdings.com")
	m.SetReplyTo(replyToEmail)

	m.SetTemplateID("d-1b95eb212cb7460f9ada4d0db09a2b3f")
	m.SetCustomArg("su_first_name", "Scott")
	m.SetCustomArg("short_url", "yyyyyyyyyyyyy")

	var Body = mail.GetRequestBody(m)
	request.Body = Body
	response, err := sendgrid.API(request)

	if err != nil {
		log.Println(request.Headers["Authorization"])
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		fmt.Println(response.Body)
		fmt.Println(response.Headers)
	}
}
