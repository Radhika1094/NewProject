package smsService

import (
	model "MahaVan/MahaVanServer/api/model"
	"net/http"
	"net/url"
	"strings"

	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/confighelper"
	"corelab.mkcl.org/MKCLOS/coredevelopmentplatform/coreospackage/logginghelper"
)

func SendSMSService(smsBody model.SmsBody) {
	var Url *url.URL
	Url, parseErr := url.Parse(confighelper.GetConfig("SMSAPIURL"))
	if parseErr != nil {
		logginghelper.LogError("Unable to connect with SMS gateway : ", parseErr)
		// return false, parseErr
	}
	parameters := url.Values{}
	parameters.Add("username", confighelper.GetConfig("UserName"))
	parameters.Add("password", confighelper.GetConfig("password"))
	parameters.Add("to", "91"+smsBody.MobileNumber)
	parameters.Add("from", confighelper.GetConfig("CDMAHeader"))
	// parameters.Add("CDMAHeader", confighelper.GetConfig("CDMAHeader"))
	parameters.Add("text", smsBody.SmsText)
	// parameters.Add("text", "test otp :123456")
	logginghelper.LogInfo("length og sms:", len(smsBody.SmsText))
	Url.RawQuery = parameters.Encode()
	logginghelper.LogDebug("Encoded url ", Url.String())
	response, errSendSMS := http.Get(Url.String())
	if errSendSMS != nil {
		logginghelper.LogError("sendSMS : ", errSendSMS.Error())
		if response != nil {
			response.Body.Close()
		}
		// return false, errSendSMS
	}
	// response.Body.Close()
	// return true, nil
}

// (author: avadhutp) function to call sms sending function based on action
func SendActionSpecificSMSService(smsData model.SMSData) {
	CreateAndSendSMSService(smsData)
}

// (author: avadhutp) function to generate the actual text to be sent to receiver
func CreateAndSendSMSService(smsData model.SMSData) {
	logginghelper.LogDebug("CreateAndSendSMSService called with data: ", smsData)
	smsBody := model.SmsBody{}
	smsBody.SmsText = confighelper.GetConfig("smsTemplates." + smsData.Action)
	logginghelper.LogDebug(smsBody.SmsText)
	for key, value := range smsData.SMSData {
		smsBody.SmsText = strings.Replace(smsBody.SmsText, "$"+key, value, 1)
	}
	logginghelper.LogDebug("SMS text to be send :" + smsBody.SmsText)
	smsBody.MobileNumber = smsData.MobileNumber
	SendSMSService(smsBody)
	return
}
