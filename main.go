package main

import (
	"github.com/oliverbenns/whatsapp-chatgpt/internal/env"
	"github.com/oliverbenns/whatsapp-chatgpt/internal/prompt"
	"github.com/oliverbenns/whatsapp-chatgpt/internal/publish"
	"github.com/oliverbenns/whatsapp-chatgpt/internal/service"
	"github.com/oliverbenns/whatsapp-chatgpt/internal/subscribe"
	"github.com/sashabaranov/go-openai"
	"github.com/twilio/twilio-go"
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}

	env, err := env.GetEnv()
	if err != nil {
		log.Fatal("could not get env", zap.Error(err))
	}

	openAiClient := openai.NewClient(env.OpenAiApiKey)

	twilioClient := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: env.TwilioAccountSid,
		Password: env.TwilioAuthToken,
	})

	subscriber := subscribe.NewTwilioSubscriber(&subscribe.NewTwilioSubscriberParams{
		Client: twilioClient,
		// yes, invert
		SendTo:      env.TwilioSendFrom,
		SendFrom:    env.TwilioSendTo,
		WebhookPort: 8080,
		WebhookPath: "/webhooks/twilio",
	})

	prompter := prompt.NewOpenAiPrompter(&prompt.NewOpenAiPrompterParams{
		Client: openAiClient,
	})

	publisher := publish.NewTwilioPublisher(&publish.NewTwilioPublisherParams{
		Client:   twilioClient,
		SendTo:   env.TwilioSendTo,
		SendFrom: env.TwilioSendFrom,
	})

	svc := service.NewService(&service.NewServiceParams{
		Subscriber: subscriber,
		Publisher:  publisher,
		Prompter:   prompter,
	})

	log.Info("starting service")
	err = svc.Start()
	if err != nil {
		log.Fatal("service error", zap.Error(err))

	}
}
