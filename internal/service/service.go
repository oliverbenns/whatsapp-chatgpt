package service

import (
	"fmt"

	"github.com/oliverbenns/whatsapp-chatgpt/internal/prompt"
	"github.com/oliverbenns/whatsapp-chatgpt/internal/publish"
	"github.com/oliverbenns/whatsapp-chatgpt/internal/subscribe"
)

type Service struct {
	subscriber subscribe.Subscriber
	publisher  publish.Publisher
	prompter   prompt.Prompter
}

type NewServiceParams struct {
	Subscriber subscribe.Subscriber
	Publisher  publish.Publisher
	Prompter   prompt.Prompter
}

func NewService(params *NewServiceParams) *Service {
	return &Service{
		subscriber: params.Subscriber,
		publisher:  params.Publisher,
		prompter:   params.Prompter,
	}
}

func (s *Service) Start() <-chan bool {
	done := make(chan bool)
	msgs, err := s.subscriber.Subscribe()
	if err != nil {
		fmt.Printf("Subscribe error: %v\n", err)
		go func() { done <- true }()
		return done
	}

	for {
		select {
		case msg := <-msgs:
			s.processMsg(msg)
		}
	}
}

func (s *Service) processMsg(msg string) {
	res, err := s.prompter.Prompt(msg)
	if err != nil {
		fmt.Printf("prompt error: %v\n", err)
		return
	}

	err = s.publisher.Publish(res)
	if err != nil {
		fmt.Println("Error sending message: " + err.Error())
	}
}
