package main

import (
	"fmt"
	"os"

	"github.com/oliverbenns/whatsapp-chatgpt/internal/prompt"
	"github.com/oliverbenns/whatsapp-chatgpt/internal/publish"
	"github.com/sashabaranov/go-openai"
	"github.com/twilio/twilio-go"
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

	openAiClient := openai.NewClient(env.OpenAiApiKey)

	twilioRestClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: env.TwilioAccountSid,
		Password: env.TwilioAuthToken,
	})

	prompter := prompt.NewOpenAiPrompter(&prompt.NewOpenAiPrompterParams{
		Client: openAiClient,
	})

	publisher := publish.NewTwilioPublisher(&publish.NewTwilioPublisherParams{
		Client:   twilioRestClient,
		SendTo:   env.TwilioSendTo,
		SendFrom: env.TwilioSendFrom,
	})

	res, err := prompter.Prompt("Hello!")
	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return
	}

	err = publisher.Publish(res)
	if err != nil {
		fmt.Println("Error sending message: " + err.Error())
	}
}
