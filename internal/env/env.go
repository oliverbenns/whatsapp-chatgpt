package env

import (
	"fmt"
	"os"
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

	if env.TwilioAccountSid == "" {
		return nil, fmt.Errorf("TWILIO_ACCOUNT_SID not present")
	}

	if env.TwilioAuthToken == "" {
		return nil, fmt.Errorf("TWILIO_AUTH_TOKEN not present")
	}

	if env.TwilioSendFrom == "" {
		return nil, fmt.Errorf("TWILIO_SEND_FROM not present")
	}
	if env.TwilioSendTo == "" {
		return nil, fmt.Errorf("TWILIO_SEND_TO not present")
	}

	if env.OpenAiApiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY not present")
	}

	return env, nil
}
