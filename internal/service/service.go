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

func (s *Service) Start() error {
	msgs, errs := s.subscriber.Subscribe()

	for {
		select {
		case msg := <-msgs:
			err := s.processMsg(msg)
			if err != nil {
				return fmt.Errorf("could not process msg: %w", err)
			}

		case err := <-errs:
			errProcess := s.processErr(err)
			if errProcess != nil {
				return fmt.Errorf("could not process err: %w", errProcess)
			}

		}
	}
}

// @TODO: Lock and only process 1 msg at once
func (s *Service) processMsg(msg string) error {
	res, err := s.prompter.Prompt(msg)

	var msgToPublish string
	if err != nil {
		msgToPublish = err.Error()
	} else {
		msgToPublish = res
	}

	err = s.publisher.Publish(msgToPublish)
	if err != nil {
		return fmt.Errorf("could not publish msg: %w", err)
	}

	return nil
}

func (s *Service) processErr(err error) error {
	errPublish := s.publisher.Publish(err.Error())
	if errPublish != nil {
		return fmt.Errorf("could not publish err: %w", errPublish)
	}

	return nil
}
