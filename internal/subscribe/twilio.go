package subscribe

import (
	"log"
	"net/http"

	"github.com/gorilla/schema"

	"github.com/twilio/twilio-go"
)

type twilioSubscriber struct {
	client   *twilio.RestClient
	sendFrom string
	sendTo   string
	msgs     chan string
}

type NewTwilioSubscriberParams struct {
	Client   *twilio.RestClient
	SendFrom string
	SendTo   string
}

func NewTwilioSubscriber(params *NewTwilioSubscriberParams) *twilioSubscriber {
	return &twilioSubscriber{
		client:   params.Client,
		sendFrom: params.SendFrom,
		sendTo:   params.SendTo,
		msgs:     make(chan string),
	}
}

// @TODO: return struct with errs?
func (s twilioSubscriber) Subscribe() (<-chan string, error) {
	go func() {
		http.HandleFunc("/", s.onWebhook)
		// @TODO: err handle
		http.ListenAndServe(":8080", nil)
	}()

	return s.msgs, nil
}

func (s twilioSubscriber) onWebhook(w http.ResponseWriter, r *http.Request) {
	// @TODO: validate req
	msg, err := s.parseRequest(r)
	if err != nil {
		log.Printf("an error occurred: %v", err)
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

var _ Subscriber = twilioSubscriber{}
