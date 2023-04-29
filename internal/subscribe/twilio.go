package subscribe

import (
	"fmt"
	"net/http"

	"github.com/gorilla/schema"

	"github.com/twilio/twilio-go"
)

type twilioSubscriber struct {
	client      *twilio.RestClient
	sendFrom    string
	sendTo      string
	webhookPath string
	webhookPort int
	msgs        chan string
	errs        chan error
}

type NewTwilioSubscriberParams struct {
	Client      *twilio.RestClient
	SendFrom    string
	SendTo      string
	WebhookPath string
	WebhookPort int
}

func NewTwilioSubscriber(params *NewTwilioSubscriberParams) *twilioSubscriber {
	return &twilioSubscriber{
		client:      params.Client,
		sendFrom:    params.SendFrom,
		sendTo:      params.SendTo,
		webhookPath: params.WebhookPath,
		webhookPort: params.WebhookPort,
		msgs:        make(chan string),
		errs:        make(chan error),
	}
}

func (s twilioSubscriber) Subscribe() (<-chan string, <-chan error) {
	go func() {
		path := s.webhookPath
		addr := fmt.Sprintf(":%d", s.webhookPort)

		http.HandleFunc(path, s.onWebhook)

		err := http.ListenAndServe(addr, nil)
		if err != nil {
			s.errs <- err
		}
	}()

	return s.msgs, s.errs
}

func (s twilioSubscriber) onWebhook(w http.ResponseWriter, r *http.Request) {
	err := s.validateRequest(r)
	if err != nil {
		s.errs <- fmt.Errorf("could not validate request: %w", err)
		return
	}

	msg, err := s.parseRequest(r)
	if err != nil {
		s.errs <- fmt.Errorf("could not parse request: %w", err)
		return
	}

	err = s.validateMsg(msg)
	if err != nil {
		s.errs <- fmt.Errorf("msg failed validation: %w", err)
		return
	}

	s.msgs <- msg.Body
}

type Msg struct {
	AccountSid       string
	ApiVersion       string
	Body             string
	From             string
	To               string
	MessageSid       string
	NumMedia         int
	NumSegments      int
	ProfileName      string
	ReferralNumMedia int
	SmsMessageSid    string
	SmsSid           string
	SmsStatus        string
	WaId             string
}

func (s twilioSubscriber) validateRequest(r *http.Request) error {
	// @TODO: validate req signature with auth token.
	// https://www.twilio.com/docs/usage/webhooks/webhooks-security
	// We need to know the webhook url ahead of time
	// in order to do this. Difficult when ngrok tunnelling or using
	// generated cloud run urls that happen after deploy.
	return nil
}

func (s twilioSubscriber) parseRequest(r *http.Request) (*Msg, error) {
	err := r.ParseForm()
	if err != nil {
		return nil, err
	}
	msg := &Msg{}
	decoder := schema.NewDecoder()

	err = decoder.Decode(msg, r.Form)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func (s twilioSubscriber) validateMsg(msg *Msg) error {
	to := fmt.Sprintf("whatsapp:%s", s.sendTo)
	from := fmt.Sprintf("whatsapp:%s", s.sendFrom)

	if to != msg.To {
		return fmt.Errorf("invalid 'to' param: %s", msg.To)
	}

	if from != msg.From {
		return fmt.Errorf("invalid 'from' param: %s", msg.From)
	}

	return nil
}

var _ Subscriber = twilioSubscriber{}
