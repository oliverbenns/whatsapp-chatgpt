package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

type Env struct {
	TwilioAccountSid string
	TwilioAuthToken  string
	TwilioSendFrom   string
	TwilioSendTo     string
	OpenAiApiKey     string
}

func GetEnv() (*Env, error) {
	env := &Env{
		TwilioAccountSid: os.Getenv("TWILIO_ACCOUNT_SID"),
		TwilioAuthToken:  os.Getenv("TWILIO_AUTH_TOKEN"),
		TwilioSendFrom:   os.Getenv("TWILIO_SEND_FROM"),
		TwilioSendTo:     os.Getenv("TWILIO_SEND_TO"),
		OpenAiApiKey:     os.Getenv("OPENAI_API_KEY"),
	}

	invalid := env.TwilioAccountSid == "" || env.TwilioAuthToken == "" || env.TwilioSendFrom == "" || env.TwilioSendTo == "" || env.OpenAiApiKey == ""

	if invalid {
		return nil, fmt.Errorf("missing credentials")
	}

	return env, nil
}

func main() {
	env, err := GetEnv()
	if err != nil {
		panic(err)
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: env.TwilioAccountSid,
		Password: env.TwilioAuthToken,
	})

	params := &twilioApi.CreateMessageParams{}
	to := fmt.Sprintf("whatsapp:%s", env.TwilioSendTo)
	params.SetTo(to)
	from := fmt.Sprintf("whatsapp:%s", env.TwilioSendFrom)
	params.SetFrom(from)
	params.SetBody("Hello from Go!")

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		fmt.Println("Error sending message: " + err.Error())
	} else {
		response, _ := json.Marshal(*resp)
		fmt.Println("Response: " + string(response))
	}
}
